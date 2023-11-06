package main

import (
	"github.com/xdblab/xdb-golang-samples/cmd/server/xdb"
	"os"
)

// main entry point for the iwf server
func main() {
	app := xdb.BuildCLI()
	app.Run(os.Args)
}
