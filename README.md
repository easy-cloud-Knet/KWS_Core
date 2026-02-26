# KWS_Core

A VM lifecycle management REST API built on libvirt (QEMU/KVM). Provides HTTP endpoints for creating, booting, shutting down, deleting, and snapshotting virtual machines.

## Prerequisites

- **Linux** with systemd
- **Go** 1.24+
- **libvirt** and **QEMU/KVM**
- **Open vSwitch** (for network bridging)
- **genisoimage** (cloud-init ISO generation)
- **mkpasswd** (VM user password hashing)

## Quick Start

```bash
make conf                        # Install Go, dependencies, monitoring agents
make build                       # Compile binary
sudo systemctl start kws_core    # Start the service
```

## Build Commands

| Command      | Description                                          |
|--------------|------------------------------------------------------|
| `make conf`  | First-time setup: installs Go, node_exporter, promtail, configures libvirt hooks |
| `make build` | Compile to `KWS_Core` binary in project root         |
| `make run`   | Build and execute locally                            |
| `make clean` | Remove the binary                                    |

## Architecture

```
HTTP Server (:8080)
    |
API Handlers (api/)          -- methods on InstHandler, JSON request/response
    |
Service Layer (vm/service/)  -- creation, termination, status, snapshot
    |
Domain Controller (DomCon/)  -- in-memory VM registry, thread-safe
    |
Libvirt (QEMU/KVM)
```

## API Endpoints

| Method | Path                     | Description                        |
|--------|--------------------------|------------------------------------|
| POST   | `/createVM`              | Create a new VM from a base image  |
| POST   | `/BOOTVM`                | Boot an existing VM                |
| POST   | `/forceShutDownUUID`     | Force shutdown a VM by UUID        |
| POST   | `/DeleteVM`              | Delete a VM and its artifacts      |
| GET    | `/getStatusUUID`         | Get status of a specific VM        |
| GET    | `/getStatusHost`         | Get host system metrics            |
| GET    | `/getInstAllInfo`        | Get detailed info for all VMs      |
| GET    | `/getAllUUIDs`            | List all VM UUIDs                  |
| GET    | `/getAll-uuidstatusList` | List all VMs with their states     |
| POST   | `/CreateSnapshot`        | Create a VM snapshot               |
| GET    | `/ListSnapshots`         | List snapshots for a VM            |
| POST   | `/RevertSnapshot`        | Revert a VM to a snapshot          |
| POST   | `/DeleteSnapshot`        | Delete a snapshot                  |

## Configuration

| Item          | Value                              |
|---------------|------------------------------------|
| Listen port   | `8080`                             |
| VM storage    | `/var/lib/kws/{uuid}/`             |
| Log directory | `/var/log/kws/`                    |
| Log files     | `{YYYYMMDD}.log` (daily rotation)  |

## Logging

Structured logging via [zap](https://github.com/uber-go/zap). Logs are written to `/var/log/kws/{YYYYMMDD}.log` and stdout. See [`logger/Readme.md`](logger/Readme.md) for conventions and usage guidelines.

## Project Structure

```
KWS_Core/
  main.go                    -- Entry point
  api/                       -- HTTP handlers and request/response types
  server/                    -- HTTP mux setup and route registration
  DomCon/                    -- In-memory domain registry (mutex-guarded)
    domain_status/           -- CPU/memory accounting with atomic ops
  vm/
    service/
      creation/              -- VM creation (VMCreator interface, factory)
      termination/           -- Shutdown and deletion
      status/                -- Host and VM metrics
      snapshot/              -- Snapshot lifecycle
    parsor/                  -- XML config and cloud-init YAML generation
  net/                       -- Network configuration parsing
  error/                     -- Custom VirError types and error chaining
  logger/                    -- zap logger setup and HTTP middleware
  build/                     -- Setup scripts (go.sh, autoconfig.sh, etc.)
```

## Contributing

See [`CLAUDE.md`](CLAUDE.md) for code conventions, patterns, and architectural guidelines.

## License

MIT License. See [`LICENSE`](LICENSE) for details.
