package util

import "fmt"

func GetCRIEndpoint(runtime string) string {
	return fmt.Sprintf("unix:///var/run/%s/%s.sock", runtime, runtime)
}
