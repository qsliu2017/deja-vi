package token

var (
	_ Token     = (*token)(nil)
	_ Tokenizer = (*txtTokenizer)(nil)
	_ Tokenizer = (*xmlTokenizer)(nil)

	TheTokenizer Tokenizer
)

const (
	tok_word = iota
	tok_other
)

type Token interface {
	IsWord() bool
	GetContent() string
}

type token struct {
	kind    int
	content string
}

func (t token) IsWord() bool { return t.kind == tok_word }

func (t token) GetContent() string { return t.content }

type Tokenizer interface {
	Tokenize(string) []Token
}

type txtTokenizer struct{}

func isalpha(c byte) bool { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }

func (txtTokenizer) Tokenize(context string) (toks []Token) {
	for i, j, n, tok_kind := 0, 0, len(context), tok_other; j < n; i, toks = j, append(toks, &token{tok_kind, context[i:j]}) {
		if isalpha(context[i]) {
			for tok_kind = tok_word; j < n && isalpha(context[j]); j++ {
			}
		} else {
			for tok_kind = tok_other; j < n && !isalpha(context[j]); j++ {
			}
		}
	}
	return
}

type xmlTokenizer struct{}

func (xmlTokenizer) Tokenize(context string) (toks []Token) {
	for i, j, n, tok_kind := 0, 0, len(context), tok_other; j < n; i, toks = j, append(toks, &token{tok_kind, context[i:j]}) {
		switch {
		case context[i] == '<':
			for tok_kind = tok_other; j < n && context[j] != '>'; j++ {
			}
			if j < n {
				j++ //eat '>'
			}
		case isalpha(context[i]):
			for tok_kind = tok_word; j < n && isalpha(context[j]); j++ {
			}
		default:
			for tok_kind = tok_other; j < n && !isalpha(context[j]); j++ {
			}
		}
	}
	return
}
