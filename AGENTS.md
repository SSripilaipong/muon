# Project notes for future tasks

- Muon is a control plane for running the "muto" language across a cluster. The server exposes an HTTP API and coordinates execution through a raft-backed event source and a runner that spawns goroutines for compiled objects.
- The `server/runner` package mirrors request flows in separate files. Each file starts with the external API (service or controller method) followed by the internal processor handlers in the order calls are triggered at runtime. Keep this ordering when extending the package.
- The runner talks to the in-memory event source (`server/eventsource`) through method calls rather than actor messages. Append and commit requests originate in the coordinator and are proxied by the runner.
- Run `go test ./...` after making changes to validate behaviour, and always format Go sources with `gofmt` before committing.

