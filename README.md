# Park Security Patrol Admin

Park Security Patrol Admin is a deterministic Go backend that models the operational core of a multi-campus patrol system. It uses an in-memory repository and fixed fixtures, so it needs no database, network service, system clock, random source, or platform-specific facility.

## Requirements

- Go 1.24.13
- `CGO_ENABLED=0`

The module only uses the Go standard library.

## Run

```sh
CGO_ENABLED=0 go run ./cmd/patrol-admin
```

The command creates a route and shift, performs ordered checkpoint check-ins, reports and resolves an incident, records supervisor approval, executes a fixed remote verification step, and emits menu, role, dictionary, and CSV export results as JSON.

## Test

```sh
CGO_ENABLED=0 go test -count=1 ./...
```

The suite covers the complete patrol and incident lifecycle, campus data isolation, role-based menus, dictionaries, and deterministic CSV export. It also contains the acceptance regression for an immediately expired remote-step deadline. That regression currently fails by design: the remote step reports `completed` with no error instead of returning `context deadline exceeded` and `timed_out`.

## Packages

- `cmd/patrol-admin`: executable deterministic demonstration
- `internal/app`: application composition
- `internal/domain`: business entities and statuses
- `internal/fixture`: fixed actors, route, and shift
- `internal/service`: patrol workflow, access configuration, export, and remote step
- `internal/store`: campus-scoped in-memory persistence

## Portability

All production paths are pure Go. Builds use no CGO and are suitable for `linux/amd64` and `linux/arm64`. The source tree is not intended to contain compiled output.
