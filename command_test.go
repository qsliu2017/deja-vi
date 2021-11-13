package dejavi

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestShow(t *testing.T) {
	var context string = "I remember it all too well"

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	show := &Show{
		context: &context,
		out:     w,
	}
	go show.Execute() //PipeWriter is a blocking operation
	buffer := make([]byte, 1024)
	n, err := r.Read(buffer)
	if err != nil {
		t.Fatalf("Error reading from pipe: %s", err)
	}

	if actual := string(buffer[:n]); strings.Compare(actual, "I remember it all too well\n") != 0 {
		t.Errorf("Expected 'I remember it all too well\\n' but got '%s'", actual)
	}

}

func TestShowToStdout(t *testing.T) {
	var context string = "I remember it all too well"

	show := &Show{
		context: &context,
		out:     os.Stdout,
	}

	show.Execute()
}

func TestAppendFront(t *testing.T) {
	var context string = "I remember it all too well"

	appendFront := &AppendFront{
		context: &context,
		arg:     "I was there, ",
	}

	if appendFront.Execute(); strings.Compare(context, "I was there, I remember it all too well") != 0 {
		t.Errorf("Expected 'I was there, I remember it all too well' but got '%s'", context)
	}

	if appendFront.Unexecuted(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestAppendBack(t *testing.T) {
	var context string = "I remember it all too well"

	appendBack := &AppendBack{
		context: &context,
		arg:     " I was there",
	}

	if appendBack.Execute(); strings.Compare(context, "I remember it all too well I was there") != 0 {
		t.Errorf("Expected 'I remember it all too well I was there' but got '%s'", context)
	}

	if appendBack.Unexecuted(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestDeleteFront(t *testing.T) {
	var context string = "I remember it all too well"

	deleteFront := &DeleteFront{
		context: &context,
		arg:     14,
	}

	if deleteFront.Execute(); strings.Compare(context, "all too well") != 0 {
		t.Errorf("Expected 'all too well' but got '%s'", context)
	}

	if deleteFront.Unexecuted(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestDeleteBack(t *testing.T) {
	var context string = "I remember it all too well"

	deleteBack := &DeleteBack{
		context: &context,
		arg:     9,
	}

	if deleteBack.Execute(); strings.Compare(context, "I remember it all") != 0 {
		t.Errorf("Expected 'I remember it' but got '%s'", context)
	}

	if deleteBack.Unexecuted(); strings.Compare(context, "I remember it all too well") != 0 {
		t.Errorf("Expected 'I remember it all too well' but got '%s'", context)
	}
}

func TestList(t *testing.T) {
	var context string = "I remember it all too well"
	var invoker Invoker = NewInvoker()

	// execute some commands and finally list
	invoker.Execute(&AppendBack{
		context: &context,
		arg:     " I was there",
	})

	invoker.Execute(&AppendFront{
		context: &context,
		arg:     "I was there, ",
	})

	invoker.Execute(&Show{
		context: &context,
		out:     os.Stdout,
	})

	invoker.Execute(&DeleteBack{
		context: &context,
		arg:     9,
	})

	invoker.Execute(&DeleteFront{
		context: &context,
		arg:     14,
	})

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	go invoker.Execute(&List{
		context: invoker,
		arg:     3,
		out:     w,
	})
	buffer := make([]byte, 1024)

	if n, err := r.Read(buffer); err != nil {
		t.Fatalf("Error reading from pipe: %s", err)
	} else if actual := string(buffer[:n]); strings.Compare(actual,
		"\t1\td 14\n",
	) != 0 {
		t.Errorf("Expected '\t1\td 14\n' but got '%s'", actual)
	}

	if n, err := r.Read(buffer); err != nil {
		t.Fatalf("Error reading from pipe: %s", err)
	} else if actual := string(buffer[:n]); strings.Compare(actual,
		"\t2\tD 9\n",
	) != 0 {
		t.Errorf("Expected '\t2\tD 9\n' but got '%s'", actual)
	}

	if n, err := r.Read(buffer); err != nil {
		t.Fatalf("Error reading from pipe: %s", err)
	} else if actual := string(buffer[:n]); strings.Compare(actual,
		"\t3\ta \"I was there, \"\n",
	) != 0 {
		t.Errorf("Expected '\t3\ta \"I was there, \"\n' but got '%s'", actual)
	}
}
