# Manage iOS Code Signing

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-manage-ios-code-signing?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-manage-ios-code-signing/releases)

TODO

<details>
<summary>Description</summary>

TODO
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `apple_service_connection` | This input determines which Bitrise Apple service connection should be used for automatic code signing. Available values: - `api-key`: [Bitrise Apple Service connection with API Key.](https://devcenter.bitrise.io/getting-started/connecting-to-services/setting-up-connection-to-an-apple-service-with-api-key/) - `apple-id`: [Bitrise Apple Service connection with Apple ID.](https://devcenter.bitrise.io/getting-started/connecting-to-services/connecting-to-an-apple-service-with-apple-id/) | required | `api-key` |
| `distribution_method` | Describes how Xcode should export the archive. | required | `development` |
| `project_path` | Xcode Project (.xcodeproj) or Workspace (.xcworkspace) path. | required | `$BITRISE_PROJECT_PATH` |
| `scheme` | Xcode Scheme name. | required | `$BITRISE_SCHEME` |
| `configuration` | Xcode Build Configuration.  If not specified, the default Build Configuration will be used. |  |  |
| `sign_uitest_targets` | If this input is set, the Step will manage the codesign assets of the UITest targets (of the main Application) among with the main Application codesign assets. | required | `no` |
| `register_test_devices` | If this input is set, the Step will register the known test devices on Bitrise from team members with the Apple Developer Portal.  Note that setting this to yes may cause devices to be registered against your limited quantity of test devices in the Apple Developer Portal, which can only be removed once annually during your renewal window. | required | `no` |
| `min_profile_validity` | If this input is set to >0, the managed Provisioning Profile will be renewed if it expires within the configured number of days.  Otherwise the Step renews the managed Provisioning Profile if it is expired. | required | `0` |
| `certificate_url_list` | URL of the code signing certificate to download.  Multiple URLs can be specified, separated by a pipe (\|) character.  Local file path can be specified, using the file:// URL scheme. | required, sensitive | `$BITRISE_CERTIFICATE_URL` |
| `passphrase_list` | Passphrases for the provided code signing certificates.  Specify as many passphrases as many Code signing certificate URL provided, separated by a pipe (\|) character. | required, sensitive | `$BITRISE_CERTIFICATE_PASSPHRASE` |
| `keychain_path` | Path to the Keychain where the code signing certificates will be installed. | required | `$HOME/Library/Keychains/login.keychain` |
| `keychain_password` | Password for the provided Keychain. | required, sensitive | `$BITRISE_KEYCHAIN_PASSWORD` |
| `build_url` | URL of the current Bitrise build. |  | `$BITRISE_BUILD_URL` |
| `build_api_token` | API token to access Bitrise resources during the current build. | sensitive | `$BITRISE_BUILD_API_TOKEN` |
| `verbose_log` | If this input is set, the Step will produce verbose level log messages. | required | `no` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_EXPORT_METHOD` | Distribution type can be one of the following: `development`, `app-store`, `ad-hoc` or `enterprise`. |
| `BITRISE_DEVELOPER_TEAM` | The development team's ID, for example, `1MZX23ABCD4`. |
| `BITRISE_DEVELOPMENT_CODESIGN_IDENTITY` | The development codesign identity's name, for example, `iPhone Developer: Bitrise Bot (VV2J4SV8V4)`. |
| `BITRISE_PRODUCTION_CODESIGN_IDENTITY` | The production codesign identity's name, for example, `iPhone Distribution: Bitrise Bot (VV2J4SV8V4)`. |
| `BITRISE_DEVELOPMENT_PROFILE` | The development provisioning profile's UUID which belongs to the main target, for example, `c5be4123-1234-4f9d-9843-0d9be985a068`. |
| `BITRISE_PRODUCTION_PROFILE` | The production provisioning profile's UUID which belongs to the main target, for example, `c5be4123-1234-4f9d-9843-0d9be985a068`. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-manage-ios-code-signing/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-manage-ios-code-signing/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
