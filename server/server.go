package server

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"bossm8.ch/portfolio/handler"
	"bossm8.ch/portfolio/models"
	"bossm8.ch/portfolio/models/config"
	"bossm8.ch/portfolio/models/content"
	"bossm8.ch/portfolio/server/messages"
	"github.com/microcosm-cc/bluemonday"
)

const (
	templateDir = "templates"
	staticDir   = "static"
)

var (
	baseTpl = filepath.Join(templateDir, "html", "base.html")

	cfg *config.Config
)

func loadConfiguration() {
	var err error
	cfg, err = config.GetConfig()
	if err != nil && errors.Is(err, config.ErrInvalidSMTPConfig) {
		log.Printf("[WARNING]: No smtp configuration loaded, will not render contact form")
	} else if err != nil {
		log.Fatalf("[ERROR] Aborting due to previous error")
	}
	messages.Compile(cfg.Profile.Email.Address)
}

func StartServer(addr string, configDir string) {

	models.SetConfigDir(configDir)
	loadConfiguration()

	fs := http.FileServer(http.Dir(staticDir))

	_http := &handler.RegexHandler{}

	_http.Handle("/favicon.ico", fs)
	_http.Handle("/static/", http.StripPrefix("/static", fs))
	_http.HandleFunc("/mail", sendMail)
	_http.HandleFunc("/"+messages.GetRoutingRegexString(), serveStatus)
	_http.HandleFunc("/"+content.GetRoutingRegexString(), serveContent)
	_http.HandleFunc(".*", serveGeneric)

	err := http.ListenAndServe(addr, _http)
	if err != nil {
		log.Fatal(err)
	}

}

func serveGeneric(w http.ResponseWriter, r *http.Request) {
	htmlFile := "index"
	if r.URL.Path != "/" {
		htmlFile = filepath.Base(r.URL.Path)
	}
	sendTemplate(w, r, htmlFile, nil, nil)
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

	content, err := content.GetRenderedContent(contentType)
	if nil != err {
		fail(w, r, messages.MsgGeneric)
		return
	}

	data := &struct {
		HTML []template.HTML
		Kind string
	}{
		HTML: content,
		Kind: cases.Title(language.English).String(contentType),
	}
	sendTemplate(w, r, "cards", data, nil)

}

func sendMail(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/contact", http.StatusSeeOther)
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

	sendTemplate(w, r, "status", msg, &msg.HttpStatus)

}

func sendTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}, status *int) {

	var tpl *template.Template
	var err error

	// someone might enter /contact manually - make sure it is not returned if disabled
	if !cfg.RenderContact && templateName == "contact" {
		fail(w, r, messages.MsgContact)
		return
	}

	htmlTpl := filepath.Join(templateDir, "html", templateName+".html")

	if res, err := os.Stat(htmlTpl); os.IsNotExist(err) || res.IsDir() {
		fail(w, r, messages.MsgNotFound)
		return
	}

	// Title is used in templates to title case content kind names
	funcMap := template.FuncMap{
		"Title": cases.Title(language.English).String,
	}
	if tpl, err = template.New(htmlTpl).Funcs(funcMap).ParseFiles(baseTpl, htmlTpl); nil != err {
		log.Printf("[ERROR] Failed to parse template: %s with error %s\n", templateName, err)
		fail(w, r, messages.MsgGeneric)
		return
	}

	tplData := &models.TemplateData{
		Data:          data,
		Profile:       cfg.Profile,
		RenderContact: cfg.RenderContact,
	}

	// We cannot pass w to ExecuteTemplate directly
	// if the template fails we cannot redirect because there would be superfluous call to w.WriteHeader
	resp := bytes.Buffer{}
	if err = tpl.ExecuteTemplate(&resp, "base", tplData); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
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
