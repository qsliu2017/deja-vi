package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/qsliu2017/deja-vi/command"
	"github.com/qsliu2017/deja-vi/content"
	"github.com/qsliu2017/deja-vi/spellcheck"
)

func init() {
	initDict()
	initContent()
}

func initDict() {
	langs, _ := os.ReadDir("lang")
	for _, lang := range langs {
		f, _ := os.Open("lang/" + lang.Name())
		defer f.Close()
		words := make([]string, 0)
		rd := bufio.NewReader(f)
		for {
			word, _, err := rd.ReadLine()
			if err != nil {
				break
			}
			words = append(words, strings.TrimSpace(string(word)))
		}
		spellcheck.TheLangManager.AddLang(strings.TrimSuffix(lang.Name(), ".txt"), words)
	}
	spellcheck.TheLangManager.SetLang(strings.TrimSuffix(langs[0].Name(), ".txt"))
}

func initContent() {
	flag.StringVar(content.TheContent, "t", "", "init content")
	flag.Parse()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		print("deja-vi> ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if text == "exit" {
			break
		}

		cmd := command.TheParser.Parse(text)
		command.TheInvoker.Execute(cmd)
	}
}
