package fdisker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
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
	fdiskReader, fdiskWriter := io.Pipe()
	fdisk.Stdin = fdiskReader
	err = fdisk.Start()
	if err != nil {
		return fmt.Errorf("starting command: %v", err)
	}
	// pipe commands to fdisk
	errs := make([]error, 0)
	for i, command := range commands {
		time.Sleep(500 * time.Millisecond)
		err := executeCommand(command, fdiskWriter)
		if err != nil {
			writeFlag = false
			errs = append(errs, fmt.Errorf("failed on command \"%v\" which was command number %v: %v", command, i+1, err))
		}
	}
	err = quitFdisk(fdiskWriter, fdisk, writeFlag)
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return compositeError(errs)
	}
	errChan := make(chan error, 1)
	go func(c chan error) {
		c <- fdisk.Wait()
	}(errChan)
	fdiskWriter.Close()
	err = <-errChan
	if err != nil {
		return err
	}
	return nil
}

func quitFdisk(fdiskWriter io.Writer, fdisk *exec.Cmd, writeFlag bool) error {
	time.Sleep(1 * time.Second)
	if writeFlag {
		err := executeCommand(writeCommand, fdiskWriter)
		if err != nil {
			return fmt.Errorf("failed to write changes during exit: %v", err)
		}
	} else {
		err := executeCommand(quitCommand, fdiskWriter)
		if err != nil {
			return fmt.Errorf("failed to successfully exit on quit without write: %v", err)
		}
	}
	return nil
}

func executeCommand(command string, w io.Writer) error {
	newLine := []byte("\n")
	if command != "DEF" {
		fmt.Println(command)
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

func compositeError(errs []error) error {
	err := ""
	for i := len(errs) - 1; i >= 0; i-- {
		if err != "" {
			err = err + ": " + errs[i].Error()
		} else {
			err = errs[i].Error()
		}
	}
	return fmt.Errorf(err)
}
