package agent

import (
	"os/exec"
	"strings"

	"github.com/nicktitle/rc/pkg/fileformat"
	"github.com/pkg/errors"
)

func startProcess(name string, args []string) (Event, error) {
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return Event{}, errors.Wrapf(err, "process failed: %s", name)
	}

	if err := cmd.Wait(); err != nil {
		return Event{}, errors.Wrapf(err, "failure waiting for process to complete, %s", name)
	}

	return Event{
		"activity":             fileformat.StartProcess,
		"started_process_name": name,
		"started_process_args": strings.Join(args, " "),
	}, nil
}
