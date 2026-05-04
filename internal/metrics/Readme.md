## matrics

Purpose
-------
Manages registration of custom Prometheus collectors.
`Initialize()` creates a registry and enrolls all collectors at once.

Adding a new collector
----------------------
1. Create a package under `internal/matrics/<name>/`.
2. Implement the `Matrics` interface.

```go
type Matrics interface {
    Enroll(prometheus.Registerer) error
}
```

3. Inside `Enroll`, create and register your GaugeVec/CounterVec, then start a collection goroutine.

```go
func (c *Collector) Enroll(reg prometheus.Registerer) error {
    c.gauge = prometheus.NewGaugeVec(...)
    if err := reg.Register(c.gauge); err != nil {
        return err
    }
    go c.collect()
    return nil
}
```

4. Add the collector to `matricsList` in `init.go`.

```go
var matricsList []Matrics = []Matrics{
    &ping.Collector{},
    &yourpkg.Collector{}, // add here
}
```

Rules
-----
- Return `error` from `Enroll` on failure. `Initialize` handles the panic.
- Errors inside goroutines should skip the current cycle and wait for the next — do not crash the server.
- Use `Register` + error return instead of `MustRegister`.
- Follow the metric naming convention: `<subsystem>_<name>_<unit>` (e.g. `ping_rtt_seconds`).
