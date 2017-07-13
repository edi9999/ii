package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func PrepareTree(tree *File, limit int64) error {
	PruneTree(tree, limit)
	if len(tree.Files) == 0 {
		return fmt.Errorf("the folder '%s' doesn't contain any files bigger than %dMB", tree.Name, limit/MEGABYTE)
	}
	SortDesc(tree)
	return nil
}

func runCmd(cmdrune []rune, stdin string) []string {
	cmdstring := string(cmdrune)
	if cmdstring == "" {
		cmdstring = "cat"
	}
	cmd := exec.Command("sh", "-c", cmdstring)
	stderr, err := cmd.StderrPipe()
	cmd.Stdin = strings.NewReader(stdin)
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	outd, err2 := ioutil.ReadAll(stdout)
	stde, err3 := ioutil.ReadAll(stderr)
	err = cmd.Wait()
	if err2 != nil {
		log.Fatal(err2)
	}
	if err3 != nil {
		log.Fatal(err3)
	}

	if err != nil {
		errmsg := fmt.Sprintf("error executing command %v\n%s\n%s", err, string(outd), string(stde))
		return strings.Split(string(errmsg), "\n")
		return strings.Split("err"+string(errmsg), "\n")
	}
	lines := strings.Split(string(outd), "\n")
	return lines
}

func StartProcessing(
	commands <-chan Executer,
	states chan<- State,
	lastStateChan chan<- *State,
	input string,
	query string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	lines := []string{"waiting for command"}
	buffer := Buf{
		Lines: lines,
	}
	li := LineInput{
		Input:  []rune(query),
		Cx:     len(query),
		Yanked: []rune{},
	}
	state := State{
		Buffers:   []Buf{buffer},
		LineInput: li,
	}
	lines = runCmd(state.LineInput.Input, input)
	buffer = Buf{Lines: lines}
	state.Buffers = []Buf{buffer}
	states <- state
	for {
		command, more := <-commands
		if !more {
			close(states)
			break
		}
		if newState, err := command.Execute(state); err == nil {
			state = newState
			lines = []string{}
			lines = runCmd(state.LineInput.Input, input)
			buffer = Buf{Lines: lines}
			state.Buffers = []Buf{buffer}
			states <- state
		}
	}
	lastStateChan <- &state
}
