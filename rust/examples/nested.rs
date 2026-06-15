//! Nested struct example using `#[serde(flatten)]`.

use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize)]
struct Tenant {
    tenant_id: String,
}

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//system.com/tenants/{tenant_id}/devices/{device_id}")]
struct DeviceRef {
    #[serde(flatten)]
    tenant: Tenant,
    device_id: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let generated = DeviceRef {
        tenant: Tenant {
            tenant_id: "alpha".to_string(),
        },
        device_id: "camera-7".to_string(),
    }
    .generate()?;
    println!("generated: {generated}");

    let parsed = DeviceRef::parse("//system.com/tenants/beta/devices/lidar-42")?;
    println!("parsed: {parsed:?}");

    Ok(())
}
