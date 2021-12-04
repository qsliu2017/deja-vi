package command

import (
	"fmt"
	"io"

	"github.com/qsliu2017/deja-vi/content"
	"github.com/qsliu2017/deja-vi/spellcheck"
	"github.com/qsliu2017/deja-vi/token"
)

type Command interface {
	execute()
}

type MutCommand interface {
	Command
	unexecute()
	toString() string
	copy() MutCommand
}

var (
	_ Command    = (*show)(nil)
	_ MutCommand = (*appendFront)(nil)
	_ MutCommand = (*appendBack)(nil)
	_ MutCommand = (*deleteFront)(nil)
	_ MutCommand = (*deleteBack)(nil)
	_ Command    = (*list)(nil)
	_ Command    = (*undo)(nil)
	_ Command    = (*redo)(nil)
	_ Command    = (*defineMacro)(nil)
	_ MutCommand = (*macro)(nil)
	_ Command    = (*switchLang)(nil)
	_ Command    = (*switchContentType)(nil)
	_ Command    = (*checkSpelling)(nil)
	_ Command    = (*checkSpellingAndMark)(nil)
	_ MutCommand = (*checkSpellingAndDelete)(nil)
)

type show struct {
	_content content.Content
	w        io.Writer
}

func (c show) execute() {
	io.WriteString(c.w, *c._content+"\n")
}

type appendFront struct {
	_content content.Content
	arg      string
}

func (c appendFront) execute() {
	*c._content = c.arg + *c._content
}

func (c appendFront) unexecute() {
	*c._content = (*c._content)[len(c.arg):]
}

func (c appendFront) toString() string {
	return fmt.Sprintf("a \"%s\"", c.arg)
}

func (c appendFront) copy() MutCommand {
	return c
}

type appendBack struct {
	_content content.Content
	arg      string
}

func (c appendBack) execute() {
	*c._content = *c._content + c.arg
}

func (c appendBack) unexecute() {
	*c._content = (*c._content)[:len(*c._content)-len(c.arg)]
}

func (c appendBack) toString() string {
	return fmt.Sprintf("A \"%s\"", c.arg)
}

func (c appendBack) copy() MutCommand {
	return c
}

type deleteFront struct {
	_content content.Content
	n        int
	cache    string
}

func (c *deleteFront) execute() {
	n := c.n
	if n > len(*c._content) {
		n = len(*c._content)
	}
	c.cache = (*c._content)[:n]
	*c._content = (*c._content)[n:]
}

func (c deleteFront) unexecute() {
	*c._content = c.cache + *c._content
}

func (c deleteFront) toString() string {
	return fmt.Sprintf("d %d", c.n)
}

func (c deleteFront) copy() MutCommand {
	return &deleteFront{
		_content: c._content,
		n:        c.n,
	}
}

type deleteBack struct {
	_content content.Content
	n        int
	cache    string
}

func (c *deleteBack) execute() {
	idx := len(*c._content) - c.n
	if idx < 0 {
		idx = 0
	}
	c.cache = (*c._content)[idx:]
	*c._content = (*c._content)[:idx]
}

func (c deleteBack) unexecute() {
	*c._content = *c._content + c.cache
}

func (c deleteBack) toString() string {
	return fmt.Sprintf("D %d", c.n)
}

func (c deleteBack) copy() MutCommand {
	return &deleteBack{
		_content: c._content,
		n:        c.n,
	}
}

type list struct {
	invoker Invoker
	n       int
	w       io.Writer
}

func (c list) execute() {
	for idx, cmd := range c.invoker.history(c.n) {
		io.WriteString(c.w, fmt.Sprintf("%d %s\n", idx+1, cmd.toString()))
	}
}

type undo struct {
	invoker Invoker
}

func (c undo) execute() {
	c.invoker.undo()
}

type redo struct {
	invoker Invoker
}

func (c redo) execute() {
	c.invoker.redo()
}

type defineMacro struct {
	parser  Parser
	invoker Invoker
	name    string
	n       int
}

func (c defineMacro) execute() {
	commands := make([]MutCommand, c.n)
	for idx, cmd := range c.invoker.history(c.n) {
		commands[idx] = cmd.copy()
	}
	c.parser.defineMacro(c.name, &macro{commands: commands, name: c.name})
}

type macro struct {
	name     string
	commands []MutCommand
}

func (c macro) execute() {
	for _, cmd := range c.commands {
		cmd.execute()
	}
}

func (c macro) unexecute() {
	for i := len(c.commands) - 1; i >= 0; i-- {
		c.commands[i].unexecute()
	}
}

func (c macro) toString() string {
	return "$" + c.name
}

func (c macro) copy() MutCommand {
	commands := make([]MutCommand, len(c.commands))
	for idx, cmd := range c.commands {
		commands[idx] = cmd.copy()
	}
	return &macro{
		name:     c.name,
		commands: commands,
	}
}

type switchLang struct {
	manager spellcheck.LangManager
	lang    string
}

func (c switchLang) execute() {
	c.manager.SetLang(c.lang)
}

type switchContentType struct {
	manager     token.ContentTypeManager
	contentType string
}

func (c switchContentType) execute() {
	c.manager.SetContentType(c.contentType)
}

type checkSpelling struct {
	_content content.Content
	lexer    token.Tokenizer
	checker  spellcheck.Checker
	w        io.Writer
}

func (c checkSpelling) execute() {
	for _, token := range c.lexer.Tokenize(*c._content) {
		if content := token.GetContent(); token.IsWord() && !c.checker.Check(content) {
			io.WriteString(c.w, content+"\n")
		}
	}
	io.WriteString(c.w, "\n")
}

type checkSpellingAndMark struct {
	_content content.Content
	lexer    token.Tokenizer
	checker  spellcheck.Checker
	w        io.Writer
}

func (c checkSpellingAndMark) execute() {
	for _, token := range c.lexer.Tokenize(*c._content) {
		content := token.GetContent()
		if token.IsWord() && !c.checker.Check(content) {
			io.WriteString(c.w, "*["+content+"]")
		} else {
			io.WriteString(c.w, content)
		}
	}
	io.WriteString(c.w, "\n")
}

type checkSpellingAndDelete struct {
	_content content.Content
	lexer    token.Tokenizer
	checker  spellcheck.Checker
	w        io.Writer
	cache    string
}

func (c *checkSpellingAndDelete) execute() {
	c.cache = *c._content
	*c._content = ""
	for _, token := range c.lexer.Tokenize(c.cache) {
		content := token.GetContent()
		if !token.IsWord() || c.checker.Check(content) {
			*c._content = *c._content + content
		}
	}
}

func (c checkSpellingAndDelete) unexecute() {
	*c._content = c.cache
}

func (checkSpellingAndDelete) toString() string {
	return "spell-m"
}

func (c *checkSpellingAndDelete) copy() MutCommand {
	return &checkSpellingAndDelete{
		_content: c._content,
		lexer:    c.lexer,
		checker:  c.checker,
		w:        c.w,
	}
}
