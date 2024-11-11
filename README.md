# Manage iOS Code Signing

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-manage-ios-code-signing?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-manage-ios-code-signing/releases)

Automatically manage code signing assets before a build.

<details>
<summary>Description</summary>

The **Manage iOS Code Signing** Step takes care of setting up the required code signing assets before the project is built on Bitrise.
The Step:
- Downloads and installs certificates uploaded to Bitrise.
- Generates, updates and downloads the provisioning profiles needed for your iOS project.
- Verifies and registers the project's Bundle IDs on the Apple Developer Site.
- Registers the iOS or tvOS devices connected to your Bitrise account with the Apple Developer Site.

Use the **Manage iOS Code Signing** Step if, for example:
- You use Fastlane for your project.
- You use the **Ionic Archive** or the **Cordova Archive** build Steps in your project.
- You use a **Script** Step because your project has its own build scripts.
The **Manage iOS Code Signing** Step takes care of automatically code signing your project before it's built on Bitrise.

### Configuring the Step
Before you start, make sure:
- You've defined your Apple Developer Account to Bitrise.
- You've assigned an Apple Developer Account to your app.
- Make sure the Step is followed by another Step that needs iOS code signing.

1. **Apple services connection method**: Select the Apple service connection method you provided earlier on Bitrise; which is either the API Key or the Apple ID.
2. **Distribution method**: Select the method Xcode should sign your project: development, app-store, ad-hoc, or enterprise.
3. **Project path**: Add the path where the Xcode Project or Workspace is located.
4. **Scheme**: Add the scheme name you wish to archive your project later.
5. **Build configuration**:Specify Xcode Build Configuration. The Step will use the provided Build Configuration's Build Settings, to understand your project's code signing configuration. If not provided, the Archive action's default Build Configuration will be used.

If you want to set the Apple service connection credentials on the step-level (instead of using the one configured in the App Settings), use the Step inputs in the **App Store Connect connection override** category. Note that this only works if **Automatic code signing method** is set to `api-key`.

Under **Options**:
1. **Ensure code signing assets for UITest targets too**: If this input is set, the Step will manage the codesign settings of the UITest targets of the main Application.
2. **Register test devices on the Apple Developer Portal**: If this input is set, the Step will register known test devices from team members with the Apple Developer Portal. Note that setting this to `yes` may cause devices to be registered against your limited quantity of test devices in the Apple Developer Portal, which can only be removed once annually during your renewal window.

Under **Build environment**:
You do not have to change any sensitive Environment Variable if all your certificates are already uploaded to Bitrise. Should you store your code signing files somewhere else (for example, in a private repository), then you can set these inputs in the `bitrise.yml` file.

Under **Debugging**:
1. **Verbose logging***: You can set this input to `yes` to produce more informative logs.

### Troubleshooting:
- The **Manage iOS Code Signing** Step will fail if the correct Apple Developer Account is not connected to Bitrise or the right connection method is not selected in the **Apple service connection method** input within the Step.
- The **Manage iOS Code Signing** Step will also fail if the right code signing certificates are not uploaded to Bitrise. A Development type certificate is needed if the **Distribution method** input is set to `development`, otherwise a Distribution type certificate is needed. We recommend you upload one Development and one Distribution certificate, so that the Step can ensure code signing files for all the distribution methods.
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
| `apple_team_id` | The Apple Developer Portal team to use for downloading code signing assets.  Defining this is only required when Apple Service Connection method is set to `apple-id` and the connected account belongs to multiple teams. |  |  |
| `certificate_url_list` | URL of the code signing certificate to download.  Multiple URLs can be specified, separated by a pipe (\|) character.  Local file path can be specified, using the file:// URL scheme. | required, sensitive | `$BITRISE_CERTIFICATE_URL` |
| `passphrase_list` | Passphrases for the provided code signing certificates.  Specify as many passphrases as many Code signing certificate URL provided, separated by a pipe (\|) character.  Certificates without a passphrase: for using a single certificate, leave this step input empty. For multiple certificates, use the separator as if there was a passphrase (examples: `pass\|`, `\|pass\|`, `\|`) | sensitive | `$BITRISE_CERTIFICATE_PASSPHRASE` |
| `keychain_path` | Path to the Keychain where the code signing certificates will be installed. | required | `$HOME/Library/Keychains/login.keychain` |
| `keychain_password` | Password for the provided Keychain. | required, sensitive | `$BITRISE_KEYCHAIN_PASSWORD` |
| `build_url` | URL of the current Bitrise build. |  | `$BITRISE_BUILD_URL` |
| `build_api_token` | API token to access Bitrise resources during the current build. | sensitive | `$BITRISE_BUILD_API_TOKEN` |
| `api_key_path` | Local path or remote URL to the private key (p8 file) for App Store Connect API. This overrides the Bitrise-managed API connection, only set this input if you want to control the API connection on a step-level. Most of the time it's easier to set up the connection on the App Settings page on Bitrise. The input value can be a file path (eg. `$TMPDIR/private_key.p8`) or an HTTPS URL. This input only takes effect if the other two connection override inputs are set too (`api_key_id`, `api_key_issuer_id`). |  |  |
| `api_key_id` | Private key ID used for App Store Connect authentication. This overrides the Bitrise-managed API connection, only set this input if you want to control the API connection on a step-level. Most of the time it's easier to set up the connection on the App Settings page on Bitrise. This input only takes effect if the other two connection override inputs are set too (`api_key_path`, `api_key_issuer_id`). |  |  |
| `api_key_issuer_id` | Private key issuer ID used for App Store Connect authentication. This overrides the Bitrise-managed API connection, only set this input if you want to control the API connection on a step-level. Most of the time it's easier to set up the connection on the App Settings page on Bitrise. This input only takes effect if the other two connection override inputs are set too (`api_key_path`, `api_key_id`). |  |  |
| `api_key_enterprise_account` | Indicates if the account is an enterprise type. This overrides the Bitrise-managed API connection, only set this input if you know you have an enterprise account. | required | `no` |
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
