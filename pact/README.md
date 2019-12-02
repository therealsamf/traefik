traefik `/` pact
================

[Pact.io](https://pact.io) testing for [traefik](https://traefik.io)

# Usage

See [pact-go](https://github.com/pact-foundation/pact-go/) [installation](https://github.com/pact-foundation/pact-go/#installation) to ensure pact is set up for golang.

From within the `pact` directory, run

```terminal
go1.13 test -v ./...
```

Logs should be written to the `logs` subdirectory and pacts in the `pacts` subdirectory.
