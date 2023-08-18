package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	v1log "github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/devportalservice"
	"github.com/bitrise-io/go-xcode/utility"
	"github.com/bitrise-io/go-xcode/v2/autocodesign"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/certdownloader"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/codesignasset"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/devportalclient"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/keychain"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/localcodesignasset"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/projectmanager"
	"github.com/bitrise-io/go-xcode/v2/codesign"
	"github.com/bitrise-io/go-xcode/xcodebuild"
)

func failf(format string, args ...interface{}) {
	v1log.Errorf(format, args...)
	os.Exit(1)
}

func downloadCertificates(certDownloader autocodesign.CertificateProvider, logger log.Logger) ([]certificateutil.CertificateInfoModel, error) {
	certificates, err := certDownloader.GetCertificates()
	if err != nil {
		return nil, fmt.Errorf("failed to download certificates: %s", err)
	}

	if len(certificates) == 0 {
		logger.Warnf("No certificates are uploaded.")

		return nil, nil
	}

	logger.Printf("%d certificates downloaded:", len(certificates))
	for _, cert := range certificates {
		logger.Printf("- %s", cert)
	}

	return certificates, nil
}

func main() {
	// Parse and validate inputs
	var cfg Config
	parser := stepconf.NewInputParser(env.NewRepository())
	if err := parser.Parse(&cfg); err != nil {
		failf("Config: %s", err)
	}
	stepconf.Print(cfg)

	logger := log.NewLogger()
	logger.EnableDebugLog(cfg.VerboseLog)
	v1log.SetEnableDebugLog(cfg.VerboseLog) // for compatibility

	cmdFactory := command.NewFactory(env.NewRepository())

	xcodebuildVersion, err := utility.GetXcodeVersion()
	if err != nil {
		failf("Failed to determine Xcode version: %s", err)
	}
	logger.Printf("%s (%s)", xcodebuildVersion.Version, xcodebuildVersion.BuildVersion)

	logger.Println()
	if xcodebuildVersion.MajorVersion >= 11 {
		// Resolve Swift package dependencies, so running -showBuildSettings is faster
		// Specifying a scheme is required for workspaces
		resolveDepsCmd := xcodebuild.NewResolvePackagesCommandModel(cfg.ProjectPath, cfg.Scheme, cfg.Configuration)
		if err := resolveDepsCmd.Run(); err != nil {
			logger.Warnf("%s", err)
		}
	}

	// Analyze project
	fmt.Println()
	logger.Infof("Analyzing project")
	project, err := projectmanager.NewProject(projectmanager.InitParams{
		ProjectOrWorkspacePath: cfg.ProjectPath,
		SchemeName:             cfg.Scheme,
		ConfigurationName:      cfg.Configuration,
	})
	if err != nil {
		failf(err.Error())
	}

	appLayout, err := project.GetAppLayout(cfg.SignUITestTargets)
	if err != nil {
		failf(err.Error())
	}

	authType, err := parseAuthType(cfg.BitriseConnection)
	if err != nil {
		failf("Invalid input: unexpected value for Bitrise Apple Developer Connection (%s)", cfg.BitriseConnection)
	}

	codesignInputs := codesign.Input{
		AuthType:                  authType,
		DistributionMethod:        cfg.Distribution,
		CertificateURLList:        cfg.CertificateURLList,
		CertificatePassphraseList: cfg.CertificatePassphraseList,
		KeychainPath:              cfg.KeychainPath,
		KeychainPassword:          cfg.KeychainPassword,
	}

	codesignConfig, err := codesign.ParseConfig(codesignInputs, cmdFactory)
	if err != nil {
		failf(err.Error())
	}

	var connection *devportalservice.AppleDeveloperConnection
	if cfg.BuildURL != "" && cfg.BuildAPIToken != "" {
		f := devportalclient.NewFactory(logger)
		connection, err = f.CreateBitriseConnection(cfg.BuildURL, cfg.BuildAPIToken)
		if err != nil {
			failf(err.Error())
		}
	} else {
		logger.Warnf(`Connected Apple Developer Portal Account not found: BITRISE_BUILD_URL and BITRISE_BUILD_API_TOKEN envs are not set. 
			The step will use the connection override inputs as a fallback. 
			For testing purposes please provide BITRISE_BUILD_URL as json file (file://path-to-json) while setting BITRISE_BUILD_API_TOKEN to any non-empty string.`)
	}

	connectionInputs := codesign.ConnectionOverrideInputs{
		APIKeyPath:     cfg.APIKeyPath,
		APIKeyID:       cfg.APIKeyID,
		APIKeyIssuerID: cfg.APIKeyIssuerID,
	}
	appleAuthCredentials, err := codesign.SelectConnectionCredentials(authType, connection, connectionInputs, logger)
	if err != nil {
		failf(err.Error())
	}

	keychain, err := keychain.New(cfg.KeychainPath, cfg.KeychainPassword, cmdFactory)
	if err != nil {
		failf(fmt.Sprintf("failed to initialize keychain: %s", err))
	}

	devPortalClientFactory := devportalclient.NewFactory(logger)
	certDownloader := certdownloader.NewDownloader(codesignConfig.CertificatesAndPassphrases, retry.NewHTTPClient().StandardClient())
	assetWriter := codesignasset.NewWriter(*keychain)
	localCodesignAssetManager := localcodesignasset.NewManager(localcodesignasset.NewProvisioningProfileProvider(), localcodesignasset.NewProvisioningProfileConverter())

	devPortalClient, err := devPortalClientFactory.Create(appleAuthCredentials, cfg.TeamID)
	if err != nil {
		failf(err.Error())
	}

	fmt.Println()
	logger.TDebugf("Downloading certificates")
	certs, err := downloadCertificates(certDownloader, logger)
	if err != nil {
		failf(err.Error())
	}

	typeToLocalCerts, err := autocodesign.GetValidLocalCertificates(certs)
	if err != nil {
		failf(err.Error())
	}

	// Create codesign manager
	manager := autocodesign.NewCodesignAssetManager(devPortalClient, assetWriter, localCodesignAssetManager)

	// Auto codesign
	distribution := cfg.DistributionType()
	var testDevices []devportalservice.TestDevice
	if cfg.RegisterTestDevices && connection != nil {
		testDevices = connection.TestDevices
	}
	codesignAssetsByDistributionType, err := manager.EnsureCodesignAssets(appLayout, autocodesign.CodesignAssetsOpts{
		DistributionType:        distribution,
		TypeToLocalCertificates: typeToLocalCerts,
		BitriseTestDevices:      testDevices,
		MinProfileValidityDays:  cfg.MinProfileDaysValid,
		VerboseLog:              cfg.VerboseLog,
	})
	if err != nil {
		failf(fmt.Sprintf("Automatic code signing failed: %s", err))
	}

	if err := project.ForceCodesignAssets(distribution, codesignAssetsByDistributionType); err != nil {
		failf(fmt.Sprintf("Failed to force codesign settings: %s", err))
	}

	// Export output
	fmt.Println()
	logger.Infof("Exporting outputs")

	teamID := codesignAssetsByDistributionType[distribution].Certificate.TeamID
	outputs := map[string]string{
		"BITRISE_EXPORT_METHOD":  cfg.Distribution,
		"BITRISE_DEVELOPER_TEAM": teamID,
	}

	settings, ok := codesignAssetsByDistributionType[autocodesign.Development]
	if ok {
		outputs["BITRISE_DEVELOPMENT_CODESIGN_IDENTITY"] = settings.Certificate.CommonName

		bundleID, err := project.MainTargetBundleID()
		if err != nil {
			failf("Failed to read bundle ID for the main target: %s", err)
		}
		profile, ok := settings.ArchivableTargetProfilesByBundleID[bundleID]
		if !ok {
			failf("No provisioning profile ensured for the main target")
		}

		outputs["BITRISE_DEVELOPMENT_PROFILE"] = profile.Attributes().UUID
	}

	if distribution != autocodesign.Development {
		settings, ok := codesignAssetsByDistributionType[distribution]
		if !ok {
			failf("No codesign settings ensured for the selected distribution type: %s", distribution)
		}

		outputs["BITRISE_PRODUCTION_CODESIGN_IDENTITY"] = settings.Certificate.CommonName

		bundleID, err := project.MainTargetBundleID()
		if err != nil {
			failf(err.Error())
		}
		profile, ok := settings.ArchivableTargetProfilesByBundleID[bundleID]
		if !ok {
			failf("No provisioning profile ensured for the main target")
		}

		outputs["BITRISE_PRODUCTION_PROFILE"] = profile.Attributes().UUID
	}

	for k, v := range outputs {
		logger.Donef("%s=%s", k, v)
		if err := tools.ExportEnvironmentWithEnvman(k, v); err != nil {
			failf("Failed to export %s=%s: %s", k, v, err)
		}
	}
}
