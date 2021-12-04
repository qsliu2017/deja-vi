package command

import (
	"os"
	"strconv"
	"strings"

	"github.com/qsliu2017/deja-vi/content"
	"github.com/qsliu2017/deja-vi/spellcheck"
	"github.com/qsliu2017/deja-vi/token"
)

var (
	TheParser Parser

	_ Parser = (*parser)(nil)
)

func init() {
	TheParser = &parser{make(map[string]*macro)}
}

type Parser interface {
	Parse(string) Command

	defineMacro(string, *macro)
}

type parser struct {
	macros map[string]*macro
}

func (p *parser) Parse(input string) Command {
	args := strings.Split(input, " ")
	switch args[0] {
	case "s":
		c := new(show)
		c._content = content.TheContent
		c.w = os.Stdout
		return c
	case "A":
		start, end := strings.Index(input, "\""), strings.LastIndex(input, "\"")
		c := new(appendBack)
		c._content = content.TheContent
		c.arg = input[start+1 : end]
		return c
	case "a":
		start, end := strings.Index(input, "\""), strings.LastIndex(input, "\"")
		c := new(appendFront)
		c._content = content.TheContent
		c.arg = input[start+1 : end]
		return c
	case "D":
		n, _ := strconv.Atoi(args[1])
		c := new(deleteBack)
		c._content = content.TheContent
		c.n = n
		return c
	case "d":
		n, _ := strconv.Atoi(args[1])
		c := new(deleteFront)
		c._content = content.TheContent
		c.n = n
		return c
	case "l":
		n, _ := strconv.Atoi(args[1])
		c := new(list)
		c.invoker = TheInvoker
		c.n = n
		c.w = os.Stdout
		return c
	case "u":
		c := new(undo)
		c.invoker = TheInvoker
		return c
	case "r":
		c := new(redo)
		c.invoker = TheInvoker
		return c
	case "m":
		n, _ := strconv.Atoi(args[1])
		name := args[2]
		c := new(defineMacro)
		c.parser = TheParser
		c.name = name
		c.n = n
		c.invoker = TheInvoker
		return c
	case "lang":
		lang := args[1]
		c := new(switchLang)
		c.manager = spellcheck.TheLangManager
		c.lang = lang
		return c
	case "content":
		contentType := args[1]
		c := new(switchContentType)
		c.manager = token.TheContentTypeManager
		c.contentType = contentType
		return c
	case "spell":
		c := new(checkSpelling)
		c._content = content.TheContent
		c.lexer = token.TheTokenizer
		c.checker = spellcheck.TheChecker
		c.w = os.Stdout
		return c
	case "spell-a":
		c := new(checkSpellingAndMark)
		c._content = content.TheContent
		c.lexer = token.TheTokenizer
		c.checker = spellcheck.TheChecker
		c.w = os.Stdout
		return c
	case "spell-m":
		c := new(checkSpellingAndDelete)
		c._content = content.TheContent
		c.lexer = token.TheTokenizer
		c.checker = spellcheck.TheChecker
		c.w = os.Stdout
		return c
	default: //defined macro
		if macro, ok := p.macros[args[0]]; ok {
			return macro.copy()
		}
	}
	return nil
}

func (p *parser) defineMacro(name string, command *macro) {
	p.macros["$"+name] = command
}
