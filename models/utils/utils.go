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

package utils

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// The directory where all (static and dynamic) configuration files are read from
var yamlDir string

// SetYAMLDir sets the configuration directory where dynamic and static configurations
// should be read from
func SetYAMLDir(dir string) {
	yamlDir = dir
}

// LoadFromYAMLFile loads the file with filename into obj
func LoadFromYAMLFile(filename string, obj interface{}) (err error) {
	// log.Printf("[INFO] Loading yaml file '%s' from directory '%s'\n", filename, configDir)
	var yamlFile []byte
	if yamlFile, err = os.ReadFile(
		filepath.Join(yamlDir, filename),
	); err != nil {
		log.Printf("[ERROR]: Failed to load yaml file '%s': %s\n", filename, err)
		return
	}
	if err = yaml.Unmarshal(yamlFile, obj); err != nil {
		log.Printf("[ERROR]: Failed to parse yaml file '%s': %s\n", filename, err)
	}
	return
}
