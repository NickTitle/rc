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
	Kind      StepKind
	FSettings *FileSettings       `yaml:"file-settings"`
	PSettings *ProcessSettings    `yaml:"proc-settings"`
	CSettings *ConnectionSettings `yaml:"conn-settings"`
}

type ProcessSettings struct {
	Name string
	Args []string
	Wait bool
}

type FileSettings struct {
	Path string
}

type ConnectionSettings struct {
	Source      ConnectionAddress
	Destination ConnectionAddress
	PayloadSize int `yaml:"payload-size"`
	Scheme      string
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
		case FileCreate, FileModify, FileDelete:
			if step.FSettings == nil {
				return rb, errors.Errorf("missing file settings for step %v", i)
			}
		case StartProcess:
			if step.PSettings == nil {
				return rb, errors.Errorf("missing process settings for step %v", i)
			}
		case SendData:
			if step.CSettings == nil {
				return rb, errors.Errorf("missing connection settings for step %v", i)
			}
		}
	}

	return rb, nil
}
