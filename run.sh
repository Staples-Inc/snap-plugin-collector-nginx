go install ./...
$GOPATH/bin/snap-v0.14.0-beta/bin/snapctl plugin load $GOPATH/bin/snap-plugin-collector-nginx
$GOPATH/bin/snap-v0.14.0-beta/bin/snapctl plugin load $GOPATH/bin/snap-plugin-processor-staplesfmt
$GOPATH/bin/snap-v0.14.0-beta/bin/snapctl plugin load $GOPATH/bin/snap-plugin-publisher-metrics
$GOPATH/bin/snap-v0.14.0-beta/bin/snapctl task create -t $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-nginx/examples/tasks/task.json
