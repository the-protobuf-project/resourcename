# `resourcename`

Resource-name parsing and generation for templates like:

`//system.com/devices/{device_id}`

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
#[resource_name(template = "//system.com/devices/{device_id}")]
struct DeviceKey {
    device_id: String,
}

let generated = DeviceKey {
    device_id: "sensor-22".to_string(),
}
.generate()?;
let parsed = DeviceKey::parse("//system.com/devices/router-01")?;
assert_eq!(parsed.device_id, "router-01");
# Ok::<(), resourcename::ResourceNameError>(())
```

## Serde rename usage

```rust
use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//system.com/devices/{device_id}")]
struct DeviceKey {
    device_id: String,
}

let generated = DeviceKey {
    device_id: "camera-7".to_string(),
}
.generate()?;

let parsed = DeviceKey::parse(&generated)?;
assert_eq!(parsed.device_id, "camera-7");
# Ok::<(), resourcename::ResourceNameError>(())
```

## Manual map mode (advanced)

```rust
use resourcename::ResourceTemplate;
use std::collections::BTreeMap;

let t = ResourceTemplate::new("//system.com/devices/{device_id}")?;
let mut values = BTreeMap::new();
values.insert("device_id".to_string(), "sensor-22".to_string());
let generated = t.generate(&values)?;
assert_eq!(generated, "//system.com/devices/sensor-22");
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
