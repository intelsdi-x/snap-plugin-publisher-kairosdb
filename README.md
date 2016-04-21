# snap publisher plugin - KairosDB 

This plugin supports pushing metrics into an KairosDB instance.

It's used in the [snap framework](http://github.com/intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Task manifest](#task-manifest)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

* [golang 1.5+](https://golang.org/dl/) for building plugin from source code

Support Matrix

- KairosDB V1 REST API

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation

#### Download the plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-kairosdb.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/blob/master/README.md#examples).

## Documentation

For details on KairosDB, please refer to [documentation](https://kairosdb.github.io/docs/build/html/index.html).

###Task manifest
Task manifest configuration is described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md). In order to use KairosDB publisher you have to add section "publish" then specify following options:
- `"host"` - KairosDB host address (ex. `"127.0.0.1"`)
- `"port"` -  KairosDB REST API port (ex. `"8083"`)

See example task manifest in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/blob/master/examples/tasks/).

### Examples
Example of use snap-collector-mock1 collector plugin and KairosDB publisher plugin to save collecting data in KairosDB.

Install and configure KairosDB, for help read [KarirosDB's documentation](https://kairosdb.github.io/docs/build/html/index.html).

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

In one terminal window, open the snap daemon (in this case with logging set to 1 and trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```

In another terminal window:

Load snap-collector-mock1 collector plugin:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-collector-mock1
```

Load snap-plugin-publisher-kairosdb publisher plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-publisher-kairosdb
```

See available metrics:
```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file to use snap-plugin-publisher-kairosdb plugin (exemplary files in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/blob/master/examples/tasks/)):
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/mock/foo": {},
                "/intel/mock/bar": {},
                "/intel/mock/*/baz": {}
            },
            "config": {
                "/intel/mock": {
                    "user": "root",
                    "password": "secret"
                }
            },
            "process": [
                {
                    "plugin_name": "passthru",
                    "process": null,
                    "publish": [
                        {
                            "plugin_name": "kairosdb",
                            "config": {
                                "host": "127.0.0.1",
                                "port": 2003
                            }
                        }
                    ]
                }
            ]
        }
    }
}
```
Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t task.json
```

### Roadmap

- alternative publishing method via telnet
- alternative publishing method via Graphite protocol

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions! 

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik)