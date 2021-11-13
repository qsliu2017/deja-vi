package dejavi

import (
	"fmt"
	"io"
)

type Command interface {
	Execute()
}

type MutCommand interface {
	Command
	Unexecuted()
	toString() string
}

var (
	_ Command    = (*Show)(nil)
	_ MutCommand = (*AppendFront)(nil)
	_ MutCommand = (*AppendBack)(nil)
	_ MutCommand = (*DeleteFront)(nil)
	_ MutCommand = (*DeleteBack)(nil)
	_ Command    = (*List)(nil)
	_ Command    = (*Undo)(nil)
	_ Command    = (*Redo)(nil)
)

type Show struct {
	context *string
	out     io.Writer
}

func (command *Show) Execute() {
	io.WriteString(command.out, *command.context+"\n")
}

type AppendFront struct {
	context *string
	arg     string
	stash   string
}

func (command *AppendFront) Execute() {
	command.stash = *command.context
	*command.context = command.arg + *command.context
}

func (command *AppendFront) Unexecuted() {
	*command.context = command.stash
}

func (command *AppendFront) toString() string {
	return fmt.Sprintf("a \"%s\"", command.arg)
}

type AppendBack struct {
	context *string
	arg     string
	stash   string
}

func (command *AppendBack) Execute() {
	command.stash = *command.context
	*command.context = *command.context + command.arg
}

func (command *AppendBack) Unexecuted() {
	*command.context = command.stash
}

func (command *AppendBack) toString() string {
	return fmt.Sprintf("A \"%s\"", command.arg)
}

type DeleteFront struct {
	context *string
	arg     int
	stash   string
}

func (command *DeleteFront) Execute() {
	command.stash = *command.context
	*command.context = (*command.context)[command.arg:]
}

func (command *DeleteFront) Unexecuted() {
	*command.context = command.stash
}

func (command *DeleteFront) toString() string {
	return fmt.Sprintf("d %d", command.arg)
}

type DeleteBack struct {
	context *string
	arg     int
	stash   string
}

func (command *DeleteBack) Execute() {
	command.stash = *command.context
	*command.context = (*command.context)[:len(*command.context)-command.arg]
}

func (command *DeleteBack) Unexecuted() {
	*command.context = command.stash
}

func (command *DeleteBack) toString() string {
	return fmt.Sprintf("D %d", command.arg)
}

type List struct {
	context Invoker
	arg     int
	out     io.Writer
}

func (command *List) Execute() {
	for idx, cmd := range command.context.history(command.arg) {
		io.WriteString(command.out, fmt.Sprintf("\t%d\t%s\n", idx+1, cmd.toString()))
	}
}

type Undo struct {
	context Invoker
}

func (command *Undo) Execute() {
	command.context.undo()
}

type Redo struct {
	context Invoker
}

func (command *Redo) Execute() {
	command.context.redo()
}
