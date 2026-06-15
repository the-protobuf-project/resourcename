//! Common validation errors (missing key, extra key, slash in value).

use resourcename::ResourceTemplate;
use std::collections::BTreeMap;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let template = ResourceTemplate::new("//system.com/devices/{device_id}")?;

    let missing = BTreeMap::new();
    println!("missing => {:?}", template.generate(&missing));

    let mut extra = BTreeMap::new();
    extra.insert("device_id".to_string(), "router-01".to_string());
    extra.insert("other".to_string(), "x".to_string());
    println!("extra => {:?}", template.generate(&extra));

    let mut slash = BTreeMap::new();
    slash.insert("device_id".to_string(), "region/a".to_string());
    println!("slash => {:?}", template.generate(&slash));

    Ok(())
}
