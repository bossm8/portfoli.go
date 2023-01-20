package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

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

var successMessages = map[string]*models.SuccessMsg{
	"contact": {
		Header:  "Message sent successfully",
		Message: "I will get in touch with you shortly",
	},
}

var baseTpl = filepath.Join("templates", "html", "base.html")

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

	exp := models.GetExperiences()
	fmt.Println(exp)

	startServer(fmt.Sprintf("%s:%d", *addr, *port))

}

func startServer(addr string) {

	fs := http.FileServer(http.Dir(staticDir))

	_http := &handler.RegexHandler{}

	_http.Handle("/static/", http.StripPrefix("/static", fs))
	_http.Handle("/favicon.ico", fs)
	_http.HandleFunc("/mail", sendMail)
	_http.HandleFunc("/success", serveSuccess)
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
	sendTemplate(w, htmlFile, nil)
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/contact", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	sanitizer := bluemonday.UGCPolicy()

	name := sanitizer.Sanitize(
		r.FormValue("name"),
	)
	email := sanitizer.Sanitize(
		r.FormValue("email"),
	)
	body := sanitizer.Sanitize(
		r.FormValue("message"),
	)
	subject := fmt.Sprintf(mailSubject, name)

	mail := gomail.NewMessage()
	mail.SetHeaders(map[string][]string{
		"To":       {email},
		"From":     {"asgard.mcathome@outlook.com", "[Portfolio]: " + name},
		"Reply-To": {email},
		"Subject":  {subject},
	})
	mail.SetBody("text/plain", body)

	dialer := gomail.NewDialer("smtp.office365.com", 587, "asgard.mcathome@outlook.com", "jiWbvXbJvVccL3H")
	if err := dialer.DialAndSend(mail); err != nil {
		log.Fatal(err)
	}

	// redirect, so form gets cleared and a refresh does not trigger another send
	http.Redirect(w, r, "/success?kind=contact", http.StatusSeeOther)
}

func serveSuccess(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	kind := vals.Get("kind")
	sendTemplate(w, "success", successMessages[kind])
}

func sendTemplate(w http.ResponseWriter, templateName string, data any) {

	var tpl *template.Template
	var err error

	htmlTpl := filepath.Join(templateDir, "html", templateName+".html")

	res, err := os.Stat(htmlTpl)

	if os.IsNotExist(err) {

	}

	if res.IsDir() {

	}

	if tpl, err = template.ParseFiles(baseTpl, htmlTpl); err != nil {
		log.Fatal(err)
	}
	if err = tpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Fatal(err)
	}

}
