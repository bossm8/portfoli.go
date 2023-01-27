// porfoli.go the simple and dynamic portfolio written with Go and Bootstrap
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bossm8.ch/portfolio/server"
	"bossm8.ch/portfolio/static"
)

// convertToAbsPath takes the path of the directory
// containing the configuration (either relative or absolute) and converts it
// to an absolute path
func convertToAbsPath(path *string) {
	if filepath.IsAbs(*path) {
		return
	}
	if exe, err := os.Executable(); err != nil {
		log.Fatalf("[ERROR] parsing configuration directory: %s\n", err)
	} else {
		*path = filepath.Join(
			filepath.Dir(exe),
			*path,
		)
	}
}

func main() {

	addr := flag.String(
		"srv.address",
		"127.0.0.1",
		"Listen address for the protfolio server",
	)
	port := flag.Int(
		"srv.port",
		8080,
		"Listen port for the portfolio server",
	)
	configDir := flag.String(
		"config.dir",
		"configs",
		"Path to the directory containing the yaml configurations",
	)
	verbose := flag.Bool(
		"verbose",
		false,
		"Print more verbose logging information (filenames)",
	)
	dist := flag.Bool(
		"static",
		false,
		"Create a static website build to e.g. host on GitLab pages",
	)
	flag.Parse()

	if *verbose {
		log.SetFlags(log.Lshortfile)
	}

	convertToAbsPath(configDir)
	log.Printf("[INFO] using config path %s\n", *configDir)

	if *dist {
		static.Build(*configDir)
	} else {
		server.StartServer(fmt.Sprintf("%s:%d", *addr, *port), *configDir)
	}

}
