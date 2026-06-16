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
//! struct ArtistKey {
//!     #[serde(rename = "artist_id")]
//!     id: String,
//! }
//!
//! let template = ResourceTemplate::new("//music.example.com/artists/{artist_id}")?;
//! let parsed: ArtistKey = template.parse_into("//music.example.com/artists/radiohead")?;
//! assert_eq!(parsed.id, "radiohead");
//! # Ok::<(), resourcename::ResourceNameError>(())
//! ```

mod resource_error;
mod resource_template;
mod template_helpers;

pub use macros::Resource;
pub use resource_error::ResourceNameError;
pub use resource_template::ResourceTemplate;

#[cfg(test)]
mod tests;
