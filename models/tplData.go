package models

import (
	"fmt"
	"html/template"
)

type TPLData struct {
	RenderContact bool
	Data          interface{}
}

type AlertMsg struct {
	Title   string
	Image   string
	Header  string
	Message template.HTML
	Kind    string
}

func GetMessages(emailAddress string) map[string]map[string]*AlertMsg {
	return map[string]map[string]*AlertMsg{
		"success": {
			"contact": {
				Title:   "Success",
				Header:  "Message sent successfully",
				Message: "I will get in touch with you shortly",
				Kind:    "success",
			},
		},
		"fail": {
			"address": {
				Title:   "Error",
				Header:  "Oops, something went went wrong",
				Message: "I could not understand your email address, please try again",
				Kind:    "danger",
			},
			"contact": {
				Title:  "Error",
				Header: "Oops, something went wrong",
				Message: template.HTML(fmt.Sprintf(
					"I could not process your contact request, please contact me here: <a href=\"mailto:%s\">%s</a>",
					emailAddress,
					emailAddress,
				)),
				Kind: "warning",
			},
			"notfound": {
				Title:   "404",
				Header:  "Oops, something went wrong",
				Message: "I could not find the page you are looking for",
				Kind:    "danger",
			},
			"generic": {
				Title:  "Sumthin wong",
				Header: "Oops, something went wrong",
				Message: template.HTML(fmt.Sprintf(
					"There was an error on my end, please try again or contact me on <a href=\"mailto:%s\">%s</a>",
					emailAddress,
					emailAddress,
				)),
				Kind: "warning",
			},
		},
	}
}
