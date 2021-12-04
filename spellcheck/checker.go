package spellcheck

var (
	TheChecker Checker

	_ Checker = (*dictChecker)(nil)
	_ Checker = (*dummyChecker)(nil)
)

func init() {
	TheChecker = dummyChecker{}
}

type Checker interface {
	Check(word string) bool
}

type dict = map[string]struct{}

type dictChecker struct {
	_dict dict
}

func (dc *dictChecker) Check(word string) bool {
	_, ok := dc._dict[word]
	return ok
}

type dummyChecker struct{}

func (dummyChecker) Check(string) bool {
	return false
}
