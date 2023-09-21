# swiftbar-jira

A simple plugin for [SwiftBar](https://swiftbar.app/) to run a JIRA search and display the resulting 
issues in the MacOS menu bar.

## Building

1. Set environment variable `SWIFTBAR_PLUGIN_DIR` to your SwiftBar plugin directory.
2. Create a configuration file `swiftbar-jira.yaml` based on the example. Set JIRA url and access token.
3. Run `make all` to compile the binary and copy the plugin to the plugin directory.
