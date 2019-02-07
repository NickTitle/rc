package agent

import (
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log/level"
)

// decorateEventAndEmit takes a unix timestamp and an Event from a fileformat.Step execution,
// and adds the common fields that should be output by all types of activity
//
func (a *Agent) decorateEventAndEmit(unixtime int64, event Event) {
	// by translating to a slice here, we guarantee key alignment in the logs
	// since indexing is reliably ordered, and maps offer no such guarantee

	var username string
	if u, err := user.LookupId(strconv.Itoa(os.Getuid())); err == nil {
		username = u.Name
	}

	// note that we are styling this like log output for clarity (meaning "keyname" and
	// "key val" on the same line), though typical syntax would be one argument per line
	eventPairs := []interface{}{
		"ts", unixtime,
		"pid", os.Getpid(),
		"process_name", os.Args[0],
		"command", strings.Join(os.Args, " "),
		"username", username,
	}

	for k, v := range event {
		eventPairs = append(eventPairs, k, v)
	}

	level.Info(a.logger).Log(eventPairs...)
}
