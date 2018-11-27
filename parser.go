package fdisker

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func parseFdiskCommands(path string) ([]string, error) {
	fbytes, err := ioutil.ReadFile(path)
	f := string(fbytes)
	if err != nil {
		return []string{}, fmt.Errorf("reading file: %v", err)
	}
	lines := strings.Split(f, "\n")
	commands := []string{}
	for _, line := range lines {
		line = strings.Trim(line, " ")
		// skip comment lines
		if strings.HasPrefix(line, "#") {
			continue
		}
		toks := strings.Split(line, "#")
		command := strings.Trim(toks[0], " ")
		if command != "" {
			commands = append(commands, command)
		}
	}
	return commands, nil
}
