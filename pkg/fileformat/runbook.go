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
	ConnSettings *ConnectionSettings `yaml:"conn-settings,omitempty"`
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
	SpawnProcess StepKind = "SpawnProcess"
	ExfilData    StepKind = "ExfilData"
)

func ExportExample() string {
	rb := Runbook{
		Name: "make-mod-del",
		Steps: []Step{
			Step{
				Kind: FileCreate,
				Path: "~/tmp/foo.txt",
			},
			Step{
				Kind: FileModify,
				Path: "~/tmp/foo.txt",
			},
			Step{
				Kind: FileDelete,
				Path: "~/tmp/foo.txt",
			},
		},
	}

	bytes, err := yaml.Marshal(rb)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func Validate(path string) (Runbook, error) {
	var rb Runbook

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return rb, errors.Wrap(err, "unable to read file")
	}

	if err := yaml.Unmarshal(bytes, &rb); err != nil {
		return rb, errors.New("unable to parse runbook from file")
	}

	return rb, nil
}
