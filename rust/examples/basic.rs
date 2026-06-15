//! Basic derive-driven parse + generate.

use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//system.com/devices/{device_id}")]
struct DeviceKey {
    device_id: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let generated = DeviceKey {
        device_id: "sensor-22".to_string(),
    }
    .generate()?;
    println!("generated: {generated}");

    let parsed = DeviceKey::parse("//system.com/devices/router-01")?;
    println!("parsed: {parsed:?}");

    Ok(())
}
