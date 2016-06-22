go install ./...
~/work/bin/snap-v0.14.0-beta/bin/snapctl plugin load ~/work/bin/snap-plugin-collector-jolokia
~/work/bin/snap-v0.14.0-beta/bin/snapctl plugin load ~/work/bin/snap-v0.14.0-beta/plugin/snap-processor-passthru
~/work/bin/snap-v0.14.0-beta/bin/snapctl plugin load ~/work/bin/snap-v0.14.0-beta/plugin/snap-publisher-file
~/work/bin/snap-v0.14.0-beta/bin/snapctl task create -t ~/work/src/github.com/intelsdi-x/snap-plugin-collector-jolokia/examples/tasks/task.json
