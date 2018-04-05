package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

func runCmd(cmdstring string, stdin string) (int, []string) {
	if cmdstring == "" {
		cmdstring = "head -n 30"
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
	waitErr := cmd.Wait()
	if err2 != nil {
		log.Fatal(err2)
	}
	if err3 != nil {
		log.Fatal(err3)
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			errmsg := fmt.Sprintf("error executing command %v\n%s\n%s", err, string(outd), string(stde))
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), strings.Split(string(errmsg), "\n")
			}
		}
	}
	lines := strings.Split(string(outd), "\n")
	return 0, lines
}

func runMultipleCmds(cmdrune []rune, stdin string) []Buf {
	result := []Buf{}
	cmdstring := string(cmdrune)
	cmdList := strings.Split(cmdstring, "|")
	index := 0
	for _, cmdstring := range cmdList {
		status, lines := runCmd(cmdstring, stdin)
		if len(lines) >= 30 {
			lines = lines[0:30]
		}
		stdin = strings.Join(lines, "\n")
		result = append(result, Buf{Lines: lines, Status: status, Cmd: cmdstring, Index: index})
		index = index + len(cmdstring)
	}
	return result
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
	stdin := []string{}
	if len(input) > 0 {
		stdin = strings.Split(input, "\n")
	}

	state := State{
		Buffers:   []Buf{buffer},
		LineInput: li,
		Stdin:     stdin,
	}
	state.Buffers = runMultipleCmds(state.LineInput.Input, input)
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
			state.Buffers = runMultipleCmds(state.LineInput.Input, input)
			states <- state
		}
	}
	lastStateChan <- &state
}
