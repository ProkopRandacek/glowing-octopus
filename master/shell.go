package main

import (
	"fmt"
	"github.com/lmorg/readline"
	"strings"
)

var commands = []string{
	"build",
	"cleararea",
	"craft",
	"drawbox",
	"drawpoint",
	"mine",
	"mineresource",
	"place",
	"put",
	"take",
	"walkto",
	"writeresrc",
}

var help = map[string]string{
	"build":        "{ TODO }",
	"cleararea":    "{ \"area\": [ [min_x, min_y], [max_x, max_y] ], \"t\": \"all\"/\"nature\"}",
	"craft":        "{ \"recipe\", \"count\"}",
	"drawbox":      "{ \"color\": [ r, g, b ], \"x1\", \"y1\", \"x2\", \"y2\"}",
	"drawpoint":    "{ \"color\": [ r, g, b ], \"x\", \"y\"}",
	"mine":         "[ x, y ]",
	"mineresource": "{ \"pos\": [ x, y ], \"amount\", \"name\" }",
	"place":        "{ \"pos\": [ x, y ], \"item\" }",
	"put":          "{ \"pos\": [ x, y ], \"item\", \"amount\", \"slot\" }",
	"take":         "{ \"pos\": [ x, y ], \"item\", \"amount\", \"slot\" }",
	"walkto":       "[ x, y ]",
	"writeresrc":   "[ [min_x, min_y], [max_x, max_y] ]",
}

func tabCompleter(text []rune, pos int, dtc readline.DelayedTabContext) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string

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

	toAdd := 'n'
	switch text[pos] {
	case '"', '\'':
		toAdd = text[pos]
	case '(':
		toAdd = ')'
	case '[':
		toAdd = ']'
	case '{':
		toAdd = '}'
	}

	if toAdd != 'n' {
		text = append(text[:pos], append([]rune{toAdd}, text[pos+1:]...)...)
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

func (b *bot) runShell() error {
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

		_, err = b.conn.Execute("/" + text)
		if err != nil {
			return err
		}
		b.waitForTaskDone()
		fmt.Println("done")
	}
	return nil
}
