package codesign

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-xcode/v2/autocodesign"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/certdownloader"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/codesignasset"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/keychain"
)

// Input ...
type Input struct {
	AuthType                     AuthType
	DistributionMethod           string
	CertificateURLList           string
	CertificatePassphraseList    stepconf.Secret
	KeychainPath                 string
	KeychainPassword             stepconf.Secret
	FallbackProvisioningProfiles string
}

// Config ...
type Config struct {
	CertificatesAndPassphrases   []certdownloader.CertificateAndPassphrase
	Keychain                     keychain.Keychain
	DistributionMethod           autocodesign.DistributionType
	FallbackProvisioningProfiles []string
}

// ParseConfig validates and parses step inputs related to code signing and returns with a Config
func ParseConfig(input Input, cmdFactory command.Factory) (Config, error) {
	certificatesAndPassphrases, err := parseCertificates(input)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse certificate URL and passphrase inputs: %s", err)
	}

	keychainWriter, err := keychain.New(input.KeychainPath, input.KeychainPassword, cmdFactory)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open keychain: %s", err)
	}

	fallbackProfiles, err := validateAndExpandProfilePaths(input.FallbackProvisioningProfiles)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse provisioning profiles: %w", err)
	}

	return Config{
		CertificatesAndPassphrases:   certificatesAndPassphrases,
		Keychain:                     *keychainWriter,
		DistributionMethod:           autocodesign.DistributionType(input.DistributionMethod),
		FallbackProvisioningProfiles: fallbackProfiles,
	}, nil
}

// parseCertificates returns an array of p12 file URLs and passphrases
func parseCertificates(input Input) ([]certdownloader.CertificateAndPassphrase, error) {
	if strings.TrimSpace(input.CertificateURLList) == "" {
		return nil, fmt.Errorf("code signing certificate URL: required input is not present")
	}
	if strings.TrimSpace(input.KeychainPath) == "" {
		return nil, fmt.Errorf("keychain path: required input is not present")
	}
	if strings.TrimSpace(string(input.KeychainPassword)) == "" {
		return nil, fmt.Errorf("keychain password: required input is not present")
	}

	pfxURLs, passphrases, err := validateCertificates(input.CertificateURLList, string(input.CertificatePassphraseList))
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

// validateCertificates validates if the number of certificate URLs matches those of passphrases
func validateCertificates(certURLList string, certPassphraseList string) ([]string, []string, error) {
	pfxURLs := splitAndClean(certURLList, "|", true)
	passphrases := splitAndClean(certPassphraseList, "|", false) // allow empty items because passphrase can be empty

	if len(pfxURLs) != len(passphrases) {
		return nil, nil, fmt.Errorf("certificate count (%d) and passphrase count (%d) should match", len(pfxURLs), len(passphrases))
	}

	return pfxURLs, passphrases, nil
}

// SplitAndClean ...
func splitAndClean(list string, sep string, omitEmpty bool) (items []string) {
	return sliceutil.CleanWhitespace(strings.Split(list, sep), omitEmpty)
}

// validateAndExpandProfilePaths validates and expands profilesList.
// profilesList must be a list of paths separated either by `|` or `\n`.
// List items must be a remote (https://) or local (file://) file paths,
// or a local directory (with no scheme).
// For directory list items, the contained profiles' path will be returned.
func validateAndExpandProfilePaths(profilesList string) ([]string, error) {
	profiles := splitAndClean(profilesList, "\n", true)
	if len(profiles) == 1 {
		profiles = splitAndClean(profiles[0], "|", true)
	}

	var validProfiles []string
	for _, profile := range profiles {
		profileURL, err := url.Parse(profile)
		if err != nil {
			return []string{}, fmt.Errorf("invalid provisioning profile URL (%s): %w", profile, err)
		}

		// When file or https scheme provided, will fetch as a file
		if profileURL.Scheme != "" {
			validProfiles = append(validProfiles, profile)
			continue
		}

		// If no scheme is provided, assuming it is a local directory
		profilesInDir, err := listProfilesInDirectory(profile)
		if err != nil {
			return []string{}, err
		}

		validProfiles = append(validProfiles, profilesInDir...)
	}

	return validProfiles, nil
}

func listProfilesInDirectory(dir string) ([]string, error) {
	exists, err := pathutil.IsDirExists(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to check if provisioning profile path (%s) exists: %w", dir, err)
	} else if !exists {
		return nil, fmt.Errorf("please provide remote (https://) or local (file://) provisioning profile file paths with a scheme, or a local directory without a scheme: profile directory (%s) does not exist", dir)
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to list profile directory: %w", err)
	}

	var profiles []string
	for _, dirEntry := range dirEntries {
		if dirEntry.Type().IsDir() || !dirEntry.Type().IsRegular() {
			continue
		}

		if strings.HasSuffix(dirEntry.Name(), codesignasset.ProfileIOSExtension) {
			profileURL := fmt.Sprintf("file://%s", filepath.Join(dir, dirEntry.Name()))
			profiles = append(profiles, profileURL)
		}
	}

	return profiles, nil
}
