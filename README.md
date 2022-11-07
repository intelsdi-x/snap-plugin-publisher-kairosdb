DISCONTINUATION OF PROJECT. 

This project will no longer be maintained by Intel.

This project has been identified as having known security escapes.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.

# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**




# Snap publisher plugin - KairosDB

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

* [golang 1.6+](https://golang.org/dl/) for building plugin from source code

Support Matrix

- KairosDB V1 REST API

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation

#### Download the plugin binary:
You can get the pre-built binaries for your OS and architecture at plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/releases) page.

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
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/blob/master/README.md#examples).

## Documentation

For details on KairosDB, please refer to [documentation](https://kairosdb.github.io/docs/build/html/index.html).

### Task manifest
Task manifest configuration is described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md). In order to use KairosDB publisher you have to add section "publish" then specify following options:
- `"host"` - KairosDB host address (ex. `"127.0.0.1"`)
- `"port"` -  KairosDB REST API port (ex. `"8083"`)
- `"useDynamic"` -  Use dynamic elements of namespace as tags for metric (ex. `true`)

To know more about Dynamic Metric Namespace, please refer to [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/METRICS.md#dynamic-metrics)

See example task manifest in [examples/tasks/](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/blob/master/examples/tasks/).

### Examples
Example of running [psutil collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-psutil) and publishing data to KairosDB.

**Sample** *`useDynamic=false`*

| metric | tags | value |
|--------------------------|-------------------------|----|
| /intel/psutil/cpu/0/user	| plugin_running_on=host1 | 12 |
| /intel/psutil/cpu/1/user | plugin_running_on=host1 | 42 |


**Sample** *`useDynamic=true`*

| metric | tags | value |
|------------------------|---------------------------------------------|----|
| /intel/psutil/cpu/user | host=host1 cpu_id=0 plugin_running_on=host1 | 12 |
| /intel/psutil/cpu/user | host=host1 cpu_id=1 plugin_running_on=host1 | 42 |


Another example of running [psutil collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-psutil) and publishing data to KairosDB.

Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`

Download and load Snap plugins (paths to binary files for Linux/amd64):
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-kairosdb/latest/linux/x86_64/snap-plugin-publisher-kairosdb
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-psutil/latest/linux/x86_64/snap-plugin-collector-psutil
$ snaptel plugin load snap-plugin-publisher-kairosdb
$ snaptel plugin load snap-plugin-collector-psutil
```

Create a [task manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) (see [exemplary tasks](examples/tasks/)),
for example `psutil-kairosdb.json` with following content:
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
        "/intel/psutil/load/load1": {},
        "/intel/psutil/load/load15": {}
      },
      "publish": [
        {
          "plugin_name": "kairos",
          "config": {
           "host": "127.0.0.1",
           "port": 8080,
           "useDynamic": true
          }
        }
      ]
    }
  }
}

```
Create a task:
```
$ snaptel task create -t psutil-kairosdb.json
```

Watch created task:
```
$ snaptel task watch <task_id>
```

To stop previously created task:
```
$ snaptel task stop <task_id>
```

### Roadmap
- alternative publishing method via telnet
- alternative publishing method via Graphite protocol

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions! 

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik)
