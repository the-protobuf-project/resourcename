//! Basic derive-driven parse + generate.

use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//music.example.com/artists/{artist_id}")]
struct ArtistKey {
    artist_id: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let generated = ArtistKey {
        artist_id: "bjork".to_string(),
    }
    .generate()?;
    println!("generated: {generated}");

    let parsed = ArtistKey::parse("//music.example.com/artists/radiohead")?;
    println!("parsed: {parsed:?}");

    Ok(())
}
