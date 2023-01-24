package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/gomail.v2"

	"bossm8.ch/portfolio/handler"
	"bossm8.ch/portfolio/models"
	"github.com/microcosm-cc/bluemonday"
)

const (
	templateDir = "templates"
	staticDir   = "static"
	mailSubject = "[Portfolio] New message from %s"
)

var (
	baseTpl = filepath.Join("templates", "html", "base.html")

	cfg *models.Config

	renderContact bool = true

	messages map[string]map[string]*AlertMsg
)

func loadConfiguration() {
	var err error
	cfg, err = models.GetConfig()
	if err != nil {
		log.Printf("[WARNING]: No smtp configuration loaded, will not render contact form")
		renderContact = false
	}
	messages = getMessages(&cfg.Profile.Email)
}

func StartServer(addr string, configDir string) {

	models.SetConfigDir(configDir)
	loadConfiguration()

	fs := http.FileServer(http.Dir(staticDir))

	_http := &handler.RegexHandler{}

	_http.Handle("/favicon.ico", fs)
	_http.Handle("/static/", http.StripPrefix("/static", fs))
	_http.HandleFunc("/mail", sendMail)
	_http.HandleFunc("/(success|fail)", serveStatus)
	_http.HandleFunc(models.GetRoutingRegex(), serveWithContent)
	_http.HandleFunc(".*", serveParamless)

	err := http.ListenAndServe(addr, _http)
	if err != nil {
		log.Fatal(err)
	}

}

func serveParamless(w http.ResponseWriter, r *http.Request) {
	htmlFile := "index"
	if r.URL.Path != "/" {
		htmlFile = filepath.Base(r.URL.Path)
	}
	sendTemplate(w, r, htmlFile, nil, nil)
}

func serveWithContent(w http.ResponseWriter, r *http.Request) {
	// TODO: allow disabling of content types (check yaml content key)
	tplName := filepath.Base(r.URL.Path)
	content, err := models.GetContent(tplName)
	if nil != err {
		fail(w, r, "generic")
	}
	data := &struct {
		HTML []template.HTML
		Kind string
	}{
		HTML: content,
		Kind: cases.Title(language.English).String(tplName),
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
		fail(w, r, "contact")
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	form := map[string]string{
		"name":    "",
		"email":   "",
		"message": "",
	}

	for key := range form {
		form[key] = sanitizer.Sanitize(
			r.FormValue(key),
		)
	}

	if _, err := mail.ParseAddress(form["email"]); nil != err {
		log.Println("[ERROR] Received invalid email address for contact form, will not send mail")
		fail(w, r, "address")
		return
	}

	mail := gomail.NewMessage()
	subject := fmt.Sprintf(mailSubject, form["name"])
	mail.SetHeaders(map[string][]string{
		"To":       {cfg.Profile.Email},
		"From":     {cfg.SMTP.User, "[Portfolio]: " + form["name"]},
		"Reply-To": {form["email"]},
		"Subject":  {subject},
	})
	mail.SetBody("text/plain", form["message"])

	smtp := cfg.SMTP
	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.Pass)
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("[ERROR] Could not send email: %s\n", err)
		fail(w, r, "contact")
		return
	}
	log.Printf("[INFO] Sent contact email to %s\n", cfg.Profile.Email)

	// redirect, so form gets cleared and a refresh does not trigger another send
	http.Redirect(w, r, "/success?kind=contact", http.StatusSeeOther)
}

func serveStatus(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	kind := vals.Get("kind")
	status := strings.Replace(r.URL.Path, "/", "", -1)
	msg := messages[status][kind]
	sendTemplate(w, r, "status", msg, &msg.HttpStatus)
}

func sendTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}, status *int) {

	var tpl *template.Template
	var err error

	// someone might enter /contact manually - make sure it is not returned if disabled
	if !renderContact && templateName == "contact" {
		fail(w, r, "notfound")
		return
	}

	htmlTpl := filepath.Join(templateDir, "html", templateName+".html")

	if res, err := os.Stat(htmlTpl); os.IsNotExist(err) || res.IsDir() {
		fail(w, r, "notfound")
		return
	}

	if tpl, err = template.ParseFiles(baseTpl, htmlTpl); nil != err {
		log.Printf("[ERROR] Failed to parse template: %s with error %s\n", templateName, err)
		fail(w, r, "generic")
		return
	}

	tplData := &models.TemplateData{
		Data:          data,
		RenderContact: renderContact,
		Profile:       cfg.Profile,
	}

	// We cannot pass w to ExecuteTemplate directly
	// if the template fails we cannot redirect because there would be superfluous call to w.WriteHeader
	resp := bytes.Buffer{}
	if err = tpl.ExecuteTemplate(&resp, "base", tplData); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		fail(w, r, "generic")
		return
	}
	if nil != status && 100 <= *status {
		w.WriteHeader(*status)
	}
	w.Write(resp.Bytes())

}

func fail(w http.ResponseWriter, r *http.Request, kind string) {
	http.Redirect(w, r, "/fail?kind="+kind, http.StatusSeeOther)
}