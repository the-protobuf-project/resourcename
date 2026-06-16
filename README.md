# resourcename

Google [AIP-122](https://aip.dev/122)-style **resource name** templates with `{placeholder}`
segments — parse a full resource name into its components and generate one from components —
implemented across four languages in a single monorepo.

A template looks like:

```text
//music.example.com/artists/{artist_id}
```

Each `{placeholder}` matches exactly one path segment (`[^/]+`), so generated values may not
contain `/`.

## Implementations

| Language       | Directory                      | Canonical import / package                      |
| -------------- | ------------------------------ | ----------------------------------------------- |
| Go             | [`go/`](go/)                   | `github.com/the-protobuf-project/resourcename`  |
| Python         | [`python/`](python/)           | `resourcename`                                  |
| Rust           | [`rust/`](rust/)               | `resourcename` crate (`#[derive(Resource)]`)    |
| TypeScript     | [`typescript/`](typescript/)   | `@the-protobuf-project/resourcename`            |

### Go

```go
import "github.com/the-protobuf-project/resourcename"

type Artist struct {
    _    struct{} `resource:"//music.example.com/artists/{id}/{name}"`
    ID   string   `resource:"id"`
    Name string   `resource:"name"`
}

rn, _ := resourcename.MarshalResource(&Artist{ID: "ar-42", Name: "Radiohead"})
// "//music.example.com/artists/ar-42/Radiohead"
```

```bash
cd go && go test ./... && go run ./example
```

### Python

```python
from resourcename import resourcename

@resourcename("//music.example.com/artists/{artist_id}")
class Artist:
    pass

Artist.resourcename.parse("//music.example.com/artists/radiohead")   # {'artist_id': 'radiohead'}
Artist.resourcename.generate(artist_id="bjork")                      # "//music.example.com/artists/bjork"
```

### Rust

```rust
use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//music.example.com/artists/{artist_id}")]
struct ArtistKey {
    artist_id: String,
}

let generated = ArtistKey { artist_id: "bjork".into() }.generate()?;
let parsed = ArtistKey::parse("//music.example.com/artists/radiohead")?;
```

```bash
cargo test
cargo run -p resourcename --example basic
```

### TypeScript

```ts
import { resourceNameBase } from "@the-protobuf-project/resourcename";

class Artist extends resourceNameBase("//music.example.com/artists/{artist_id}") {}

Artist.Resource.Parse("//music.example.com/artists/radiohead");
Artist.Resource.Generate({ artist_id: "bjork" });
```

```bash
cd typescript && bun install && bun run build && bun run check
```

## Repository layout

```text
.
├── go/          # Go module: github.com/the-protobuf-project/resourcename
├── python/      # uv workspace member: resourcename
├── rust/        # Cargo workspace member: resourcename (+ rust/macros derive crate)
├── typescript/  # Bun/npm workspace member: @the-protobuf-project/resourcename
├── Cargo.toml   # Rust workspace root
├── go.work      # Go workspace
├── package.json # JS workspace root
└── pyproject.toml # uv workspace root
```

----
Copyright 2026 The Protobuf Project. Licensed under the Apache License, Version 2.0
