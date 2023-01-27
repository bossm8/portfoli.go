package messages

import (
	"fmt"
	"html/template"
	"log"
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

type MessageType string
type MessageEndpoint string

const (
	EndpointSuccess MessageEndpoint = "success"
	EndpointFail    MessageEndpoint = "fail"

	MsgContact  MessageType = "contact"
	MsgAddress  MessageType = "address"
	MsgNotFound MessageType = "notfound"
	MsgGeneric  MessageType = "generic"
)

var messages map[MessageEndpoint]map[MessageType]*AlertMsg

func Compile(emailAddress *mail.Address) {
	mailto := "<a href=\"mailto:%s\">%s</a>"
	if nil == emailAddress {
		mailto = ""
	} else {
		mailto = fmt.Sprintf(mailto, emailAddress.Address, emailAddress.Address)
	}
	messages = map[MessageEndpoint]map[MessageType]*AlertMsg{
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

func Get(endpoint string, kind string) (msg *AlertMsg) {
	var ok bool
	if msg, ok = messages[MessageEndpoint(endpoint)][MessageType(kind)]; !ok {
		log.Printf("[WARNING] Invalid message requested: %s/%s\n", endpoint, kind)
		return messages[EndpointFail][MsgGeneric]
	}
	return
}

func GetRoutingRegexString() string {
	return fmt.Sprintf("(%s|%s)", EndpointSuccess, EndpointFail)
}
