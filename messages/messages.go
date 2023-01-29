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

// Package messages contains messages which can be displayed on either success
// or failure pages
package messages

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
)

// AlertMsg is the object which can be passed down to the status template
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
				Image:      "delivered.svg",
			},
		},
		EndpointFail: {
			MsgAddress: {
				Title:      "Error",
				Header:     "Oops, something went went wrong",
				Message:    "I could not understand your email address, please try again",
				Kind:       "danger",
				HttpStatus: http.StatusBadRequest,
				Image:      "undelivered.svg",
			},
			MsgContact: {
				Title:      "Error",
				Header:     "Oops, something went wrong",
				Message:    template.HTML("I could not process your contact request, please contact me here: " + mailto),
				Kind:       "warning",
				HttpStatus: http.StatusInternalServerError,
				Image:      "delivered.svg",
			},
			MsgNotFound: {
				Title:      "404",
				Header:     "Oops, something went wrong",
				Message:    "<i class='bi-binoculars me-1'></i> I could not find the page you are looking for <i class='ms-1 bi-binoculars'></i>",
				Kind:       "danger",
				HttpStatus: http.StatusNotFound,
				Image:      "404.svg",
			},
			MsgGeneric: {
				Title:      "Sumthin Wong",
				Header:     "Oops, something went wrong",
				Message:    template.HTML("There was an error on my end, please try again or contact me on " + mailto),
				Kind:       "warning",
				HttpStatus: http.StatusInternalServerError,
				Image:      "error.svg",
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
