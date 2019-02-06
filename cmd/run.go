package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/nicktitle/rc/pkg/agent"
	"github.com/nicktitle/rc/pkg/fileformat"
)

func main() {
	runbookPath := flag.String("runbook", "", "path to the runbook yaml you wish to run")
	flag.Parse()

	if *runbookPath == "" {
		flag.Usage()
		os.Exit(0)
	}

	rb, err := fileformat.Validate(*runbookPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	a := agent.NewAgent(logger, rb)
	a.Run()
}
