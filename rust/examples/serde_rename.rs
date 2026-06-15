//! Typed parse/generate using serde field rename.

use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//system.com/devices/{device_id}")]
struct DeviceKey {
    #[serde(rename = "device_id")]
    id: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let generated = DeviceKey {
        id: "camera-7".to_string(),
    }
    .generate()?;
    println!("generated: {generated}");

    let parsed = DeviceKey::parse(&generated)?;
    println!("parsed: {parsed:?}");

    Ok(())
}
