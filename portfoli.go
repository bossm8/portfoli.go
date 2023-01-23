package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bossm8.ch/portfolio/server"
)

func getConfDirAbsPath(flag string) (dir string) {
	if filepath.IsAbs(flag) {
		return flag
	}
	if exe, err := os.Executable(); err != nil {
		log.Fatalf("[ERROR] parsing configuration directory: %s\n", err)
	} else {
		dir = filepath.Join(
			filepath.Dir(exe),
			flag,
		)
	}
	return
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
	flag.Parse()

	*configDir = getConfDirAbsPath(*configDir)

	log.SetFlags(log.Lshortfile)

	server.StartServer(fmt.Sprintf("%s:%d", *addr, *port), *configDir)

}
