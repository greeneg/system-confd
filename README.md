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
  "supportsRootTargets": true,
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
| supportsRootTargets | boolean | | Whether the plugin supports using a different root path for configuration changes |
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
| rootTarget | string | | The full path to the targeted OS root to make changes on. This is only supported on plugins that support using a different OS root than `/` |
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
    "action": "discover",
    "rootTarget": "/"
}
```

This simple input JSON tells the plugin to run the discover endpoint code, which will return a discover output JSON with the data described in the metadata for the plugin.

In a `setConfig` action request, you'll see something like so:

```json
{
    "version": 1,
    "action": "setConfig",
    "rootTarget": "/",
    "attributes": {
        "numLockEnable": "on"
    }
}
```

This simple request asks the plugin to set the `numLockEnable` keyboard setting to `on`, and stores its configuration appropriately for the given Linux distribution the computer is running.

Finally, in a `readConfig` action request, you'll see something similar to this:

```json
{
    "version": 1,
    "action": "readConfig",
    "rootTarget": "/",
    "attributes": [
        "numLockEnable"
    ]
}
```

This request returns the value of the `numLockEnable` setting.

The structure of the three input JSONs are based on the following attributes detailed in the table below.

**Fig. 1**: Structure of a discovery object

| Key | Type | Accepted Values | Description |
| --- | --- | --- | --- |
| version | Integer | 1 | The version of the plugin protocol |
| action | String | discovery | The action that maps to the internal function to call in the plugin |
| rootTarget | String | | The full path to the targeted OS root to make changes on. This is only supported on plugins that support using a different OS root than `/` |

**Fig. 2**: Structure of a setConfig object

| Key | Type | Accepted Values | Description |
| --- | --- | --- | --- |
| version | Integer | 1 | The version of the plugin protocol |
| action | String | setConfig | The action that maps to the internal function to call in the plugin |
| rootTarget | String | | The full path to the targeted OS root to make changes on. This is only supported on plugins that support using a different OS root than `/` |
| attributes | Object of key/value entries | | The attributes to set. Note that these are specific to the plugin. |

**Fig. 3**: Structure of a readConfig object

| Key | Type | Accepted Values | Description |
| --- | --- | --- | --- |
| version | Integer | 1 | The version of the plugin protocol |
| action | String | setConfig | The action that maps to the internal function to call in the plugin |
| rootTarget | String | | The full path to the targeted OS root to make changes on. This is only supported on plugins that support using a different OS root than `/` |
| attributes | List | | The list of attributes that are desired to be returned. |

#### Outputs

The outputs of the plugins are also well-defined to ensure that errors, status notifications, and successes can be propagated to whatever client application is sending requests to the service.

A typical response after a request is submitted looks similar to the following:

```json
{
    "version": 1,
    "status": 200,
    "message": "Success",
    "hasError": false,
    "originalRequestUri": "",
    "error": "",
}
```

As can be seen, standard HTTP status codes are used. The 2xx error range is typically a success. The 3xx range for states that are in transition. The 4xx range for requests that are invalid from the client. Finally, the 5xx range denotes serious errors on the host side.

Responses that are errors will contain the error text in the error attribute's value. This typically will be the error message from the service or other aspects for errors in the 5xx range. For the 4xx range, the error will denote why the request is invalid. The 1xx, 2xx, and 3xx range will not supply an error message.

If a task for a plugin will take time and has potentially more data to send to the client, the 202 status code will be returned, and the response will look like the following:

```json
{
    "version": 1,
    "status": 202,
    "message": "Accepted/Paused",
    "hasError": false,
    "originalRequestUri": "socket:///var/lib/system-confd/system-conf.sock/api/v1/hardware/keyboard/apply",
    "triggerToken": "d37b7a3f-1918-479a-8e07-5d14a27ab9fe",
    "streamSocket": "/var/lib/system-confd/tmp/sock.AUQSunzs",
    "error": ""
}
```

The `streamSocket` is a socket the client can connect to that will send streamed data from the server to the client. Generally, this will be progress indication information. Note that this requires sending an additional request from the client to the plugin to unpause the task as a way to ensure that the content from the streaming socket is captured. Also note, that it is expected that the client send the unpause request, or the task will hang waiting for permission to unpause the task.

A typical unpause request contains the original request URI and trigger token to validate that the client is "authorized" to request that the task continue. A typical unpause request looks like the following:

```json
{
    "version": 1,
    "action": "continue",
    "originalRequestUri": "socket:/var/lib/system-confd/system-conf.sock/api/v1/hardware/keyboard/apply",
    "triggerToken": "d37b7a3f-1918-479a-8e07-5d14a27ab9fe"
}
```