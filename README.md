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
| Swift          | [`swift/`](swift/)             | `Resourcename` (`@Resource` macro)              |
| C              | [`c/`](c/)                     | `resourcename.h` (`libresourcename`)            |

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
import resourcename

t = resourcename.ResourceTemplate("//music.example.com/artists/{artist_id}")
t.parse("//music.example.com/artists/radiohead")   # {'artist_id': 'radiohead'}
t.generate(artist_id="bjork")                      # "//music.example.com/artists/bjork"

# or the class decorator
@resourcename.resource("//music.example.com/artists/{artist_id}")
class Artist:
    pass
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
import resourcename from "@the-protobuf-project/resourcename";

const t = new resourcename.ResourceTemplate("//music.example.com/artists/{artist_id}");
t.parse("//music.example.com/artists/radiohead"); // { artist_id: "radiohead" }
t.generate({ artist_id: "bjork" });               // "//music.example.com/artists/bjork"

// or a typed class
class Artist extends resourcename.resourceNameBase("//music.example.com/artists/{artist_id}") {}
Artist.Resource.Generate({ artist_id: "bjork" });
```

```bash
cd typescript && bun install && bun run build && bun run check
```

### Swift

```swift
import Resourcename

let t = try ResourceTemplate("//music.example.com/artists/{artist_id}")
try t.parse("//music.example.com/artists/radiohead")  // ["artist_id": "radiohead"]
try t.generate(["artist_id": "bjork"])                // "//music.example.com/artists/bjork"

// or the @Resource attribute macro
@Resource("//music.example.com/artists/{artist_id}")
struct Artist {}
try Artist.resourcename.generate(["artist_id": "bjork"])
```

```bash
swift build && swift run ResourcenameExamples
```

### C

```c
#include "resourcename.h"
#include <stddef.h>

typedef struct { char *id; char *name; } Artist;
static const rn_field ARTIST_FIELDS[] = {
    { "id", offsetof(Artist, id) }, { "name", offsetof(Artist, name) },
};

rn_template *t = NULL;
rn_template_compile("//music.example.com/artists/{id}/{name}", &t);

char *name = NULL;
Artist a = { .id = "ar-42", .name = "Radiohead" };
rn_generate(t, &a, ARTIST_FIELDS, 2, &name);   // "//music.example.com/artists/ar-42/Radiohead"
free(name);
rn_template_free(t);
```

```bash
cd c && make && make run && make test
```

## Repository layout

```text
.
├── resourcename.go # Go public API — re-exports the ./go implementation
├── go.mod          # Go module root: github.com/the-protobuf-project/resourcename
├── go/             # Go implementation package (.../resourcename/go)
├── python/         # uv workspace member: resourcename
├── rust/           # Cargo workspace member: resourcename (+ rust/macros derive crate)
├── typescript/     # Bun/npm workspace member: @the-protobuf-project/resourcename
├── swift/          # Swift package sources: Resourcename (+ ResourcenameMacros)
├── c/              # C library + example + tests (Makefile)
├── Cargo.toml      # Rust workspace root
├── Package.swift   # Swift package manifest (sources under swift/)
├── package.json    # JS workspace root
└── pyproject.toml  # uv workspace root
```

The Go module lives at the repo root (so `go get github.com/the-protobuf-project/resourcename`
resolves), while the implementation stays in [`go/`](go/) and is re-exported by
[`resourcename.go`](resourcename.go).

----
Copyright 2026 The Protobuf Project. Licensed under the Apache License, Version 2.0
