package fileformat

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Runbook struct {
	Name  string
	Steps []Step
}

type Step struct {
	Kind         StepKind
	Path         string
	Args         []string
	ConnSettings *ConnectionSettings `yaml:"conn-settings"`
}

type ConnectionSettings struct {
	Destination ConnectionAddress
	Source      ConnectionAddress

	PayloadSize int
	Insecure    bool
}

type ConnectionAddress struct {
	Host string
	Port int
}

type StepKind string

const (
	FileCreate   StepKind = "FileCreate"
	FileDelete   StepKind = "FileDelete"
	FileModify   StepKind = "FileModify"
	SendData     StepKind = "SendData"
	StartProcess StepKind = "StartProcess"
)

func Validate(path string) (Runbook, error) {
	var rb Runbook

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return rb, errors.Wrap(err, "unable to read file")
	}

	if err := yaml.Unmarshal(bytes, &rb); err != nil {
		return rb, errors.New("unable to parse runbook from file")
	}

	for i, step := range rb.Steps {
		switch step.Kind {
		case SendData:
			if step.ConnSettings == nil {
				return rb, errors.Errorf("missing connection settings for step %v", i)
			}
		}
	}

	return rb, nil
}
