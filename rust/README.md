# `resourcename`

Resource-name parsing and generation for templates like:

`//music.example.com/artists/{artist_id}`

The crate supports:

- derive-based parse/generate API via `#[derive(Resource)]` (recommended)
- typed parse/generate via `serde`
- map-based parse/generate via `ResourceTemplate` (advanced/manual mode)

## Add to Cargo.toml

```toml
[dependencies]
resourcename = { path = "../resourcename" }
serde = { version = "1", features = ["derive"] }
```

## Basic usage (recommended)

```rust
use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//music.example.com/artists/{artist_id}")]
struct ArtistKey {
    artist_id: String,
}

let generated = ArtistKey {
    artist_id: "bjork".to_string(),
}
.generate()?;
let parsed = ArtistKey::parse("//music.example.com/artists/radiohead")?;
assert_eq!(parsed.artist_id, "radiohead");
# Ok::<(), resourcename::ResourceNameError>(())
```

## Serde rename usage

```rust
use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//music.example.com/artists/{artist_id}")]
struct ArtistKey {
    artist_id: String,
}

let generated = ArtistKey {
    artist_id: "the-cure".to_string(),
}
.generate()?;

let parsed = ArtistKey::parse(&generated)?;
assert_eq!(parsed.artist_id, "the-cure");
# Ok::<(), resourcename::ResourceNameError>(())
```

## Manual map mode (advanced)

```rust
use resourcename::ResourceTemplate;
use std::collections::BTreeMap;

let t = ResourceTemplate::new("//music.example.com/artists/{artist_id}")?;
let mut values = BTreeMap::new();
values.insert("artist_id".to_string(), "bjork".to_string());
let generated = t.generate(&values)?;
assert_eq!(generated, "//music.example.com/artists/bjork");
# Ok::<(), resourcename::ResourceNameError>(())
```

## Run examples

```bash
cargo run -p resourcename --example basic
cargo run -p resourcename --example serde_rename
cargo run -p resourcename --example nested
cargo run -p resourcename --example error_cases
```

## Limitations

- Placeholders are one path segment each (`[^/]+`), so generated values cannot contain `/`.
- Placeholder names must match `[A-Za-z_][A-Za-z0-9_]*`.
- Typed generation accepts scalar serde values (`string`, `bool`, `number`) only.
