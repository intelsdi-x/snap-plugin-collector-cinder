# snap-plugin-collector-cinder

snap plugin for collecting metrics from OpenStack Cinder module. 
It collects metrics by communicating with OpenStack by REST API.
It can be used in- as well as out-of-bands.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

Plugin collects metrics by communicating with OpenStack by REST API.
It can be used in- as well as out-of-bands.

### System Requirements

 - Linux

### Installation
#### Download <cinder> plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-<cinder>
Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-<cinder>
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description (optional)
----------|-----------|-----------------------
intel/openstack/cinder/\<tenant_name\>/volumes/count | Total number of OpenStack volumes for given tenant
intel/openstack/cinder/\<tenant_name\>/volumes/bytes | Total number of bytes used by OpenStack volumes for given tenant
intel/openstack/cinder/\<tenant_name\>/snapshots/count | Total number of OpenStack volumes snapshots for given tenant
intel/openstack/cinder/\<tenant_name\>/snapshots/bytes | Total number of bytes used by OpenStack volumes snapshots for given tenant
intel/openstack/cinder/\<tenant_name\>/limits/MaxTotalVolumeGigabytes | Tenant quota for volume size
intel/openstack/cinder/\<tenant_name\>/limits/MaxTotalVolumes | Tenant quota for number of volumes

### Examples
Example task manifest to use <cinder> plugin:
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "5s"
    },
    "workflow": {
        "collect": {
            "metrics": {
		        "/intel/openstack/cinder/demo/limits/MaxTotalVolumeGigabytes": {},
		        "/intel/openstack/cinder/demo/volumes/count": {},
		        "/intel/openstack/cinder/demo/volumes/bytes": {},
		        "/intel/openstack/cinder/demo/snapshots/count": {},
		        "/intel/openstack/cinder/demo/snapshots/bytes": {}
           },
            "config": {
            },
            "process": null,
            "publish": null
        }
    }
}
```


### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Marcin Krolik](https://github.com/marcin-krolik)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.