package main

import (
	"fmt"
	"github.com/bitrise-io/go-xcode/autocodesign/devportalclient"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/go-xcode/autocodesign"
	"github.com/bitrise-io/go-xcode/autocodesign/certdownloader"
)

// Config holds the step inputs
type Config struct {
	Distribution  string `env:"distribution_method,opt[development,app-store,ad-hoc,enterprise]"`
	ProjectPath   string `env:"project_path,dir"`
	Scheme        string `env:"scheme,required"`
	Configuration string `env:"configuration"`

	BitriseConnection string `env:"apple_service_connection,opt[api-key,apple-id]"`

	RegisterTestDevices bool `env:"register_test_devices,opt[yes,no]"`
	MinProfileDaysValid int  `env:"min_profile_validity,required"`
	SignUITestTargets   bool `env:"sign_uitest_targets,opt[yes,no]"`

	CertificateURLList        string          `env:"certificate_url_list,required"`
	CertificatePassphraseList stepconf.Secret `env:"passphrase_list,required"`
	KeychainPath              string          `env:"keychain_path,required"`
	KeychainPassword          stepconf.Secret `env:"keychain_password,required"`
	BuildAPIToken             string          `env:"build_api_token"`
	BuildURL                  string          `env:"build_url"`

	VerboseLog bool `env:"verbose_log,opt[no,yes]"`
}

// DistributionType ...
func (c Config) DistributionType() autocodesign.DistributionType {
	return autocodesign.DistributionType(c.Distribution)
}

// ValidateCertificates validates if the number of certificate URLs matches those of passphrases
func (c Config) ValidateCertificates() ([]string, []string, error) {
	pfxURLs := splitAndClean(c.CertificateURLList, "|", true)
	passphrases := splitAndClean(string(c.CertificatePassphraseList), "|", false)

	if len(pfxURLs) != len(passphrases) {
		return nil, nil, fmt.Errorf("certificates count (%d) and passphrases count (%d) should match", len(pfxURLs), len(passphrases))
	}

	return pfxURLs, passphrases, nil
}

// Certificates returns an array of p12 file URLs and passphrases
func (c Config) Certificates() ([]certdownloader.CertificateAndPassphrase, error) {
	pfxURLs, passphrases, err := c.ValidateCertificates()
	if err != nil {
		return nil, err
	}

	files := make([]certdownloader.CertificateAndPassphrase, len(pfxURLs))
	for i, pfxURL := range pfxURLs {
		files[i] = certdownloader.CertificateAndPassphrase{
			URL:        pfxURL,
			Passphrase: passphrases[i],
		}
	}

	return files, nil
}

// SplitAndClean ...
func splitAndClean(list string, sep string, omitEmpty bool) (items []string) {
	return sliceutil.CleanWhitespace(strings.Split(list, sep), omitEmpty)
}

func parseClientType(bitriseConnection string) (devportalclient.ClientType, error) {
	switch bitriseConnection {
	case "api-key":
		return devportalclient.APIKeyClient, nil
	case "apple-id":
		return devportalclient.AppleIDClient, nil
	default:
		return 0, fmt.Errorf("invalid connection input: %s", bitriseConnection)
	}
}
