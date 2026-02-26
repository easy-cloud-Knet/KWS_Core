# CLAUDE.md — KWS_Core

## Project Overview

KWS_Core is a personal cloud management system that wraps libvirt (QEMU/KVM) with a REST API for VM lifecycle management.

- **Module**: `github.com/easy-cloud-Knet/KWS_Core`
- **Language**: Go 1.24.0
- **Key dependencies**: `libvirt.org/go/libvirt`, `go.uber.org/zap`, `github.com/shirou/gopsutil`, `gopkg.in/yaml.v3`

## Build & Run

```bash
make conf          # First-time setup (installs Go, dependencies, monitoring)
make build         # Compile → produces KWS_Core binary in project root
make run           # Build + execute locally
make clean         # Remove binary
```

**Deploy**: `sudo systemctl start kws_core`
**Logs**: `/var/log/kws/{YYYYMMDD}.log` + stdout

## Architecture

```
HTTP Server (:8080)
    ↓
API Handlers (api/)         — methods on InstHandler, JSON request/response
    ↓
Service Layer (vm/service/) — creation, termination, status, snapshot
    ↓
Domain Controller (DomCon/) — in-memory VM registry, thread-safe
    ↓
Libvirt (QEMU/KVM)
```

**Key directories**:
| Directory | Purpose |
|-----------|---------|
| `api/` | HTTP handlers, request/response types, generic `BaseResponse[T]` |
| `server/` | HTTP mux setup, route registration, logger middleware |
| `DomCon/` | In-memory domain map, mutex-guarded, domain seeking/caching |
| `DomCon/domain_status/` | CPU/memory accounting with atomic ops |
| `vm/service/creation/` | VM creation via `VMCreator` interface + factory |
| `vm/service/termination/` | Shutdown/deletion via `DomainTermination`/`DomainDeletion` interfaces |
| `vm/service/status/` | Host and VM metrics, `DataTypeHandler` router |
| `vm/service/snapshot/` | Snapshot lifecycle management |
| `vm/parsor/` | XML config generation, cloud-init YAML parsing |
| `error/` | Custom `VirError` types, `ErrorDescriptor`, `ErrorGen`/`ErrorJoin` |
| `logger/` | zap-based structured logger, HTTP middleware |
| `net/` | Network configuration parsing |

**Entry point**: `main.go` → init logger → create DomListControl → connect libvirt → load domains → start server

## Code Patterns (must follow)

### Handler pattern
All HTTP endpoints are receiver methods on `InstHandler` (`api/type.go`). Standard flow:
1. Create param struct + `ResponseGen[T](message)`
2. Decode body with `HttpDecoder(r, param)`
3. Validate / retrieve domain
4. Call service layer
5. Return `resp.ResponseWriteOK(w, data)` or `resp.ResponseWriteErr(w, err, status)`

### Factory pattern
Object creation uses `*Factory()` functions. Examples:
- `DomSeekUUIDFactory()` — domain lookup
- `DomainTerminatorFactory()` / `DomainDeleterFactory()` — shutdown/delete
- `LocalConfFactory()` / `LocalCreatorFactory()` — VM creation
- `DomListConGen()` — domain list controller

### Interface-first design
Define interfaces before implementations:
- `VMCreator` — `CreateVM() (*libvirt.Domain, error)`
- `DomainTermination` — `TerminateDomain() (*libvirt.Domain, error)`
- `DomainDeletion` — `DeleteDomain() (*libvirt.Domain, error)`
- `DataTypeHandler` — `GetInfo(*domCon.Domain) error`
- `DomainSeeker` — `ReturnDomain() (*Domain, error)`

### Thread safety
- `sync.Mutex` on `DomListControl` (list-level) and `Domain` (domain-level)
- `atomic.AddInt64` for CPU counter operations in `DomainListStatus`
- Lock → defer Unlock pattern throughout

### Error handling
Use the custom error package (`error/Error.go`), not raw `errors.New()`:
```go
// Create a typed error
virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("detail: %w", err))

// Chain errors with context
virerr.ErrorJoin(existingErr, fmt.Errorf("additional context"))
```
Predefined error constants: `NoSuchDomain`, `DomainGenerationError`, `LackCapacityCPU`, `LackCapacityRAM`, `InvalidUUID`, `InvalidParameter`, etc.

### VM storage convention
All VM artifacts live in `/var/lib/kws/{uuid}/` — disk images (`{uuid}.qcow2`), cloud-init ISO (`cidata.iso`), config files.

## Logging Rules

- Use `*zap.Logger` injected via `InstHandler`. **Never** use `fmt.Println`, `fmt.Printf`, or `log.*` for runtime logs.
- Log only at **API/handler** and **service** layers. Lower layers return errors; callers decide how to log.
- **Levels**: Debug (dev detail), Info (operations), Warn (recoverable), Error (with `zap.Error(err)`), Fatal/Panic (bootstrap only).
- **Standard fields**: `component`, `uuid`, `request_id`, `error`.
- Prefix error logs with where the error occurred for debugging context.
- Never log raw pointers or entire internal structs.

```go
// Good
logger.Error("failed to create vm", zap.String("uuid", uuid), zap.Error(err))

// Bad
fmt.Println("error:", err)
```

## Behavioral Guidelines

### Think before coding
- State assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them — don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.

### Simplicity first
- Minimum code that solves the problem. No speculative features.
- No abstractions for single-use code. No "flexibility" that wasn't requested.
- No error handling for impossible scenarios.

### Surgical changes
- Touch only what you must. Don't "improve" adjacent code or formatting.
- Match existing style, even if you'd do it differently.
- Remove imports/variables/functions that YOUR changes made unused.
- Every changed line should trace directly to the request.

### Goal-driven execution
- Transform tasks into verifiable goals before implementing.
- For multi-step tasks, state a brief plan with verification steps.
