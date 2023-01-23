package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

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

	mailCfg       *models.MailConfig
	renderContact bool = true

	messages map[string]map[string]*models.AlertMsg
)

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
	flag.Parse()

	var err error
	mailCfg, err = models.GetSMTPConfig()
	if err != nil {
		log.Printf("[WARNING]: No smtp configuration loaded, will not render contact form")
		renderContact = false
		messages = models.GetMessages("")
	} else {
		messages = models.GetMessages(mailCfg.To)
	}

	startServer(fmt.Sprintf("%s:%d", *addr, *port))

}

func startServer(addr string) {

	fs := http.FileServer(http.Dir(staticDir))

	_http := &handler.RegexHandler{}

	_http.Handle("/static/", http.StripPrefix("/static", fs))
	_http.Handle("/favicon.ico", fs)
	_http.HandleFunc("/mail", sendMail)
	_http.HandleFunc("/(success|fail)", serveStatus)
	_http.HandleFunc(".*", serveParamless)

	err := http.ListenAndServe(addr, _http)
	if err != nil {
		log.Fatal(err)
	}

}

func serveParamless(w http.ResponseWriter, r *http.Request) {
	htmlFile := "index"
	if r.URL.Path != "/" {
		htmlFile = filepath.Clean(r.URL.Path)
	}
	sendTemplate(w, r, htmlFile, nil)
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

	for key, _ := range form {
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
		"To":       {mailCfg.To},
		"From":     {mailCfg.SMTP.User, "[Portfolio]: " + form["name"]},
		"Reply-To": {form["email"]},
		"Subject":  {subject},
	})
	mail.SetBody("text/plain", form["message"])

	smtp := mailCfg.SMTP
	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.Pass)
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("[ERROR] Could not send email: %s\n", err)
		fail(w, r, "contact")
		return
	}
	log.Printf("[INFO] Sent contact email to %s\n", mailCfg.To)

	// redirect, so form gets cleared and a refresh does not trigger another send
	http.Redirect(w, r, "/success?kind=contact", http.StatusSeeOther)
}

func serveStatus(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	kind := vals.Get("kind")
	status := strings.Replace(r.URL.Path, "/", "", -1)
	data := &models.TPLData{
		RenderContact: renderContact,
		Data:          messages[status][kind],
	}
	sendTemplate(w, r, "status", data)
}

func sendTemplate(w http.ResponseWriter, r *http.Request, templateName string, data *models.TPLData) {

	var tpl *template.Template
	var err error

	// someone might enter /concat manually - make sure it is not returned if disabled
	if !renderContact && strings.Contains(templateName, "contact") {
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

	if nil == data {
		data = &models.TPLData{
			RenderContact: renderContact,
		}
	}

	// We cannot pass w to ExecuteTemplate directly
	// if the template fails we cannot redirect because there would be superfluous call to w.WriteHeader
	resp := bytes.Buffer{}
	if err = tpl.ExecuteTemplate(&resp, "base", data); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		fail(w, r, "generic")
		return
	}
	w.Write(resp.Bytes())

}

func fail(w http.ResponseWriter, r *http.Request, kind string) {
	http.Redirect(w, r, "/fail?kind="+kind, http.StatusSeeOther)
}
