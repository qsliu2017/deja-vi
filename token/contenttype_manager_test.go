package token

import "testing"

func TestInit(t *testing.T) {
	if TheContentTypeManager == nil {
		t.Error("TheContentTypeManager is nil")
	}

	if TheTokenizer == nil {
		t.Error("TheTokenizer is nil")
	}
}

func TestSetContentType(t *testing.T) {
	isTxtTokenizer(t, TheTokenizer)

	TheContentTypeManager.SetContentType("xml")
	isXmlTokenizer(t, TheTokenizer)

	TheContentTypeManager.SetContentType("txt")
	isTxtTokenizer(t, TheTokenizer)
}
