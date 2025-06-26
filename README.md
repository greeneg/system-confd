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

When the service is started, it probes for installed and enabled plugins and then loads their metadata JSON files to generate the additional API endpoints specific to that plugin. These dynamic endpoints are parented under `/api/v1/system`.

For example, the `keyboard` plugin currently has the following API endpoints:

- keyboard/discover: a GET method endpoint that queries and returns attached keyboards on the system

Example keyboard.meta.json:
```json
{
  "version": "1.0.0",
  "name": "keyboard",
  "apiMountPoint": "/hardware",
  "description": "Keyboard configuration plugin",
  "author": "Gary L. Greene, Jr. <greeneg@yggdrasilsoft.com>",
  "license": "GPL-3.0",
  "apiPaths": {
    "/discover": {
      "method": "GET",
      "description": "List keyboards attached to the system",
      "invocationActionObject": {
        "pluginProtocol": "v1.0",
        "action": "discover"
      }
    },
    "/read": {
      "method": "GET",
      "description": "Retrieve the current keyboard input configuration",
      "invocationActionObject": {
        "pluginProtocol": "v1.0",
        "action": "readConfiguration"
      }
    },
    "/apply": {
      "method": "POST",
      "description": "Change settings for the current keyboard input configuration",
      "invocationActionObject": {
        "pluginProtocol": "v1.0",
        "action": "setConfiguration",
        "values": {
          "keymap": {
            "type": "string",
            "description": "The keyboard keymap to apply to the keyboard on the virtual console",
            "default": "us"
          },
          "xkbLayout": {
            "type": "string",
            "description": "The keyboard keymap to apply to the keyboard under X11 or Wayland",
            "default": "us"
          },
          "xkbModel": {
            "type": "string",
            "description": "The model of keyboard attached under X11 or Wayland",
            "default": "pc104"
          },
          "xkbOptions": {
            "type": "string",
            "description": "Extra options to apply to the keyboard for remapping keys, etc",
            "default": "terminate:ctrl_alt_bksp"
          },
          "kbdDelay": {
            "type": "int",
            "description": "Key delay time in milliseconds",
            "default": 250
          },
          "keyRepeatRate": {
            "type": "float",
            "description": "Key repeat rate",
            "default": 4.0
          },
          "numlockEnable": {
            "type": "string",
            "description": "Whether to enable or disable the numlock key on boot or let the bios setting apply",
            "default": "bios"
          },
          "scrollLockEnable": {
            "type": "bool",
            "description": "Whether to enable or disable the scroll lock key on boot",
            "default": false
          },
          "capsLockEnable": {
            "type": "bool",
            "description": "Whether to enable the caps lock key on boot",
            "default": false
          },
          "disableCapsLock": {
            "type": "bool",
            "description": "Whether to disable the caps lock key and make it act like a normal shift key",
            "default": false
          },
          "kbdTtyToApplySettings": {
            "type": "stringarray",
            "description": "Which TTYs the delay, repeat, numlock, scroll-lock, and capslock settings apply to",
            "default": [ "tty1", "tty2", "tty3", "tty4", "tty5", "tty6" ]
          }
        }
      }
    }
  }
}
```

All metadata JSON files must have the following attributes:

| key | type | description |
| --- | --- | --- |
| version | version string | The version of the plugin |
| name | string | The name of the plugin |
| apiMountPoint | string enum | The root path that the API is mounted at. The valid list of mount points for API extensions are in the table below |
| description | string | The description of what the plugin manages |
| author | name and email address string | The author name and email address |
| license | SPDX identifier string | The SPDX identifier string that matches the license of the plugin |
| apiPaths | API endpoint object | The API dynamic endpoint definition of the plugin. The object's payload is discussed in a later part of this README |

The valid list of API mount points:

| API root path | Description |
| --- | --- |
| `/api/v1/system/hardware` | Configuration related directly to hardware |
| `/api/v1/system/security` | Configuration related directly to system security |
| `/api/v1/system/services` | Configuration related to services that run on the host |
| `/api/v1/system/software` | Configuration of all other software elements |

The dynamic endpoints are generated by the main server code based on the metadata API paths to allow querying the system, the subsystem's current configuration, or applying settings to a given subsystem.

#### Inputs

Plugins require taking all inputs through `<STDIN>` as JSON data, and do not allow passing arguments via positional arguments or flags.

#### Outputs