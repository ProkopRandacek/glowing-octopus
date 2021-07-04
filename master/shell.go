package main

import (
	"strings"
	"fmt"

	"github.com/lmorg/readline"
)

var commands = []string{
	"writeresrc",
	"walkto",
	"drawbox",
	"drawpoint",
	"craft",
	"mine",
	"mineresource",
	"cleararea",
	"place",
	"put",
	"take",
	"build",
}

var help = map[string]string {
	"writeresrc": "[ [min_x, min_y], [max_x, max_y] ]",
	"walkto": "[ x, y ]",
	"drawbox": "{ \"color\": [ r, g, b ], \"x1\", \"y1\", \"x2\", \"y2\"}",
	"drawpoint": "{ \"color\": [ r, g, b ], \"x\", \"y\"}",
	"craft": "{ \"recipe\", \"count\"}",
	"mine": "[ x, y ]",
	"mineresource": "{ \"pos\": [ x, y ], \"amount\", \"name\" }",
	"cleararea": "{ \"area\": [ [min_x, min_y], [max_x, max_y] ], \"t\": \"all\"/\"nature\"}",
	"place": "{ \"pos\": [ x, y ], \"item\" }",
	"put": "{ \"pos\": [ x, y ], \"item\", \"amount\", \"slot\" }",
	"take": "{ \"pos\": [ x, y ], \"item\", \"amount\", \"slot\" }",
	"build": "{ TODO }",
}

func tabCompleter(text []rune, pos int, dtc readline.DelayedTabContext) (string, []string, map[string]string, readline.TabDisplayType) {
	suggestions := []string{}

	for _, c := range commands {
		if strings.HasPrefix(c, string(text)) {
			suggestions = append(suggestions, c)
		}
	}

	return string(text[pos:]), suggestions, nil, readline.TabDisplayGrid
}

func syntaxCompleter(text []rune, pos int) ([]rune, int) { // TODO
	if pos < 0 {
		return text, pos
	}

	to_add := 'n'
	switch text[pos] {
	case '"', '\'':
		to_add = text[pos]
	case '(':
		to_add = ')'
	case '[':
		to_add = ']'
	case '{':
		to_add = '}'
	}

	if to_add != 'n' {
		text = append(text[:pos], append([]rune{to_add}, text[pos+1:]...)...)
	}

	return text, pos
}

func hinter(text []rune, pos int) []rune {
	command := ""

	for _, c := range commands {
		if c == strings.Split(string(text), " ")[0] {
			command = c
			break
		}
	}

	return []rune(help[command])
}

func (b *Bot) runShell() {
	rl := readline.NewInstance()
	rl.TabCompleter = tabCompleter
	//rl.SyntaxCompleter = syntaxCompleter
	rl.HintText = hinter
	rl.SetPrompt("> ")

	for {
		text, err := rl.Readline()

		if err == readline.CtrlC {
			fmt.Println()
			break
		}

		b.conn.Execute("/" + text)
		b.waitForTaskDone()
		fmt.Println("done")
	}
}
