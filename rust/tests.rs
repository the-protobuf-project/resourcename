use crate::{ResourceNameError, ResourceTemplate};
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, PartialEq, Eq)]
struct DeviceKey {
    #[serde(rename = "device_id")]
    id: String,
}

#[test]
fn parse_and_generate_map_roundtrip() {
    let template =
        ResourceTemplate::new("//system.com/devices/{device_id}").expect("template should compile");
    let parsed = template
        .parse("//system.com/devices/router-01")
        .expect("resource name should parse");
    assert_eq!(
        parsed.get("device_id"),
        Some(&"router-01".to_string()),
        "placeholder should be extracted",
    );
    let generated = template.generate(&parsed).expect("values should generate");
    assert_eq!(generated, "//system.com/devices/router-01");
}

#[test]
fn serde_renamed_fields_work() {
    let template =
        ResourceTemplate::new("//system.com/devices/{device_id}").expect("template should compile");
    let parsed: DeviceKey = template
        .parse_into("//system.com/devices/sensor-22")
        .expect("typed parse should succeed");
    assert_eq!(
        parsed,
        DeviceKey {
            id: "sensor-22".to_string()
        }
    );
    let generated = template
        .generate_from(&DeviceKey {
            id: "camera-7".to_string(),
        })
        .expect("typed generate should succeed");
    assert_eq!(generated, "//system.com/devices/camera-7");
}

#[test]
fn rejects_duplicate_placeholders() {
    let err =
        ResourceTemplate::new("//x/{id}/y/{id}").expect_err("duplicate placeholders should error");
    assert!(
        matches!(err, ResourceNameError::DuplicatePlaceholders(_)),
        "expected duplicate placeholder error",
    );
}

#[test]
fn rejects_unbalanced_braces() {
    let err = ResourceTemplate::new("//x/{id").expect_err("unbalanced braces should error");
    assert!(matches!(err, ResourceNameError::UnbalancedBraces));
}

#[test]
fn rejects_invalid_placeholder_name() {
    let err = ResourceTemplate::new("//x/{device-id}")
        .expect_err("invalid placeholder name should error");
    assert!(matches!(err, ResourceNameError::InvalidPlaceholder(_)));
}

#[derive(Debug, Serialize)]
struct DeviceStatus {
    device_id: String,
    online: bool,
    retries: u8,
}

#[test]
fn typed_generate_accepts_bool_and_number_scalars() {
    let template = ResourceTemplate::new("//x/{device_id}/{online}/{retries}")
        .expect("template should compile");
    let path = template
        .generate_from(&DeviceStatus {
            device_id: "camera-7".to_string(),
            online: true,
            retries: 3,
        })
        .expect("typed generation should succeed");
    assert_eq!(path, "//x/camera-7/true/3");
}
