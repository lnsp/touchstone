package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetCRIEndpoint(runtime string) string {
	return fmt.Sprintf("unix:///var/run/%s/%s.sock", runtime, runtime)
}

func GetOutputTarget(file string) io.WriteCloser {
	var (
		out io.WriteCloser = os.Stdout
		err error
	)
	if file != "-" {
		out, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			out = os.Stdout
		}
	}
	return out
}

func FindPrefixedLine(data []byte, prefix string) (string, error) {
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix)), nil
		}
	}
	return "", errors.New("not found")
}
