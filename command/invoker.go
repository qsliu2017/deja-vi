package command

import "github.com/qsliu2017/deja-vi/content"

var (
	TheInvoker Invoker

	_ Invoker = (*invoker)(nil)
)

func init() {
	TheInvoker = &invoker{make([]MutCommand, 0), make([]MutCommand, 0)}
}

type Invoker interface {
	Execute(Command)
	history(int) []MutCommand
	undo()
	redo()
}

type invoker struct {
	stack []MutCommand
	stash []MutCommand
}

func (invoker *invoker) Execute(command Command) {
	command.execute()
	if mutCommand, ok := command.(MutCommand); ok {
		// push back to stack
		invoker.stack = append(invoker.stack, mutCommand)
		// clear the stash
		invoker.stash = invoker.stash[:0]
		// print out context
		println(*content.TheContent)
	}
}

func (invoker *invoker) history(n int) []MutCommand {
	if n > len(invoker.stack) {
		n = len(invoker.stack)
	}
	reverseStack := make([]MutCommand, n)
	for i := 0; i < n; i++ {
		reverseStack[i] = invoker.stack[len(invoker.stack)-i-1]
	}
	return reverseStack
}

func (invoker *invoker) undo() {
	if len(invoker.stack) > 0 {
		command := invoker.stack[len(invoker.stack)-1]
		invoker.stack = invoker.stack[:len(invoker.stack)-1]
		invoker.stash = append(invoker.stash, command)
		command.unexecute()
	}
}

func (invoker *invoker) redo() {
	if len(invoker.stash) > 0 {
		command := invoker.stash[len(invoker.stash)-1]
		invoker.stash = invoker.stash[:len(invoker.stash)-1]
		invoker.stack = append(invoker.stack, command)
		command.execute()
	}
}
