# Snap collector plugin - cinder

Snap plugin for collecting metrics from OpenStack Cinder module. 

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Snap's Global Config](#snaps-global-config)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started
Plugin collects metrics by communicating with OpenStack by REST API.
It can run locally on the host, or in proxy mode (communicating with the host via HTTP(S)). 

### System Requirements
 * OpenStack deployment available
 * Cinder V2 API
 
### Operating systems
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### Download cinder plugin binary:
You can get the pre-built binaries for your OS and architecture at Snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snapd` to access.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-cinder

Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-cinder
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/${GOOS}/${GOARCH}`


### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [Snap's Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-cinder/blob/master/README.md#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-cinder/blob/master/README.md#examples).

#### Suggestions
* It is not recommended to set interval for task less than 20 seconds. This may lead to overloading Cinder API with requests.

## Documentation
### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
intel/openstack/cinder/\<tenant_name\>/volumes/count | int | Total number of OpenStack volumes for given tenant
intel/openstack/cinder/\<tenant_name\>/volumes/bytes | int  | Total number of bytes used by OpenStack volumes for given tenant
intel/openstack/cinder/\<tenant_name\>/snapshots/count | int | Total number of OpenStack volumes snapshots for given tenant
intel/openstack/cinder/\<tenant_name\>/snapshots/bytes | int | Total number of bytes used by OpenStack volumes snapshots for given tenant
intel/openstack/cinder/\<tenant_name\>/limits/MaxTotalVolumeGigabytes | int64 | Tenant quota for volume size
intel/openstack/cinder/\<tenant_name\>/limits/MaxTotalVolumes | int64 | Tenant quota for number of volumes

### Snap's Global Config
Global configuration files are described in [Snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "cinder" in "collector" section and then specify following options:
- `"endpoint"` - URL for OpenStack Identity endpoint aka Keystone (ex. `"http://keystone.public.org:5000"`)
- `"user"` -  user name which has access to OpenStack. It is highly prefer to provide user with administrative privileges. Otherwise returned metrics may not be complete.
- `"password"` -  user password 
- `"tenant"` - name of project admin project. This parameter is optional for global config. It can be provided at later stage, in task manifest configuration section for metrics. 

See example Global Config in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-cinder/blob/master/examples/cfg/).

### Examples
Example running snap-plugin-collector-cinder plugin and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build/linux/x86_64
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

Create Global Config, see example in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-cinder/blob/master/examples/cfg/).

In one terminal window, open the Snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in cfg.json):
```
$ $SNAP_PATH/snapd -l 1 -t 0 --config examples/cfg/cfg.json
```
In another terminal window:

Load snap-plugin-collector-cinder plugin:
```
$ $SNAP_PATH/snapctl plugin load build/linux/x86_64/snap-plugin-collector-cinder
```
Download desired publisher plugin eg.
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
```
Load file plugin for publishing:
```
$ $SNAP_PATH/snapctl plugin load snap-plugin-publisher-file
```
See available metrics for your system:
```
$ $SNAP_PATH/snapctl metric list
```
Create a task manifest file to use snap-plugin-collector-cinder plugin (exemplary file in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-cinder/blob/master/examples/tasks/)):
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "60s"
    },
    "workflow": {
        "collect": {
            "metrics": {
		        "/intel/openstack/cinder/admin/limits/MaxTotalVolumeGigabytes": {},
		        "/intel/openstack/cinder/admin/volumes/count": {},
		        "/intel/openstack/cinder/admin/volumes/bytes": {},
		        "/intel/openstack/cinder/admin/snapshots/count": {},
		        "/intel/openstack/cinder/admin/snapshots/bytes": {}
           },
            "config": {
                "/intel/openstack/cinder": {
                    "tenant": "admin"
                }
            },
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/snap-cinder-file.log"
                    }
                }
            ]
        }
    }
}
```
Create a task:
```
$ $SNAP_PATH/snapctl task create -t examples/tasks/task.json
```

### Roadmap
There are few items on current roadmap for this plugin:
- quotable Cinder resources like backups and consistency groups
- number of volumes per volume type
- handling wildcard for tenant
- support for Cinder V1 API

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. The full project is at http://github.com/intelsdi-x/snap.

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik)
