package level

import (
	"errors"
	"strings"
)

var levels = []string{"a", "b", "c", "d", "e", "f", "g"}

func CheckContest(name string) ([]string, error) {
	if strings.HasPrefix(name, "abc") {
		return levels, nil
	} else if strings.HasPrefix(name, "arc") {
		return levels[:6], nil
	} else if strings.HasPrefix(name, "agc") {
		return levels[:6], nil
	} else if strings.HasPrefix(name, "ahc") {
		return levels[:1], nil
	} else {
		return nil, errors.New("Contest Name is unknown")
	}
}
