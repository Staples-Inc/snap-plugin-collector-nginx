# Snap Collector plugin - NginxPlus
Collector get nginx metrics from the nginx-plus endpoint.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating Systems](#operating-systems)
  * [Build](#build)
  * [Run](#run)
  * [Verify](#verify)
  * [Configuration and Usage](#configuration-and-usage)
2. [Community Support](#community-support)
3. [Contributing](#contributing)
4. [License](#license)

## Getting Started
A working snap agent and a running instance of nginx application which expose a rest api in json format to get access to real time nignx metrics.

### System Requirements
* [golang 1.5+](https://golang.org/dl/) - needed only for building
* [snap](https://github.com/intelsdi-x/snap)
* [nginx](https://www.nginx.com/)
* [nginx status module](http://nginx.org/en/docs/http/ngx_http_status_module.html)

### Operating Systems
* Linux
* Mac OS X

### Build
Fork https://github.com/Staples-Inc/snap-plugin-collector-nginx
Clone repo into `$GOPATH/src/github.com/Staples-Inc/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-nginx.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Run
Make sure that your $SNAP_PATH is set, e.g.:

$ export SNAP_PATH=\<snapDirectoryPath\>/build/linux/x86_64

Run the snap agent with the config file:
> $SNAP_PATH/snapteld --plugin-trust 0 --log-level 1 --config $GOPATH/src/github.com/Staples-Inc/snap-plugin-collector-nginx/config.json

Load the collector plugin:
> $SNAP_PATH/snaptel plugin load $GOPATH/src/github.com/Staples-Inc/snap-plugin-collector-nginx/build/rootfs/snap-plugin-collector-nginx

### Verify
To Verify nginx mertics:
> $SNAP_PATH/snaptel metric list

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
As part of snapteld global config
* Load the plugin and create a task, you can find example config and task manifest files in "examples" directory

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. The full project is at http://github.com:intelsdi-x/snap.
To reach out on other use cases, visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We currently have no future plans for this plugin. If you have a feature request, please add it as an issue and/or submit a pull request

## License
This plugin is Open Source software released under the Apache 2.0 [License](LICENSE)
