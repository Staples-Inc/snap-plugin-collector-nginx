go install ./...
snapctl plugin load $GOPATH/bin/snap-plugin-collector-nginx
snapctl plugin load $GOPATH/bin/snap-plugin-processor-staplesfmt
snapctl plugin load $GOPATH/bin/snap-plugin-publisher-metrics
snapctl task create -t $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-nginx/examples/tasks/task.json
