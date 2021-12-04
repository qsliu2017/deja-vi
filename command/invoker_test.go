package command

import (
	"os"
	"testing"
)

func TestInvoker(t *testing.T) {
	var invoker Invoker = TheInvoker
	var context string = "I could be the one you dream."
	invoker.Execute(&appendFront{
		_content: &context,
		arg:      "You could be the one that I love, ",
	})

	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})

	invoker.Execute(&deleteBack{
		_content: &context,
		n:        1,
	})

	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})

	invoker.Execute(&appendBack{
		_content: &context,
		arg:      ", message in bottle is all I can do.",
	})

	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})

	if context != "You could be the one that I love, I could be the one you dream, message in bottle is all I can do." {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream, message in bottle is all I can do.' Got: %s", context)
	}

	if invoker.Execute(&undo{invoker: invoker}); context != "You could be the one that I love, I could be the one you dream" {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream' Got: %s", context)
	}

	if invoker.Execute(&undo{invoker: invoker}); context != "You could be the one that I love, I could be the one you dream." {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream.' Got: %s", context)
	}

	if invoker.Execute(&redo{invoker: invoker}); context != "You could be the one that I love, I could be the one you dream" {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream' Got: %s", context)
	}

	if invoker.Execute(&redo{invoker: invoker}); context != "You could be the one that I love, I could be the one you dream, message in bottle is all I can do." {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream, message in bottle is all I can do.' Got: %s", context)
	}
}

func TestInvokerHistory(t *testing.T) {
	var invoker Invoker = TheInvoker
	var context string = "I could be the one you dream."
	invoker.Execute(&appendFront{
		_content: &context,
		arg:      "You could be the one that I love, ",
	})

	invoker.Execute(&show{
		_content: &context,
		w:        os.Stdout,
	})

	invoker.Execute(&deleteBack{
		_content: &context,
		n:        1,
	})

	if history := invoker.history(2); len(history) != 2 {
		t.Errorf("Expected: 2 Got: %d", len(history))
	} else if history[0].(*deleteBack).n != 1 {
		t.Errorf("Expected: 1 Got: %d", history[0].(*deleteBack).n)
	}

	invoker.Execute(&undo{
		invoker: invoker,
	})

	if history := invoker.history(1); len(history) != 1 {
		t.Errorf("Expected: 1 Got: %d", len(history))
	} else if command, ok := history[0].(*appendFront); !ok {
		t.Errorf("Expected: AppendFront Got: %T", command)
	} else if command.arg != "You could be the one that I love, " {
		t.Errorf("Expected: 'You could be the one that I love, ' Got: %s", command.arg)
	}
}
