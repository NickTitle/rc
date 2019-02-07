package fileformat

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Runbook is an aggregtion of steps to run with the agent, to generate host activity and
// corresponding logs
//
type Runbook struct {
	Name string
	// The ordered set of Steps to execute when running the runbook
	Steps []Step
}

// Step is used to define a unit of work executed by the test agent. The "kind" of step defines
// the required settings from among the optional fields included below
//
type Step struct {
	Kind StepKind
	// Required settings for file steps
	FSettings *FileSettings `yaml:"file-settings"`
	// Required settings for process steps
	PSettings *ProcessSettings `yaml:"proc-settings"`
	// Required settings for connection steps
	CSettings *ConnectionSettings `yaml:"conn-settings"`
}

// ProcessSettings are used for controling the behavior of StartProcess steps
//
type ProcessSettings struct {
	// Name of the process to start
	Name string
	// The list of args provided to the process
	Args []string
	// Whether or not to wait for the process to complete before moving on
	Wait bool
}

// FileSettings are used for controling the behavior of File steps
//
type FileSettings struct {
	// Path to the file you intend to create/modify/destroy
	Path string
}

// ConnectionSettings are used for controling the behavior of SendData steps
//
type ConnectionSettings struct {
	// Local address for requests to originate from
	Source ConnectionAddress
	// Destination address for post requests to target
	Destination ConnectionAddress
	// Amount of bytes written to the destination on successful connection
	PayloadSize int `yaml:"payload-size"`
	// HTTP/HTTPS to account for the scheme on the other end
	Scheme string
}

// ConnectionAddress is used to define an address for either side of SendData connections
//
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

// Validate attempts to read a yaml file at the given path, and unmarshal its contents into a runbook
// An unsuccessful attempt to unmarshal will result in a fatal error. Additionally, settings are required
// for each type of step within a runbook. Missing settings for a given step will also result in a fatal error
//
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
