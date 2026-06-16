# resourcename

Google [AIP-122](https://aip.dev/122)-style **resource name** templates with `{placeholder}`
segments — parse a full resource name into its components and generate one from components —
implemented across four languages in a single monorepo.

A template looks like:

```text
//system.com/devices/{device_id}
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

type User struct {
    _    struct{} `resource:"//example.com/users/{id}/{name}"`
    ID   string   `resource:"id"`
    Name string   `resource:"name"`
}

rn, _ := resourcename.MarshalResource(&User{ID: "u42", Name: "Ria"})
// "//example.com/users/u42/Ria"
```

```bash
cd go && go test ./... && go run ./example
```

### Python

```python
from resourcename import resourcename

@resourcename("//system.com/devices/{device_id}")
class Device:
    pass

Device.resourcename.parse("//system.com/devices/router-01")   # {'device_id': 'router-01'}
Device.resourcename.generate(device_id="sensor-22")           # "//system.com/devices/sensor-22"
```

### Rust

```rust
use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//system.com/devices/{device_id}")]
struct DeviceKey {
    device_id: String,
}

let generated = DeviceKey { device_id: "sensor-22".into() }.generate()?;
let parsed = DeviceKey::parse("//system.com/devices/router-01")?;
```

```bash
cargo test
cargo run -p resourcename --example basic
```

### TypeScript

```ts
import { resourceNameBase } from "@the-protobuf-project/resourcename";

class Device extends resourceNameBase("//system.com/devices/{device_id}") {}

Device.Resource.Parse("//system.com/devices/router-01");
Device.Resource.Generate({ device_id: "sensor-22" });
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
