# Zero Hunger Contracts

Canonical protobuf contracts shared by Zero Hunger services.

## Generate Go bindings

```bash
buf lint
buf generate
go test ./...
```

The `proto/` directory is the source of truth. Generated bindings belong under `gen/` and are versioned with semantic tags.
