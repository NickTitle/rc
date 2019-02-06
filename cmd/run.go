package main

import (
	"fmt"

	"github.com/nicktitle/rc/pkg/agent"
	"github.com/nicktitle/rc/pkg/fileformat"
)

func main() {
	fmt.Println(fileformat.ExportExample())
	agent.Start()
}
