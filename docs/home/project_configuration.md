# Project Configuration

Bazooka's philosophy tries to have every build-related configuration stored and versionned in a text file besides the other project files.

However, some cases just don't fit this model, especially for configuration values required before even checking-out the project source code.

To this effect, bazooka comes with a facility to store a key value structured configuration in the server.

The CLI can be used to list, view, set and delete these configuration values using the `bzk project config` command.
Please refer to the CLI help messages for more details on using this command.

## Standard configuration keys
### bzk.scm.reuse

If this key is set to `true`, using the CLI for example:

`bzk project config set $PROJECT-NAME bzk.scm.reuse true`

Then instead of performing a fresh checkout of the project source code on every build, bazooka will instead checkout the project only once, and then perform an in-place update to the latest version.

This can be useful for large projects or when behind a slow network.

By default, the key is not set.