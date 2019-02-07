package agent

import (
	"os/exec"
	"strings"

	"github.com/nicktitle/rc/pkg/fileformat"
	"github.com/pkg/errors"
)

func startProcess(settings fileformat.ProcessSettings) (Event, error) {
	cmd := exec.Command(settings.Name, settings.Args...)
	if err := cmd.Start(); err != nil {
		return Event{}, errors.Wrapf(err, "failed to start process: %s", settings.Name)
	}
	if settings.Wait {
		if err := cmd.Wait(); err != nil {
			return Event{}, errors.Wrapf(err, "failure waiting for process to complete, %s", settings.Name)
		}
	}

	return Event{
		"activity":             fileformat.StartProcess,
		"started_process_name": settings.Name,
		"started_process_args": strings.Join(settings.Args, " "),
	}, nil
}
