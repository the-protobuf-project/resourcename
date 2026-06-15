//! Google-style resource-name templates with serde-powered typed parsing and generation.
//!
//! # Highlights
//! - Placeholder segments use `{field_name}` and match exactly one path segment.
//! - Supports both map-based and typed (`serde`) parse/generate flows.
//! - Provides derive macro support (`Resource`) for ergonomic model methods.
//!
//! # Example
//! ```rust
//! use resourcename::ResourceTemplate;
//! use serde::{Deserialize, Serialize};
//!
//! #[derive(Debug, Deserialize, Serialize)]
//! struct DeviceKey {
//!     #[serde(rename = "device_id")]
//!     id: String,
//! }
//!
//! let template = ResourceTemplate::new("//system.com/devices/{device_id}")?;
//! let parsed: DeviceKey = template.parse_into("//system.com/devices/router-01")?;
//! assert_eq!(parsed.id, "router-01");
//! # Ok::<(), resourcename::ResourceNameError>(())
//! ```

mod resource_error;
mod resource_template;
mod template_helpers;

pub use macros::Resource;
pub use macros::ResourceName;
#[doc(hidden)]
pub use macros::Resourcename;
pub use resource_error::ResourceNameError;
pub use resource_template::ResourceTemplate;

#[cfg(test)]
mod tests;
