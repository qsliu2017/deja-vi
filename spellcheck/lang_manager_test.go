package spellcheck

import "testing"

func TestInit(t *testing.T) {
	if langCheckerMap == nil {
		t.Error("langDictMap is nil")
	}

	if TheLangManager == nil {
		t.Error("TheLangManager is nil")
	}

	if TheChecker == nil {
		t.Error("TheChecker is nil")
	}
}

func TestAddLang(t *testing.T) {
	lang := "en"
	dict := []string{
		"What",
		"what",
		"Is",
		"is",
		"Your",
		"your",
		"Point",
		"point",
		"Of",
		"of",
		"View",
		"view",
		"And",
		"and",
	}

	TheLangManager.AddLang(lang, dict)
	TheLangManager.SetLang(lang)

	matrix := []struct {
		word   string
		expect bool
	}{
		{"What", true},
		{"what", true},
		{"Is", true},
		{"is", true},
		{"Your", true},
		{"your", true},
		{"Quel", false},
		{"quel", false},
		{"Est", false},
		{"est", false},
		{"Votre", false},
		{"votre", false},
	}

	for _, m := range matrix {
		if TheChecker.Check(m.word) != m.expect {
			t.Error("TheChecker.Check(", m.word, ") !=", m.expect)
		}
	}
}
