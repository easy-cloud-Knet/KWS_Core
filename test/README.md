# Testing Guide

## Run

```bash
make test       # run all tests
make test-v     # verbose output
```

Run a specific package:

```bash
go test ./vm/parsor/...
go test ./vm/parsor/cloud-init/...
go test ./DomCon/...
```

## Structure

Test files live alongside their target files (`filename_test.go`).

```
vm/parsor/
  util.go
  util_test.go
  cloud-init/
    userConf.go
    userConf_test.go
DomCon/domainList_status/
  cpu_status.go
  cpu_status_test.go
```

## Coverage

| Package | Test File | Coverage |
|---|---|---|
| `vm/parsor` | `util_test.go` | UUID validation, safe path generation |
| `vm/parsor/cloud-init` | `userConf_test.go` | cloud-init YAML generation |
| `DomCon/domainList_status` | `cpu_status_test.go` | vCPU counter concurrency |

## Guidelines

### External Dependencies

- **Filesystem**: use `t.TempDir()` — cleaned up automatically after each test
- **External binaries** (e.g. `mkpasswd`): isolate with `//go:build linux`
- **libvirt connection**: isolate with `//go:build integration`

Run integration tests:
```bash
go test -tags=integration ./...
```

### Patterns

```go
// Table-driven test
func TestFoo(t *testing.T) {
    tests := []struct{ ... }{ ... }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) { ... })
    }
}

// Filesystem test
func TestWriteFile(t *testing.T) {
    dir := t.TempDir()
    // ...
}
```
