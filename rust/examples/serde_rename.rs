//! Typed parse/generate using serde field rename.

use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//music.example.com/artists/{artist_id}")]
struct ArtistKey {
    #[serde(rename = "artist_id")]
    id: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let generated = ArtistKey {
        id: "the-cure".to_string(),
    }
    .generate()?;
    println!("generated: {generated}");

    let parsed = ArtistKey::parse(&generated)?;
    println!("parsed: {parsed:?}");

    Ok(())
}
