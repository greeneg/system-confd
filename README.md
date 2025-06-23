# system-confd - A Unified API for Linux System Configuration

A plugin based service for setting configuration settings on Linux systems.

## Plugins - How Getting and Setting Configurations are Done

At the heart of SystemConfd, the service relies on "plugin" binaries that drive configuration. Each of the plugins installed and enabled on a system generates additional dynamic API endpoints to allow discovery of hardware, requesting the current configuration of a subsystem of a Linux installation, and setting its configuration parameters including steps required to apply the changes to the underlying configuration files. All with a unified rest-like API over a local protected UNIX socket.

This API will allow client tools to query and set configurations without having to be familiar with the syntax of a subsystem, or the tasks involved in allowing the configuration change to be applied. This greatly simplifies how to manage a system, either with graphical or command-line interfaces, or via any of the numerous configuration management tools available, such as Ansible, CFEngine, Chef, Puppet, or Salt.

### Plugin Architecture

All installed plugins require registration with SystemConfd's `plugins.json` file normally installed in `/etc/system-confd/`. The registration requires the following fields defined in the table below in the `plugins` array:

| key | type | description |
| --- | --- | --- |
| name | string | The name of the plugin, as would be displayed in an interface |
| path | string | The fully-qualified path to the plugin binary |
| metadata | string | The fully-qualified path to the plugin metadata JSON file |
| enabled | boolean | Whether the plugin is enabled or not |

Each plugin binary (or script) may be written in any language as long as it adheres to the architecture as documented below.

#### Metadata and Dynamic API Endpoints

By default, SystemConfd listens on a local socket in /var/run/system-confd as the root user. This socket internally works much like any REST-ful service with API endpoints.

When the service is started, it probes for installed and enabled plugins and then loads their metadata JSON files to generate the additional API endpoints specific to that plugin. These dynamic endpoints are parented under `/api/v1/system/plugins`.

For example, the `keyboard` plugin currently has the following API endpoints:

- keyboard/discover: a GET method endpoint that queries and returns attached keyboards on the system

#### Inputs

Plugins require taking all inputs through `<STDIN>` as JSON data, and do not allow passing arguments via positional arguments or flags.

#### Outputs