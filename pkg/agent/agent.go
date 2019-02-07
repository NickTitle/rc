package agent

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nicktitle/rc/pkg/fileformat"
)

const modByteLen = 10

type Agent struct {
	logger  log.Logger
	runbook fileformat.Runbook
}

type Event map[string]interface{}

func NewAgent(logger log.Logger, runbook fileformat.Runbook) *Agent {
	return &Agent{
		logger:  logger,
		runbook: runbook,
	}
}

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
