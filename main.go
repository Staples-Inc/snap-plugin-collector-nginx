
package main

import (
"os"

"github.com/intelsdi-x/snap/control/plugin"
"github.com/intelsdi-x/snap-plugin-collector-nginx/nginx"
)

func main() {
	meta := nginx.Meta()
	plugin.Start(meta, new(nginx.Nginx), os.Args[1])
}
