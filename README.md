# Snap Collector plugin - Nginx
Collector get Nginx metrics from Nginx application and pass it to blueflood and Metric square to store it in cassandra

I. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating Systems](#operating-systems)
  * [Build](#build)
  * [Run](#run)
  * [Verify](#verify)

II. [Contributing](#contributing)

III. [License](#license)

## Getting Started
A working snap agent and a running instance of Nignx application which expose a rest api in json format to get access to real time nignx metrics.

### System Requirements
* [golang 1.5+](https://golang.org/dl/)
* [snap](https://github.com/intelsdi-x/snap)
* [blueflood](http://blueflood.io/)
* [metric](htts://github.com/square/metrics)
* [cassandra](http://cassandra.apache.org/)
* [nignx](http://demo.nginx.com/status)

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
Run the snap agent with the config file

> $GOPATH/bin/snap-v0.14.0-beta/bin/snapd --plugin-trust 0 --log-level 1 --config $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-nginx/config.json

Run the collector plugin seperately

> $GOPATH/bin/snap-v0.14.0-beta/bin/snapctl  plugin load $GOPATH/bin/snap-plugin-collector-nginx

### Verify
To Verify nginx mertics
> $GOPATH/bin/snap-v0.14.0-beta/bin/snapctl metric list

## Contributing
We currently have no future plans for this plugin. If you have a feature request, please add it as an issue and/or submit a pull request

## License
This plugin is Open Source software released uder the Apache 2.0 [License](LICENSE)
