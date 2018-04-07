package main

import (
	"flag"
	"github.com/edi9999/ii/core"
	"github.com/edi9999/ii/tui"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	flag.Parse()
	args := flag.Args()
	query := strings.Join(args, "")
	s := initScreen()
	commands := make(chan core.Executer)
	states := make(chan core.State)
	lastStateChan := make(chan *core.State, 1)
	var wg sync.WaitGroup
	wg.Add(3)
	bytes := []byte{}
	stat, err := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic("Failed reading from stdin")
		}
	}
	input := string(bytes)
	go core.StartProcessing(commands, states, lastStateChan, input, query, &wg)
	go tui.InteractiveTree(s, states, &wg)
	go core.ParseCommand(s, commands, &wg)
	wg.Wait()
	s.Fini()
}

func initScreen() tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		log.Printf("%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		log.Printf("%v\n", e)
		os.Exit(1)
	}
	s.Clear()
	return s
}
