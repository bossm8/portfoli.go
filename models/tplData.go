package models

import (
	"fmt"
	"html/template"
	"net/http"
)

type TPLData struct {
	RenderContact bool
	Data          interface{}
	Profile       *ProfileConfig
}

type AlertMsg struct {
	Title      string
	Image      string
	Header     string
	Message    template.HTML
	Kind       string
	HttpStatus int
}

func GetMessages(emailAddress *string) map[string]map[string]*AlertMsg {
	mailto := "<a href=\"mailto:%s\">%s</a>"
	if nil == emailAddress {
		mailto = ""
	} else {
		mailto = fmt.Sprintf(mailto, *emailAddress, *emailAddress)
	}
	return map[string]map[string]*AlertMsg{
		"success": {
			"contact": {
				Title:      "Success",
				Header:     "Message sent successfully",
				Message:    "I will get in touch with you shortly",
				Kind:       "success",
				HttpStatus: http.StatusOK,
			},
		},
		"fail": {
			"address": {
				Title:      "Error",
				Header:     "Oops, something went went wrong",
				Message:    "I could not understand your email address, please try again",
				Kind:       "danger",
				HttpStatus: http.StatusBadRequest,
			},
			"contact": {
				Title:      "Error",
				Header:     "Oops, something went wrong",
				Message:    template.HTML("I could not process your contact request, please contact me here: " + mailto),
				Kind:       "warning",
				HttpStatus: http.StatusInternalServerError,
			},
			"notfound": {
				Title:      "404",
				Header:     "Oops, something went wrong",
				Message:    "I could not find the page you are looking for",
				Kind:       "danger",
				HttpStatus: http.StatusNotFound,
			},
			"generic": {
				Title:      "Sumthin Wong",
				Header:     "Oops, something went wrong",
				Message:    template.HTML("There was an error on my end, please try again or contact me on " + mailto),
				Kind:       "warning",
				HttpStatus: http.StatusInternalServerError,
			},
		},
	}
}
