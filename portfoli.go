// Copyright (c) 2023, Boss Marco <bossm8@hotmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// porfoli.go the simple and dynamic portfolio written with Go and Bootstrap
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bossm8/portfoli.go/config"
	"github.com/bossm8/portfoli.go/server"
	"github.com/bossm8/portfoli.go/static"
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
	basePath := flag.String(
		"srv.base",
		"/",
		"The base path to serve content on",
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
	imageCacheDir := flag.String(
		"images.cache.dir",
		"",
		"Path to the directory where cached images are stored (default: img/cache, relative to static dir only)",
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
		static.Build(*basePath, *configDir, *imageCacheDir)
	} else {
		// Do not log the dist dir path by using nil
		config.SetPaths(templatesDir, staticDir, nil)
		server.StartServer(fmt.Sprintf("%s:%d", *addr, *port), *basePath, *configDir, *imageCacheDir)
	}

}
