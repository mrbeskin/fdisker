package fdisker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// DEFAULT is the key token for sending default to fdisk in the form of a newline
const DEFAULT = "DEF"
const fdiskCmd = "fdisk"
const writeCommand = "w"
const quitCommand = "q"

// RunFdiskCommandFile reads from a file, converts it to a set of fdisk commands,
// then starts fdisk and runs those commands on the passed in mounted volume.
//
// fdisk output will be piped to stdout and errors to stderr
// to save changes, writeFlag must be set to true, otherwise fdisk will quit without
// saving changes
func RunFdiskCommandFile(path string, mountPath string, writeFlag bool) error {
	// parse commands
	commands, err := parseFdiskCommands(path)
	if err != nil {
		return fmt.Errorf("could not parse commands from provided file: %v", err)
	}
	// start fdisk
	fdisk := exec.Command(fdiskCmd, mountPath)
	fdisk.Stdout = os.Stdout
	fdisk.Stderr = os.Stderr
	fdiskReader, fdiskInWriter := io.Pipe()
	fdisk.Stdin = fdiskInReader
	err = fdisk.Start()
	if err != nil {
		return fmt.Errorf("starting command: %v", err)
	}
	defer func() {
		if writeFlag {
			err := executeCommand(writeCommand, fdiskInWriter)
			if err != nil {
				panic(fmt.Errorf("failed to write changes during exit: %v", err))
			}
		} else {
			err := executeCommand(quitCommand, fdiskInWriter)
			if err != nil {
				panic(fmt.Errorf("failed to successfully exit on quit without write: %v", err))
			}
		}
		err := fdisk.Wait()
		if err != nil {
			panic(fmt.Errorf("failure for command to end: %v", err))
		}
	}()
	// pipe commands to fdisk
	// ANY ERROR RETURNS PAST THIS POINT SHOULD SET writeFlag to FALSE
	for i, command := range commands {
		err := executeCommand(command, fdiskInWriter)
		if err != nil {
			writeFlag = false
			return fmt.Errorf("failed on command \"%v\" which was command number %v: %v", command, i+1, err)
		}
	}
	return nil
}

func executeCommand(command string, w io.Writer) error {
	newLine := []byte("\n")
	if command != "DEF" {
		_, err := w.Write([]byte(command))
		if err != nil {
			return fmt.Errorf("writing command %v: %v", command, err)
		}
	}
	_, err := w.Write(newLine)
	if err != nil {
		return fmt.Errorf("return after writing command: %v", err)
	}
	return nil
}
