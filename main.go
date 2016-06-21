
package main

import (
   "os"

   "github.com/intelsdi-x/snap/control/plugin"
   "github.com/intelsdi-x/snap-plugin-collector-jolokia/jolokia"
)

func main() {
    meta := jolokia.Meta()
    plugin.Start(meta, new(jolokia.Jolokia), os.Args[1])
}
