# rc - test agent for host activity generation

### usage
The code must be invoked with the `-runbook` flag, and a path to a valid runbook. a collection of examples are provided [here](examples/runbooks). On macOS, feel free to use the compiled binary in build/agent. Otherwise, use `go run cmd/run.go` as your entrypoint

ex:
```
./build/agent -runbook examples/runbooks/start-proc.yml
<or>
go run cmd/run.go -runbook examples/runbooks/make-mod-del.yml
```

### building
Interacting with the repo is done via makefile, supporting the following commands.
- `make deps` -> install dependencies
- `make test` -> run test suite
- `make agent` -> compile the binary (macOS only for now)


### runbooks
`Runbooks` are yaml files which collect `Steps` that the agent will run when started with the given runbook. Steps are executed in order, outputting structured logs at each step. The runbook will stop running and emit an error log on exception in any step.

The spec for runbooks is fully documented in [runbook.go](pkg/fileformat/runbook.go)

### extra
For testing the [send data runbook](examples/runbooks/send-data.yml), a [blackhole server](examples/blackhole.go) is provided, to drain activity to if you don't have an endpoint configured already. The server accepts a port argument if you wish to modify the runbook and set the server up to listen somewhere else

### ideas for more work
- the step executor things could be easily abstracted behind a `Stepper` interface, which would implement the `Step` function, that accepts a `fileformat.Step`, and returns an `Event, error` pair just like now. this would be cool, because then we could mock each of the step implementations to test the agent (even though it's just a dumb pipe to those functions now)

- compiling binaries for multiple platforms would be straightforward, but i figured that wasn't the best place to spend my time. it'd be great to call `make agent`, and have it doing `make agent-darwin`, `make agent-windows`, etc under the hood
