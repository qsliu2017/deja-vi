package token

import "testing"

func TestTxtTokenizer(t *testing.T) {
	tokenizer := &txtTokenizer{}
	isTxtTokenizer(t, tokenizer)
}

func TestXmlTokenizer(t *testing.T) {
	tokenizer := &xmlTokenizer{}
	isXmlTokenizer(t, tokenizer)
}

func isTxtTokenizer(t *testing.T, tokenizer Tokenizer) {
	actual := tokenizer.Tokenize("What is your point of view, and, ")
	expect := []struct {
		isWord  bool
		content string
	}{
		{true, "What"},
		{false, " "},
		{true, "is"},
		{false, " "},
		{true, "your"},
		{false, " "},
		{true, "point"},
		{false, " "},
		{true, "of"},
		{false, " "},
		{true, "view"},
		{false, ", "},
		{true, "and"},
		{false, ", "},
	}

	tokenEqual(t, actual, expect)
}

func isXmlTokenizer(t *testing.T, tokenizer Tokenizer) {
	actual := tokenizer.Tokenize("<This><is>Et what is your</is><label>point of view</label></This>")
	expect := []struct {
		isWord  bool
		content string
	}{
		{false, "<This>"},
		{false, "<is>"},
		{true, "Et"},
		{false, " "},
		{true, "what"},
		{false, " "},
		{true, "is"},
		{false, " "},
		{true, "your"},
		{false, "</is>"},
		{false, "<label>"},
		{true, "point"},
		{false, " "},
		{true, "of"},
		{false, " "},
		{true, "view"},
		{false, "</label>"},
		{false, "</This>"},
	}

	tokenEqual(t, actual, expect)
}

func tokenEqual(t *testing.T, actual []Token, expect []struct {
	isWord  bool
	content string
}) {
	if len(actual) != len(expect) {
		t.Errorf("len(actual) != len(expect)")
	}
	for i, v := range actual {
		if v.IsWord() != expect[i].isWord {
			t.Errorf("v.IsWord() != expect[i].isWord")
		}

		if v.GetContent() != expect[i].content {
			t.Errorf("v.GetContent() != expect[i].content")
		}
	}
}
