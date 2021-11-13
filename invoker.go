package dejavi

type Invoker interface {
	Execute(Command)
	history(int) []MutCommand
	undo()
	redo()
}

func NewInvoker() Invoker {
	return &invokerImpl{
		stack: make([]MutCommand, 0),
		stash: make([]MutCommand, 0),
	}
}

var _ Invoker = (*invokerImpl)(nil)

type invokerImpl struct {
	stack []MutCommand
	stash []MutCommand
}

func (invoker *invokerImpl) Execute(command Command) {
	if mutCommand, ok := command.(MutCommand); ok {
		// push back to stack
		invoker.stack = append(invoker.stack, mutCommand)
		// clear the stash
		invoker.stash = invoker.stash[:0]
	}
	command.Execute()
}

func (invoker *invokerImpl) history(n int) []MutCommand {
	if n > len(invoker.stack) {
		n = len(invoker.stack)
	}
	reverseStack := make([]MutCommand, n)
	for i := 0; i < n; i++ {
		reverseStack[i] = invoker.stack[len(invoker.stack)-i-1]
	}
	return reverseStack
}

func (invoker *invokerImpl) undo() {
	if len(invoker.stack) > 0 {
		command := invoker.stack[len(invoker.stack)-1]
		invoker.stack = invoker.stack[:len(invoker.stack)-1]
		invoker.stash = append(invoker.stash, command)
		command.Unexecuted()
	}
}

func (invoker *invokerImpl) redo() {
	if len(invoker.stash) > 0 {
		command := invoker.stash[len(invoker.stash)-1]
		invoker.stash = invoker.stash[:len(invoker.stash)-1]
		invoker.stack = append(invoker.stack, command)
		command.Execute()
	}
}
