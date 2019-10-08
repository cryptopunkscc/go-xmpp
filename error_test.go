package xmpp

import (
	"encoding/xml"
	"testing"
)

const error1 = "<stream:error><system-shutdown xmlns='urn:ietf:params:xml:ns:xmpp-streams'/><text xmlns='urn:ietf:params:xml:ns:xmpp-streams' xml:lang='en'>description</text></stream:error>"
const error2 = "<stream:error><system-shutdown xmlns='urn:ietf:params:xml:ns:xmpp-streams'/><text xmlns='urn:ietf:params:xml:ns:xmpp-streams' xml:lang='en'>description</text><extra/></stream:error>"
const error3 = "<stream:error><system-shutdown xmlns='urn:ietf:params:xml:ns:xmpp-streams'/><extra/></stream:error>"

func TestError_UnmarshalXML(t *testing.T) {
	// condition + text
	var err1 = &Error{}
	if e := xml.Unmarshal([]byte(error1), err1); e != nil {
		t.Error(e)
	}
	if err1.Condition != "system-shutdown" {
		t.Error("invalid condition")
	}
	if err1.Text != "description" {
		t.Error("invalid text")
	}
	if err1.Extra != nil {
		t.Error("extra should be nil")
	}

	// condition + text + extra
	var err2 = &Error{}
	if e := xml.Unmarshal([]byte(error2), err2); e != nil {
		t.Error(e)
	}
	if err2.Condition != "system-shutdown" {
		t.Error("invalid condition")
	}
	if err2.Text != "description" {
		t.Error("invalid text")
	}
	if err2.Extra == nil {
		t.Error("missing extra")
	}

	// condition + extra
	var err3 = &Error{}
	if e := xml.Unmarshal([]byte(error3), err3); e != nil {
		t.Error(e)
	}
	if err3.Condition != "system-shutdown" {
		t.Error("invalid condition")
	}
	if err3.Text != "" {
		t.Error("invalid text:", err1.Text)
	}
	if err3.Extra == nil {
		t.Error("missing extra")
	}
}

func TestError_MarshalXML(t *testing.T) {
	orig := &Error{
		Condition: "system-shutdown",
		Text:      "description",
		Extra:     nil,
	}

	bytes, err := xml.Marshal(orig)
	if err != nil {
		t.Error(err)
	}

	clone := &Error{}
	if err := xml.Unmarshal(bytes, clone); err != nil {
		t.Error(err)
	}

	if clone.Condition != orig.Condition {
		t.Error("condition mismatch")
	}
	if clone.Text != orig.Text {
		t.Error("text mismatch")
	}
}
