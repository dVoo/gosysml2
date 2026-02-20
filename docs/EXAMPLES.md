# Examples

Runnable examples are in `examples/`.

## Quick Mapping

- `examples/basic/main.go`: basic parse, walk, and typed element queries
- `examples/requirements/main.go`: requirement-focused parsing and traceability fields
- `examples/parallel/main.go`: directory parsing in parallel
- `examples/validation/main.go`: syntax validation patterns
- `examples/visitor/main.go`: custom visitor traversal
- `examples/main.go`: broader showcase entry point

## Run Examples

From repository root:

```bash
go run ./examples/basic
go run ./examples/requirements
go run ./examples/parallel
go run ./examples/validation
go run ./examples/visitor
```

## Which Example to Start With

1. Start with `basic` for first successful parse.
2. Move to `requirements` if your domain includes requirements/verification links.
3. Use `parallel` for multi-file repositories.
4. Use `visitor` when building analyzers/linters/exporters.
