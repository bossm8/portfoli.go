// porfoli.go the simple and dynamic portfolio written with Go and Bootstrap
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bossm8.ch/portfolio/config"
	"bossm8.ch/portfolio/server"
	"bossm8.ch/portfolio/static"
)

func main() {

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf(
			"[ERROR] Could not find the  user config dir, please provide the path to the configuration files manually (%s)",
			err,
		)
	}
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
		filepath.Join(cfgDir, "portfoli.go", "configs"),
		"Path to the directory containing the yaml configurations",
	)
	staticDir := flag.String(
		"static.dir",
		config.StaticContentPath(),
		"Path to the directory containing the static content for the website",
	)
	templatesDir := flag.String(
		"templates.dir",
		config.TemplatesPath(),
		"Path to the directory containing the html templates",
	)
	distDir := flag.String(
		"dist.dir",
		config.DistDir(),
		"Path to the directory of where to output the static build",
	)
	verbose := flag.Bool(
		"verbose",
		false,
		"Print more verbose logging information (filenames)",
	)
	dist := flag.Bool(
		"dist",
		false,
		"Create a static website build to e.g. host on GitLab pages",
	)
	flag.Parse()

	if *verbose {
		log.SetFlags(log.Lshortfile)
	}

	*configDir = config.ConvertToAbsPath(configDir)
	log.Printf("[INFO] using config path %s\n", *configDir)

	if *dist {
		config.SetPaths(templatesDir, staticDir, distDir)
		static.Build(*configDir)
	} else {
		config.SetPaths(templatesDir, staticDir, nil)
		server.StartServer(fmt.Sprintf("%s:%d", *addr, *port), *configDir)
	}

}
