## Logging guidelines (syslogger)

Purpose
-------
All application logs must be emitted through the `syslogger` package (zap-based). Do not use `fmt.Println`, `fmt.Printf`, or the standard `log` package for runtime/application logs.

High level rules
----------------
- Use structured logging (zap) and prefer typed fields over printf-style messages.
- Inject `*zap.Logger` into high-level components (for example, `InstHandler`) instead of using global `fmt`/`log` calls.
- Keep the logging depth limited: only the API/handler layer and service layer should produce informational logs. Lower-level libraries should return errors and let callers decide how to log them.
- Avoid logging raw pointers or whole internal structs. Log identifiers (UUID, request_id, name) or use `zap.Any` only when necessary.

Log levels
----------
- `Debug`: development-only details and verbose internal state.
- `Info`: normal operational events (resource created, request accepted).
- `Warn`: unexpected but recoverable situations.
- `Error`: errors that affect the operation or request handling. Include `zap.Error(err)`.
- `Fatal` / `Panic`: only for unrecoverable bootstrap failures.

Standard fields
---------------
When applicable, include these structured fields to make logs searchable and consistent:

- `component` — short name of module (e.g. `vm.creator`).
- `uuid` / `request_id` — tracing identifiers.
- `error` — use `zap.Error(err)`.
- `user` / `tenant` — if multi-tenant context exists.

Examples
--------
Preferred (structured):

```go
logger.Info("vm created",
    zap.String("uuid", uuid),
    zap.String("component", "vm.creator"),
)

logger.Error("failed to create vm",
    zap.String("uuid", uuid),
    zap.Error(err),
)
```

If you need printf-style messages for quick debugging, use the sugared logger sparingly:

```go
logger.Sugar().Infof("vm created: uuid=%s", uuid)
```

Error handling
--------------
- Do not swallow errors. Return errors to caller layers and log them once at a boundary (typically in handlers/API layer) with context.
- When logging errors, add human-friendly prefix/context to indicate where the error occurred.

Do not
------
- Use `fmt.Println` / `fmt.Printf` for application logs.
- Print raw pointer addresses or entire internal objects.

Suggested enforcement
---------------------
- Add a quick code review checklist item to reject PRs that introduce `fmt.Println` or `log.Print*` in non-test code.
- Consider a static check (grep) in CI to detect stray `fmt.Println` usages.

Notes about this document
------------------------
- This file documents recommended logging style and patterns for the codebase. It is intentionally short — if you want, we can add a small example `logger.New()` wrapper and example of how to inject the logger into `InstHandler`.
### How to write log in this project.

All log emmiting from the project should be taken by syslogger(this file. - via zap.logger)

*zap.logger should be embedded to Insthandler and depth of logging should not passed more dipper than service layer.
which means, we don't print non-informational data for every instruction.

info.log, info.warn, info.error will be only that we concern.
Info.panic only occures when there's critical error from bootstrap.


All error should be passed with err argument and must be reported at the end of the api handler. 
For debugging we suggest adding prefixing comprehansible information of where the error occured.

```bash  
e.g) Error: From createvm: instanceCon: CheckInstancePersistance: Instance with the same name already exists.
```

