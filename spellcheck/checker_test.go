package spellcheck

import "testing"

func TestDictChecker(t *testing.T) {
	var checker Checker = &dictChecker{dict{"girl": {}}}
	expect := map[string]bool{
		"girl": true,
		"boy":  false,
	}
	for v, e := range expect {
		if actual := checker.Check(v); actual != e {
			t.Logf("checker.Checker(%s), expect: %t, actual: %ts\n", v, e, actual)
		}
	}
}
