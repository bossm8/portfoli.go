package server

import (
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
)

type AlertMsg struct {
	Title      string
	Image      string
	Header     string
	Message    template.HTML
	Kind       string
	HttpStatus int
}

const (
	EndpointSuccess = "success"
	EndpointFail    = "fail"

	MsgContact  = "contact"
	MsgAddress  = "address"
	MsgNotFound = "notfound"
	MsgGeneric  = "generic"
)

func getMessages(emailAddress *mail.Address) map[string]map[string]*AlertMsg {
	mailto := "<a href=\"mailto:%s\">%s</a>"
	if nil == emailAddress {
		mailto = ""
	} else {
		mailto = fmt.Sprintf(mailto, *emailAddress, *emailAddress)
	}
	return map[string]map[string]*AlertMsg{
		EndpointSuccess: {
			MsgContact: {
				Title:      "Success",
				Header:     "Message sent successfully",
				Message:    "I will get in touch with you shortly",
				Kind:       "success",
				HttpStatus: http.StatusOK,
			},
		},
		EndpointFail: {
			MsgAddress: {
				Title:      "Error",
				Header:     "Oops, something went went wrong",
				Message:    "I could not understand your email address, please try again",
				Kind:       "danger",
				HttpStatus: http.StatusBadRequest,
			},
			MsgContact: {
				Title:      "Error",
				Header:     "Oops, something went wrong",
				Message:    template.HTML("I could not process your contact request, please contact me here: " + mailto),
				Kind:       "warning",
				HttpStatus: http.StatusInternalServerError,
			},
			MsgNotFound: {
				Title:      "404",
				Header:     "Oops, something went wrong",
				Message:    "<i class='bi-binoculars me-1'></i> I could not find the page you are looking for <i class='ms-1 bi-binoculars'></i>",
				Kind:       "danger",
				HttpStatus: http.StatusNotFound,
			},
			MsgGeneric: {
				Title:      "Sumthin Wong",
				Header:     "Oops, something went wrong",
				Message:    template.HTML("There was an error on my end, please try again or contact me on " + mailto),
				Kind:       "warning",
				HttpStatus: http.StatusInternalServerError,
			},
		},
	}
}
