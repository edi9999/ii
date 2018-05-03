package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

var cmdBlackList = []string{"rm", "mv", "su", "sudo", "vim", "vi", "top", "htop", "nano", "emacs", "xargs", "trash", "ii", "fzf", "gedit"}

func runCmd(cmdstring string, stdin string) (int, []string) {
	fixedCmd := strings.Trim(cmdstring, " ")
	for _, disallowedCmd := range cmdBlackList {
		if fixedCmd == disallowedCmd || strings.HasPrefix(fixedCmd, disallowedCmd+" ") {
			return 50, []string{"Command '" + disallowedCmd + "' not allowed"}
		}
	}
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

var firstNonEmptyChar = regexp.MustCompile("[^ ]")

func runMultipleCmds(cmdList []string, oldCmdList []string, stdin string, oldBuffers []Buf) []Buf {
	buffers := []Buf{}
	if len(stdin) > 0 {
		buffers = append(buffers, Buf{
			Lines: strings.Split(stdin, "\n"),
			Stdin: true,
		})
	}
	index := 0
	same := true
	for i, cmdstring := range cmdList {
		cmdbyte := []byte(cmdstring)
		cmdLen := len(cmdstring)
		offset := firstNonEmptyChar.FindIndex(cmdbyte)
		cmdstring = normalizeCmd(cmdstring)
		if same == true {
			if len(oldCmdList) <= i {
				same = false
			} else {
				oldCmdString := normalizeCmd(oldCmdList[i])
				if oldCmdString != cmdstring {
					same = false
				}
			}
			if len(oldBuffers) <= len(buffers) {
				same = false
			}
		}
		ioutil.WriteFile("/tmp/ii.log", []byte(fmt.Sprintf("%i, %s\n", i, same)), 0644)
		if same == true {
			ioutil.WriteFile("/tmp/ii.log", []byte(fmt.Sprintf("Running %s %s \n", cmdstring, same)), 0644)
			b := oldBuffers[len(buffers)]
			b.Index = index + offset[0]
			buffers = append(buffers, b)
			stdin = strings.Join(b.Lines, "\n")
		} else {
			if len(cmdstring) > 0 {
				ioutil.WriteFile("/tmp/ii.log", []byte(fmt.Sprintf("Running %s %s \n", cmdstring, same)), 0644)
				status, lines := runCmd(cmdstring, stdin)
				stdin = strings.Join(lines, "\n")
				buffers = append(buffers, Buf{Lines: lines, Status: status, Cmd: cmdstring, Index: index + offset[0], Scroll: 0})
			}
		}
		index = index + cmdLen
	}
	return buffers
}

func getSelectedWidget(cmdList []string, cursorPosition int) int {
	totalLength := 0
	for i, str := range cmdList {
		totalLength = totalLength + len(str) + 1
		if cursorPosition < totalLength {
			return i
		}
	}
	return len(cmdList) - 1
}

func getCmdList(cmdString string) []string {
	return strings.Split(cmdString, "|")
}
func normalizeCmd(str string) string {
	return strings.Trim(str, " ")
}

func getNewState(stdin string, oldState State, newState State) State {
	update := string(newState.LineInput.Input) != string(oldState.LineInput.Input)
	cmdList := getCmdList(string(newState.LineInput.Input))
	oldCmdList := getCmdList(string(oldState.LineInput.Input))
	if update {
		newState.Buffers = runMultipleCmds(cmdList, oldCmdList, stdin, oldState.Buffers)
	}
	newState.SelectedWidget = getSelectedWidget(getCmdList(string(newState.LineInput.Input)), newState.LineInput.Cx)
	if len(newState.Buffers) > 0 && newState.Buffers[0].Stdin {
		newState.SelectedWidget = 1 + newState.SelectedWidget
	}
	return newState
}

func ProcessCommands(
	commands <-chan Executer,
	states chan<- State,
	lastStateChan chan<- *State,
	stdin string,
	query string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	oldState := State{
		LineInput: LineInput{
			Input:  []rune{},
			Cx:     0,
			Yanked: []rune{},
		},
	}
	newState := State{
		LineInput: LineInput{
			Input:  []rune(query),
			Cx:     len(query),
			Yanked: []rune{},
		},
	}
	state := getNewState(stdin, oldState, newState)

	states <- state
	for {
		command, more := <-commands
		if !more {
			close(states)
			break
		}
		if newState, err := command.Execute(state); err == nil {
			oldState := state
			state = getNewState(
				stdin,
				oldState,
				newState)
			states <- state
		}
	}
	lastStateChan <- &state
}
