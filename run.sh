go install ./...
snapctl plugin load ~/work/bin/snap-plugin-collector-jolokia
snapctl plugin load ~/work/bin/snap-processor-passthru
snapctl plugin load ~/work/bin/snap-publisher-file
snapctl task create -t ~/work/src/github.com/intelsdi-x/snap-plugin-collector-jolokia/examples/tasks/task.json
