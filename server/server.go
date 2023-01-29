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

// Package server contains the implementation of the dynamic server which has a
// contact form rendered (if configured)
package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"

	appconfig "github.com/bossm8/portfoli.go/config"

	"github.com/bossm8/portfoli.go/handler"
	"github.com/bossm8/portfoli.go/messages"
	"github.com/bossm8/portfoli.go/models"
	"github.com/bossm8/portfoli.go/models/config"
	"github.com/bossm8/portfoli.go/models/content"
	"github.com/bossm8/portfoli.go/utils"

	"github.com/microcosm-cc/bluemonday"
)

var (
	cfg         *config.Config
	srvBasePath string
)

// StartServer will attempt to start and listen the server on the specified address
func StartServer(addr string, basePath string, configDir string) {

	var err error
	cfg, err = models.LoadConfiguration(configDir)
	if err != nil && errors.Is(err, config.ErrInvalidSMTPConfig) {
		log.Printf("[WARNING] No smtp configuration loaded, will not render contact form")
		err = nil
	} else if err != nil {
		log.Fatalf("[WARNING] Aborting due to previous error")
	}

	srvBasePath = basePath
	utils.Init(basePath)

	messages.Compile(cfg.Profile.Email.Address)

	fs := http.FileServer(http.Dir(appconfig.StaticContentPath()))

	_http := &handler.RegexHandler{}
	_http.SetBasePath(basePath)

	_http.Handle("/favicon.ico", fs)
	_http.Handle("/static/", http.StripPrefix("/static", fs))
	_http.HandleFunc("/mail", sendMail)
	_http.HandleFunc("/"+messages.RoutingRegexString(), serveStatus)
	_http.HandleFunc("/"+content.GetRoutingRegexString(), serveContent)
	_http.HandleFunc(".*", serveGeneric)

	log.Printf("[INFO] Listening on %s", addr)
	err = http.ListenAndServe(addr, _http)
	if err != nil {
		log.Fatal(err)
	}

}

func serveGeneric(w http.ResponseWriter, r *http.Request) {
	tplName := "index"
	if r.URL.Path != "/" && r.URL.Path != "" {
		tplName = filepath.Base(r.URL.Path)
	}
	sendTemplate(w, r, tplName, nil, nil)
}

func isContentEnabled(requestedContent string) bool {
	enabled := false
	for _, configuredContent := range cfg.Profile.ContentTypes {
		if configuredContent == requestedContent {
			enabled = true
			break
		}
	}
	return enabled
}

func serveContent(w http.ResponseWriter, r *http.Request) {

	contentType := filepath.Base(r.URL.Path)
	if !isContentEnabled(contentType) {
		fail(w, r, messages.MsgNotFound)
		return
	}

	data, err := content.GetRenderedContent(contentType)
	if nil != err {
		fail(w, r, messages.MsgGeneric)
		return
	}

	sendTemplate(w, r, appconfig.ContentTemplateName, data, nil)

}

func sendMail(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/"+appconfig.ContactTemplateName, http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("[ERROR] Could not parse contact form: %s\n", err)
		fail(w, r, messages.MsgContact)
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	form := map[string]string{
		"name":    "",
		"email":   "",
		"message": "",
	}

	// just sanitize everything
	for key := range form {
		form[key] = sanitizer.Sanitize(
			r.FormValue(key),
		)
	}

	var addr *mail.Address
	var err error

	if addr, err = mail.ParseAddress(form["email"]); nil != err {
		log.Println("[ERROR] Received invalid email address for contact form, will not send mail")
		fail(w, r, messages.MsgAddress)
		return
	}

	err = cfg.SMTP.SendMail(
		cfg.Profile.Email.Address,
		addr,
		form["name"],
		form["message"])
	if nil != err {
		fail(w, r, messages.MsgContact)
		return
	}

	log.Printf("[INFO] Sent contact email to %s\n", cfg.Profile.Email)

	// redirect, so form gets cleared and a refresh does not trigger another send
	success(w, r, messages.MsgContact)

}

func serveStatus(w http.ResponseWriter, r *http.Request) {

	vals := r.URL.Query()
	kind := vals.Get("kind")

	status := filepath.Base(r.URL.Path)
	msg := messages.Get(status, kind)

	sendTemplate(w, r, appconfig.StatusTemplateName, msg, &msg.HttpStatus)

}

func abortWithStatusTplCheck(templateName string, w http.ResponseWriter, r *http.Request) {
	if templateName == appconfig.StatusTemplateName {
		// This catches the case when the server cant find or has an error with the status template
		// if this check is not made we end up having an infinite amount of requests
		// because we would again be redirected to the status template
		log.Printf(
			"[WARNING] The template failed is the status template, aborting with %d\n",
			http.StatusInternalServerError,
		)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fail(w, r, messages.MsgGeneric)
	}
}

func sendTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}, status *int) {

	// someone might enter /contact manually - make sure it is not returned if disabled
	if !cfg.RenderContact && templateName == appconfig.ContactTemplateName {
		fail(w, r, messages.MsgNotFound)
		return
	}

	htmlTpl := filepath.Join(appconfig.HTMLTemplatesPath(), templateName+".html")

	if res, err := os.Stat(htmlTpl); os.IsNotExist(err) || res.IsDir() {
		log.Printf("[ERROR] Could not find or read from %s\n", htmlTpl)
		abortWithStatusTplCheck(templateName, w, r)
		return
	}

	tplData := &models.TemplateData{
		Data:          data,
		Profile:       cfg.Profile,
		RenderContact: cfg.RenderContact,
	}

	resp, err := utils.RenderTemplate(appconfig.BaseTemplateName, tplData, appconfig.BaseTemplatePath(), htmlTpl)
	if nil != err {
		abortWithStatusTplCheck(templateName, w, r)
		return
	}

	if nil != status && 100 <= *status {
		w.WriteHeader(*status)
	}
	w.Write(resp)

}

func fail(w http.ResponseWriter, r *http.Request, kind messages.MessageType) {
	http.Redirect(w, r, fmt.Sprintf("%s%s?kind=%s", srvBasePath, messages.EndpointFail, kind), http.StatusSeeOther)
}

func success(w http.ResponseWriter, r *http.Request, kind messages.MessageType) {
	http.Redirect(w, r, fmt.Sprintf("%s%s?kind=%s", srvBasePath, messages.EndpointSuccess, kind), http.StatusSeeOther)
}
