# system-confd - A Unified API for Linux System Configuration

A plugin based service for setting configuration settings on Linux systems.

## Plugins - How Getting and Setting Configurations are Done

At the heart of SystemConfd, the service relies on "plugin" binaries that drive configuration. Each of the plugins installed and enabled on a system generates additional dynamic API endpoints to allow discovery of hardware, requesting the current configuration of a subsystem of a Linux installation, and setting its configuration parameters including steps required to apply the changes to the underlying configuration files. All with a unified rest-like API over a local protected UNIX socket.

This API will allow client tools to query and set configurations without having to be familiar with the syntax of a subsystem, or the tasks involved in allowing the configuration change to be applied. This greatly simplifies how to manage a system, either with graphical or command-line interfaces, or via any of the numerous configuration management tools available, such as Ansible, CFEngine, Chef, Puppet, or Salt.

### Plugin Architecture

All installed plugins require registration with SystemConfd's `plugins.json` file normally installed in `/etc/system-confd/`.

The file has a required `version` field that takes an integer and for now is set to `1` to denote the schema version of the plugin registration file.

The registration requires the following fields defined in the table below in the `plugins` array:

| Key | Type | Description |
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

| Path | Description |
| --- | --- |
| /api/v1/system/hardware/keyboard/discover | a GET method endpoint that queries and returns attached keyboards on the system |
| /api/v1/system/hardware/keyboard/read | a GET method endpoint to retrieve the current keyboard configuration |
| /api/v1/system/hardware/keyboard/apply | a POST method endpoint to set the desired keyboard configuration |

The described metadata JSON file uses a modified OAS v3 format.

Example keyboard.meta.json:
```json
{
  "version": 1,
  "name": "keyboard",
  "apiMountPoint": "/hardware",
  "apiName": "keyboard",
  "description": "Keyboard configuration plugin",
  "author": "Gary L. Greene, Jr. <greeneg@yggdrasilsoft.com>",
  "license": "GPL-3.0",
  "apiPaths": {
    "/discover": {
      "method": "GET",
      "description": "List keyboards attached to the system",
      "actionObject": {
        "pluginProtocol": 1,
        "action": "discover"
      }
    },
    "/read": {
      "method": "GET",
      "description": "Retrieve the current keyboard input configuration",
      "actionObject": {
        "pluginProtocol": 1,
        "action": "readConfig"
      }
    },
    "/apply": {
      "method": "POST",
      "description": "Change settings for the current keyboard input configuration",
      "actionObject": {
        "pluginProtocol": 1,
        "action": "setConfig",
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
            "type": "integer",
            "description": "Key delay time in milliseconds",
            "default": 250,
            "enum": [
              250,
              500,
              750,
              1000
            ]
          },
          "keyRepeatRate": {
            "type": "number",
            "description": "Key repeat rate",
            "default": 4.0,
            "enum": [
              2.0,
              2.1,
              2.3,
              2.5,
              2.7,
              3.0,
              3.3,
              3.7,
              4.0,
              4.3,
              4.6,
              5.0,
              5.5,
              6.0,
              6.7,
              7.5,
              8.0,
              8.6,
              9.2,
              10.0,
              10.9,
              12.0,
              13.3,
              15.0,
              16.0,
              17.1,
              18.5,
              20.0,
              21.8,
              24.0,
              26.7,
              30.0
            ]
          },
          "numlockEnable": {
            "type": "string",
            "description": "Whether to enable or disable the numlock key on boot or let the bios setting apply",
            "default": "bios",
            "enum": [
              "bios",
              "no",
              "yes"
            ]
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
            "type": "array",
            "description": "Which TTYs the delay, repeat, numlock, scroll-lock, and capslock settings apply to",
            "default": [ "tty1", "tty2", "tty3", "tty4", "tty5", "tty6" ],
            "items": {
              "type": "string",
              "enum": [
                "tty1",
                "tty2",
                "tty3",
                "tty4",
                "tty5",
                "tty6"
              ]
            }
          }
        }
      }
    }
  }
}
```

All metadata JSON files must have the following attributes:

| Key | Type | Accepted Values | Description |
| --- | --- | --- | --- |
| version | integer | 1 | The version of the plugin |
| name | string | | The name of the plugin |
| apiMountPoint | string | Described below | The root path that the API is mounted at. The valid list of mount points for API extensions are in the table below |
| apiName | string | | The name of the plugin in the URI for the REST-ful API |
| description | string | | The description of what the plugin manages |
| author | name and email address string | | The author name and email address |
| license | SPDX identifier string | | The SPDX identifier string that matches the license of the plugin |
| apiPaths | API endpoint object | | The API dynamic endpoint definition of the plugin. The object's payload is discussed in a later part of this README |

The valid list of API mount points:

| API root path | Description |
| --- | --- |
| `/api/v1/system/hardware` | Configuration related to system hardware |
| `/api/v1/system/networking` | Configuration related to system networking |
| `/api/v1/system/security` | Configuration related to system security |
| `/api/v1/system/services` | Configuration related to services that run on the host |
| `/api/v1/system/software` | Configuration of all other software elements |

The dynamic endpoints are generated by the main server code based on the metadata API paths to allow querying the system, the subsystem's current configuration, or applying settings to a given subsystem.

The `apiPaths` are organized as named endpoints from their path endpoints. Each of the paths have the following required attributes:

| Key | Type | Accepted Values | Description |
| --- | --- | --- | --- |
| method | string | One of the following:<br />- GET<br />- POST | The type of RESTful API method to accept on the API endpoint |
| description | string | | The description of what the API endpoint does |
| actionObject | action object | Described below | Describes the functions to call internally and their required payload if needed |

The **action object**, as described above, details the internal function call and payload required for it. The attributes of a configuration are described with the following attributes inside the action object:

| Key | Type | Accepted Values | Description |
| --- | --- | --- | --- |
| pluginProtocol | integer | 1 | The version of the plugin protocol to use for enforcing the attributes which must be present |
| action | string | one of the following:<br />- discovery<br />- readConfig<br />- setConfig | The action that maps to the internal plugin function to call |
| values | attribute object | | The object that describes the attributes that are allowed to be set |

The **attribute objects** are detailed in the table below:

| Key | Type | Accepted Values | Description |
| --- | --- | --- | --- |
| type | string | See the types table below | The data type accepted for the given attribute |
| description | string | | The description of the attribute item |
| default | any valid type | | The default value the attribute will return when initially queried |
| items | item object | | If the attribute contains an array, the `items` key defines the constituent members of the array |
| properties | properyy object | | If the attribute contains an object, the `properties` key defines the elements of the object |
| naximum | integer or number | | The maximum integer or number value that is valid to set the attribute to |
| minimum | integer or number | | The minimum integer or number value that is valid to set the attribute to |
| enum | array of values | | The specific allowed values for the attribute |

#### Inputs

Plugins require taking all inputs through `<STDIN>` as JSON data, and do not allow passing arguments via positional arguments or flags.

Using our previous example of the `keyboard` plugin, lets take a look at the format of the input JSON.

```json
{
    "version": 1,
    "action": "discover"
}
```

This simple input JSON tells the plugin to run the discover endpoint code, which will return a discover output JSON with the data described in the metadata for the plugin.

In a `setConfiguration` action request, you'll see something like so:

```json
{
    "version": 1,
    "action": "setConfiguration",
    "payload": {

    }
}
```

#### Outputs