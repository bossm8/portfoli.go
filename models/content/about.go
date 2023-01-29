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

package content

import (
	"html/template"
	"log"
	"path/filepath"

	"github.com/bossm8/portfoli.go/config"
	apputils "github.com/bossm8/portfoli.go/utils"
)

const (
	meTpl = "about.html"
)

type AboutMeConfig struct {
	AboutMe template.HTML `yaml:"me"`
}

// Make sure the interface is implemented
var _ ContentConfig = &AboutMeConfig{}

func (a *AboutMeConfig) ConfigName() string {
	return a.ContentType() + ".yml"
}

func (a *AboutMeConfig) ContentType() string {
	return ContentTypes[typeAbout]
}

func (a *AboutMeConfig) Title() string {
	return "Bio'n'Skills"
}

func (a *AboutMeConfig) Render() (*template.HTML, error) {
	baseTpl := filepath.Join(config.ContentTemplatesPath(), meTpl)
	result, err := apputils.RenderTemplate("about", a.AboutMe, baseTpl)
	if err != nil {
		log.Printf("[ERROR] Failed to render %s\n", baseTpl)
		return nil, err
	}
	html := template.HTML(result)
	return &html, nil
}
