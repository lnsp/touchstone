package util

import (
	"fmt"
	"io"
	"os"
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
