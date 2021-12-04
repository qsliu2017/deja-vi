package token

var (
	TheContentTypeManager ContentTypeManager

	_ ContentTypeManager = (*contentTypeManager)(nil)

	contentTypeTokenizerMap = map[string]Tokenizer{
		"txt": txtTokenizer{},
		"xml": xmlTokenizer{},
	}
)

func init() {
	TheContentTypeManager = contentTypeManager{}
	TheContentTypeManager.SetContentType("txt")
}

type ContentTypeManager interface {
	SetContentType(contentType string)
}

type contentTypeManager struct{}

func (contentTypeManager) SetContentType(contentType string) {
	if tokennizer, ok := contentTypeTokenizerMap[contentType]; ok {
		TheTokenizer = tokennizer
	}
}
