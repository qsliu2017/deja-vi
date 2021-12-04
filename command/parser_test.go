package command

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	var tests = []struct {
		input  string
		expect Command
	}{
		{"s", &show{_content: nil, w: nil}},
		{`A "last word"`, &appendBack{_content: nil, arg: "last word"}},
		{`a "first word"`, &appendFront{_content: nil, arg: "first word"}},
		{"D 5", &deleteBack{_content: nil, n: 5, cache: ""}},
		{"d 5", &deleteFront{_content: nil, n: 5, cache: ""}},
		{"l 10", &list{invoker: nil, n: 10}},
		{"u", &undo{invoker: nil}},
		{"r", &redo{invoker: nil}},
	}

	for _, test := range tests {
		command := TheParser.Parse(test.input)
		if !reflect.DeepEqual(command, test.expect) {
			t.Errorf("%s: expected %v, got %v", test.input, test.expect, command)
		}

	}
}
