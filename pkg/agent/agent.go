package agent

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nicktitle/rc/pkg/fileformat"
)

const modByteLen = 10

// Agent is the abstract type that runs runbooks and outputs logs about the steps inside
//
type Agent struct {
	logger  log.Logger
	runbook fileformat.Runbook
}

// Event is used to collect attributes about different steps, for eventual passing to an agent's logger
//
type Event map[string]interface{}

// NewAgent returns an Agent with a logger and runbook ready to go
//
func NewAgent(logger log.Logger, runbook fileformat.Runbook) *Agent {
	return &Agent{
		logger:  logger,
		runbook: runbook,
	}
}

// Run will attempt to run the runbook for an instance of Agent, exiting
// gracefully if something goes wrong, and logging the exception
//
func (a *Agent) Run() {
	for _, step := range a.runbook.Steps {
		if err := a.runStep(step); err != nil {
			level.Error(a.logger).Log("msg", "step failed", "err", err)
			os.Exit(0)
		}
	}
}

func (a *Agent) runStep(step fileformat.Step) error {
	var (
		event Event
		err   error
	)
	start := time.Now().Unix()

	switch step.Kind {
	case fileformat.FileCreate:
		event, err = createFile(*step.FSettings)
	case fileformat.FileModify:
		event, err = modifyFile(*step.FSettings)
	case fileformat.FileDelete:
		event, err = deleteFile(*step.FSettings)
	case fileformat.StartProcess:
		event, err = startProcess(*step.PSettings)
	case fileformat.SendData:
		event, err = sendData(*step.CSettings)
	}

	if err != nil {
		return err
	}

	a.decorateEventAndEmit(start, event)
	return nil
}
