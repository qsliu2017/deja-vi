package spellcheck

var (
	TheLangManager LangManager

	langCheckerMap map[string]Checker

	_ LangManager = (*langManager)(nil)
)

func init() {
	langCheckerMap = make(map[string]Checker)
	TheLangManager = langManager{}
}

type Dict = []string

type LangManager interface {
	SetLang(lang string)
	AddLang(lang string, dict Dict)
}

type langManager struct{}

func (langManager) SetLang(lang string) {
	if checker, ok := langCheckerMap[lang]; ok {
		TheChecker = checker
	}
}

func (langManager) AddLang(lang string, _dict Dict) {
	mapDict := make(dict)
	for _, word := range _dict {
		mapDict[word] = struct{}{}
	}
	checker := &dictChecker{mapDict}
	langCheckerMap[lang] = checker
}
