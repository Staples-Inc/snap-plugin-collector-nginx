
## Nginx Collector Plugin Structure
---
Nginx plugins  

#### Plugin binary

./main.go

##### Collector Implementation

./nginx/nginx.go

##### Testing #########
Run below command in seperate terminal
Run the snap agent

> $GOPATH/bin/snap-v0.14.0-beta/bin/snapd --plugin-trust 0 --log-level 1

Run the snap agent with the config file

> $GOPATH/bin/snap-v0.14.0-beta/bin/snapd --plugin-trust 0 --log-level 1 --config $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-nginx/config.json

Run the collector plugin seperately

> $GOPATH/bin/snap-v0.14.0-beta/bin/snapctl  plugin load $GOPATH/bin/snap-plugin-collector-nginx 

Or Run the collector using run.sh

To verify metrics

> $GOPATH/bin/snap-v0.14.0-beta/bin/snapctl metric list


