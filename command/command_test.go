package command

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestShow(t *testing.T) {
	var content string = "I remember it all too well"

	r, w := io.Pipe()
	go func() { //PipeWriter is a blocking operation
		defer w.Close()
		show := &show{
			_content: &content,
			w:        w,
		}
		show.execute()
	}()
	defer r.Close()

	output, _ := io.ReadAll(r)

	if actual := string(output); strings.Compare(actual, "I remember it all too well\n") != 0 {
		t.Errorf("Expected 'I remember it all too well\\n' but got '%s'", actual)
	}
}

func TestAppendFront(t *testing.T) {
	var context string = "I remember it all too well"

	appendFront := &appendFront{
		_content: &context,
		arg:      "I was there, ",
	}

	if appendFront.execute(); strings.Compare(context, "I was there, I remember it all too well") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well' but got '%s'", context)
	}

	if appendFront.unexecute(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestAppendBack(t *testing.T) {
	var context string = "I remember it all too well"

	appendBack := &appendBack{
		_content: &context,
		arg:      " I was there",
	}

	if appendBack.execute(); strings.Compare(context, "I remember it all too well I was there") != 0 {
		t.Errorf("Expected 'I remember it all too well I was there' but got '%s'", context)
	}

	if appendBack.unexecute(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestDeleteFront(t *testing.T) {
	var context string = "I remember it all too well"

	deleteFront := &deleteFront{
		_content: &context,
		n:        14,
	}

	if deleteFront.execute(); strings.Compare(context, "all too well") != 0 {
		t.Errorf("Expected 'all too well' but got '%s'", context)
	}

	if deleteFront.unexecute(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestDeleteBack(t *testing.T) {
	var context string = "I remember it all too well"

	deleteBack := &deleteBack{
		_content: &context,
		n:        9,
	}

	if deleteBack.execute(); strings.Compare(context, "I remember it all") != 0 {
		t.Errorf("Expected 'I remember it' but got '%s'", context)
	}

	if deleteBack.unexecute(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestList(t *testing.T) {
	var context string = "I remember it all too well"
	var invoker Invoker = TheInvoker

	// execute some commands and finally list
	invoker.Execute(&appendBack{
		_content: &context,
		arg:      " I was there",
	})
	invoker.Execute(&appendFront{
		_content: &context,
		arg:      "I was there, ",
	})
	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})
	invoker.Execute(&deleteBack{
		_content: &context,
		n:        9,
	})
	invoker.Execute(&deleteFront{
		_content: &context,
		n:        14,
	})

	r, w := io.Pipe()
	go func() { //PipeWriter is a blocking operation
		defer w.Close()
		invoker.Execute(&list{
			invoker: invoker,
			n:       3,
			w:       w,
		})
	}()
	defer r.Close()

	output, _ := io.ReadAll(r)
	expect := `1 d 14
2 D 9
3 a "I was there, "
`

	if actual := string(output); strings.Compare(actual, expect) != 0 {
		t.Errorf("Expected '%s' but got '%s'", expect, actual)
	}
}

func TestUndo(t *testing.T) {
	var context string = "I remember it all too well"
	var invoker Invoker = TheInvoker

	// execute some commands and finally list
	invoker.Execute(&appendBack{
		_content: &context,
		arg:      " I was there",
	})
	invoker.Execute(&appendFront{
		_content: &context,
		arg:      "I was there, ",
	})
	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})
	invoker.Execute(&deleteBack{
		_content: &context,
		n:        12,
	})

	if strings.Compare(context, "I was there, I remember it all too well") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well' but got '%s'", context)
	}

	if invoker.Execute(&undo{invoker}); strings.Compare(context, "I was there, I remember it all too well I was there") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well I was there' but got '%s'", context)
	}

	if invoker.Execute(&undo{invoker}); strings.Compare(context, "I remember it all too well I was there") != 0 {
		t.Errorf("Expected 'I remember it all too well I was there' but got '%s'", context)
	}

	if invoker.Execute(&undo{invoker}); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestRedo(t *testing.T) {
	var context string = "I remember it all too well"
	var invoker Invoker = TheInvoker

	// execute some commands and finally list
	invoker.Execute(&appendBack{
		_content: &context,
		arg:      " I was there",
	})
	invoker.Execute(&appendFront{
		_content: &context,
		arg:      "I was there, ",
	})
	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})
	invoker.Execute(&deleteBack{
		_content: &context,
		n:        12,
	})

	if strings.Compare(context, "I was there, I remember it all too well") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well' but got '%s'", context)
	}

	if invoker.Execute(&undo{invoker}); strings.Compare(context, "I was there, I remember it all too well I was there") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well I was there' but got '%s'", context)
	}

	if invoker.Execute(&redo{invoker}); strings.Compare(context, "I was there, I remember it all too well") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well' but got '%s'", context)
	}
}

func TestDefineMacro(t *testing.T) {
	var context string = "I remember it all too well"
	var invoker Invoker = TheInvoker
	var parser Parser = TheParser

	// execute some commands, then define macro and execute that macro
	invoker.Execute(&appendBack{
		_content: &context,
		arg:      " I was there",
	})
	invoker.Execute(&appendFront{
		_content: &context,
		arg:      "I was there, ",
	})
	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})
	invoker.Execute(&defineMacro{parser, invoker, "myMacro", 2})

	macro := parser.Parse("$myMacro")

	if invoker.Execute(macro); strings.Compare(context, "I was there, I was there, I remember it all too well I was there I was there") != 0 {
		t.Errorf("Expected 'I was there, I was there, I remember it all too well I was there I was there' but got '%s'", context)
	}

	if invoker.Execute(&undo{invoker}); strings.Compare(context, "I was there, I remember it all too well I was there") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well I was there' but got '%s'", context)
	}

	if invoker.Execute(&undo{invoker}); strings.Compare(context, "I remember it all too well I was there") != 0 {
		t.Errorf("Expected 'I remember it all too well I was there' but got '%s'", context)
	}

	if invoker.Execute(&undo{invoker}); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}
