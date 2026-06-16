# `@the-protobuf-project/resourcename` examples

From the **package root** (`typescript/`), build the library first (published `exports` point at `dist/`):

```bash
bun run build
```

Then:

```bash
cd examples
bun install
bun ./index.ts
```

Optional diagnostics:

```bash
RESOURCE_NAME_LOG_LEVEL=debug bun ./index.ts
```
