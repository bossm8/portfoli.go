package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"

	"bossm8.ch/portfolio/handler"
	"bossm8.ch/portfolio/messages"
	"bossm8.ch/portfolio/models"
	"bossm8.ch/portfolio/models/config"
	"bossm8.ch/portfolio/models/content"
	"bossm8.ch/portfolio/utils"
	"github.com/microcosm-cc/bluemonday"
)

var (
	cfg *config.Config
)

func StartServer(addr string, configDir string) {

	var err error
	cfg, err = utils.LoadConfiguration(configDir)
	if err != nil && errors.Is(err, config.ErrInvalidSMTPConfig) {
		log.Printf("[WARNING]: No smtp configuration loaded, will not render contact form")
		err = nil
	} else if err != nil {
		log.Fatalf("[ERROR] Aborting due to previous error")
	}
	messages.Compile(cfg.Profile.Email.Address)

	fs := http.FileServer(http.Dir(cfg.StaticDir))

	_http := &handler.RegexHandler{}

	_http.Handle("/favicon.ico", fs)
	_http.Handle("/static/", http.StripPrefix("/static", fs))
	_http.HandleFunc("/mail", sendMail)
	_http.HandleFunc("/"+messages.GetRoutingRegexString(), serveStatus)
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
	if r.URL.Path != "/" {
		tplName = filepath.Base(r.URL.Path)
	}
	sendTemplate(w, r, tplName, nil, nil)
}

func isContentEnabled(requestedContent string) bool {
	enabled := false
	for _, configuredContent := range cfg.Profile.ContentKinds {
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

	sendTemplate(w, r, cfg.ContentTemplateName, data, nil)

}

func sendMail(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/"+cfg.ContactTemplateName, http.StatusSeeOther)
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

	sendTemplate(w, r, cfg.StatusTemplateName, msg, &msg.HttpStatus)

}

func sendTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}, status *int) {

	// someone might enter /contact manually - make sure it is not returned if disabled
	if !cfg.RenderContact && templateName == cfg.ContactTemplateName {
		fail(w, r, messages.MsgContact)
		return
	}

	htmlTpl := filepath.Join(cfg.HTMLTeplatesDir, templateName+".html")

	if res, err := os.Stat(htmlTpl); os.IsNotExist(err) || res.IsDir() {
		log.Printf("[ERROR] Could not find or read from %s\n", htmlTpl)
		fail(w, r, messages.MsgNotFound)
		return
	}

	tplData := &models.TemplateData{
		Data:          data,
		Profile:       cfg.Profile,
		RenderContact: cfg.RenderContact,
	}

	resp, err := utils.RenderTemplate(cfg.BaseTemplateName, cfg.BaseTemplatePath, htmlTpl, tplData)
	if nil != err {
		fail(w, r, messages.MsgGeneric)
		return
	}

	if nil != status && 100 <= *status {
		w.WriteHeader(*status)
	}
	w.Write(resp.Bytes())

}

func fail(w http.ResponseWriter, r *http.Request, kind messages.MessageType) {
	http.Redirect(w, r, fmt.Sprintf("/%s?kind=%s", messages.EndpointFail, kind), http.StatusSeeOther)
}

func success(w http.ResponseWriter, r *http.Request, kind messages.MessageType) {
	http.Redirect(w, r, fmt.Sprintf("/%s?kind=%s", messages.EndpointSuccess, kind), http.StatusSeeOther)
}
