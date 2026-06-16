# `@the-protobuf-project/resourcename`

TypeScript helpers for **resource name** strings with `{placeholder}` segments (Google-style API paths): compile templates, **parse** a full name into segments, **generate** a name from values, and optional **stage-3 `@resourceName`** class decorator or **`resourceNameBase`** `extends` for typed `Artist.Resource.Parse` / `Generate`.

## Requirements

- **Runtime:** Node 20+, Bun, or any JS environment your bundler targets.
- **Peer:** `typescript` ≥ 5 &lt; 7 (for types and decorator semantics in source consumers).

## Install

Published **compiled** output lives under `dist/` (see `tsconfig.build.json` and `package.json` `exports`). For monorepo `examples/` with `file:..`, run **`bun run build`** in this folder first.

## Quick start

```ts
import { resourceNameBase } from "@the-protobuf-project/resourcename";

class Artist extends resourceNameBase("//music.example.com/artists/{artist_id}") {}

Artist.Resource.Parse("//music.example.com/artists/radiohead");
Artist.Resource.Generate({ artist_id: "bjork" });
```

**Decorator** (add `declare static readonly Resource: ClassResource` for typings):

```ts
import {
  resourceName,
  type ClassResource,
} from "@the-protobuf-project/resourcename";

@resourceName("//music.example.com/artists/{artist_id}")
class Artist {
  declare static readonly Resource: ClassResource;
}
```

## Logging

Diagnostics use **`ResourceNameLogger`** and **`PackageLogger`** (`shared/`). They use `console.log` / `console.warn` / `console.error` only — **no Node-only APIs**, so the same code runs in **browsers** and **Node/Bun**.

- **Default threshold:** `debug` when `NODE_ENV !== "production"`, otherwise **`info`**.
- **Env** (from `process.env` when present; bundlers may inject for browser builds):

| Variable                             | Purpose                                            |
| ------------------------------------ | -------------------------------------------------- |
| `RESOURCE_NAME_LOG_LEVEL`       | `silent` \| `debug` \| `info` \| `warn` \| `error` |
| `RESOURCE_NAME_COLOR`           | `auto` \| `always` \| `never`                      |
| `RESOURCE_NAME_PACKAGE_NAME`    | Override package name in log lines                 |
| `RESOURCE_NAME_PACKAGE_VERSION` | Override version in log lines                      |

**Levels:** mostly **`debug`** (template compile, parse/generate success, decorator attach). **`info`:** two one-time lines — runtime ready on first `ResourceTemplate`, and “class decorator pipeline” on first `@resourceName` application. **`warn`:** bad templates, parse mismatch, invalid `Generate` inputs. **`error`:** reserved for the same sink as hard failures (optional future use).

```ts
import {
  LogLevel,
  ResourceNameLogger,
} from "@the-protobuf-project/resourcename";

ResourceNameLogger.configure({ minLevel: LogLevel.WARN });
```

## Build

```bash
bun run build
bun run check
```

## API docs (TypeDoc + TSDoc)

Comments use **TSDoc**; **TypeDoc** turns them into browsable HTML:

```bash
bun install
bun run docs        # writes docs/api/ (gitignored)
bun run docs:open   # generate and open docs/api/index.html
```

`typedoc.json` points at `index.ts` and this package `tsconfig.json`. **`tsdoc.json`** (if present) configures editors for `@remarks` / `@throws`.
