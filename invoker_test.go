package dejavi

import (
	"os"
	"testing"
)

func TestInvoker(t *testing.T) {
	var invoker Invoker = NewInvoker()
	var context string = "I could be the one you dream."
	invoker.Execute(&AppendFront{
		context: &context,
		arg:     "You could be the one that I love, ",
	})

	invoker.Execute(&Show{
		context: &context,
		out:     os.Stdout,
	})

	invoker.Execute(&DeleteBack{
		context: &context,
		arg:     1,
	})

	invoker.Execute(&Show{
		context: &context,
		out:     os.Stdout,
	})

	invoker.Execute(&AppendBack{
		context: &context,
		arg:     ", message in bottle is all I can do.",
	})

	invoker.Execute(&Show{
		context: &context,
		out:     os.Stdout,
	})

	if context != "You could be the one that I love, I could be the one you dream, message in bottle is all I can do." {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream, message in bottle is all I can do.' Got: %s", context)
	}

	if invoker.Execute(&Undo{context: invoker}); context != "You could be the one that I love, I could be the one you dream" {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream' Got: %s", context)
	}

	if invoker.Execute(&Undo{context: invoker}); context != "You could be the one that I love, I could be the one you dream." {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream.' Got: %s", context)
	}

	if invoker.Execute(&Redo{context: invoker}); context != "You could be the one that I love, I could be the one you dream" {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream' Got: %s", context)
	}

	if invoker.Execute(&Redo{context: invoker}); context != "You could be the one that I love, I could be the one you dream, message in bottle is all I can do." {
		t.Errorf("Expected: 'You could be the one that I love, I could be the one you dream, message in bottle is all I can do.' Got: %s", context)
	}
}

func TestInvokerHistory(t *testing.T) {
	var invoker Invoker = NewInvoker()
	var context string = "I could be the one you dream."
	invoker.Execute(&AppendFront{
		context: &context,
		arg:     "You could be the one that I love, ",
	})

	invoker.Execute(&Show{
		context: &context,
		out:     os.Stdout,
	})

	invoker.Execute(&DeleteBack{
		context: &context,
		arg:     1,
	})

	if history := invoker.history(2); len(history) != 2 {
		t.Errorf("Expected: 2 Got: %d", len(history))
	} else if history[0].(*DeleteBack).arg != 1 {
		t.Errorf("Expected: 1 Got: %d", history[0].(*DeleteBack).arg)
	}

	invoker.Execute(&Undo{
		context: invoker,
	})

	if history := invoker.history(1); len(history) != 1 {
		t.Errorf("Expected: 1 Got: %d", len(history))
	} else if command, ok := history[0].(*AppendFront); !ok {
		t.Errorf("Expected: AppendFront Got: %T", command)
	} else if command.arg != "You could be the one that I love, " {
		t.Errorf("Expected: 'You could be the one that I love, ' Got: %s", command.arg)
	}
}
