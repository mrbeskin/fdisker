package fdisker

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	resources := "./test_resources"
	tests := []struct {
		path string
	}{
		{fmt.Sprintf("%v/without_comments", resources)},
		{fmt.Sprintf("%v/with_comments", resources)},
		{fmt.Sprintf("%v/with_new_lines", resources)},
	}

	expectedCommands := []string{
		"d",
		"d",
		"n",
		"p",
		"DEF",
		"+30G",
		"w",
	}

	for _, test := range tests {
		commands, err := parseFdiskCommands(test.path)
		if err != nil {
			t.Errorf("failure in parsing commands: %v", err)
		}
		if !slicesEqual(commands, expectedCommands) {
			t.Errorf("failure in test for file %v: expected:\n%v\nbut got:\n%v\n", test.path, expectedCommands, commands)
		}
	}
}

func slicesEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, val := range s1 {
		if val != s2[i] {
			return false
		}
	}
	return true
}
