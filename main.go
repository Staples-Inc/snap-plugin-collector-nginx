package main

import (
	"os"

	"github.com/intelsdi-x/snap-plugin-collector-nginx/nginx"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {
	meta := nginx.Meta()
	plugin.Start(meta, nginx.NewNginxCollector(), os.Args[1])
}
