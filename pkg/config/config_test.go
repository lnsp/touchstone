package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	tt := []struct {
		Name    string
		Content []byte
		Config  *Config
	}{
		{
			Name: "performance",
			Content: []byte(`
output: performance.yaml
oci: ["runc", "runsc"]
cri: ["containerd", "crio"]
filter:
- performance
runs: 10
scale: 1
`),
			Config: &Config{
				Output: "performance.yaml",
				OCIs:   []string{"runc", "runsc"},
				CRIs:   []string{"containerd", "crio"},
				Filter: []string{"performance"},
				Runs:   10,
				Scale:  1,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			// write tmp file
			tmpFile, err := ioutil.TempFile("", "config_test")
			if err != nil {
				t.Fatalf("could not create tmpfile: %v", err)
			}
			_, err = tmpFile.Write(tc.Content)
			if err != nil {
				t.Fatalf("could not write to tmpfile: %v", err)
			}
			tmpFile.Close()
			cfg, err := Parse(tmpFile.Name())
			if err != nil {
				t.Fatalf("could not parse config: %v", err)
			}
			if !reflect.DeepEqual(cfg, tc.Config) {
				t.Errorf("expected %v, got %v", tc.Config, cfg)
			}
			if err := os.Remove(tmpFile.Name()); err != nil {
				t.Fatalf("could not remove tmpfile: %v", err)
			}
		})
	}

}
