//! Nested struct example using `#[serde(flatten)]`.

use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize)]
struct Artist {
    artist_id: String,
}

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//music.example.com/artists/{artist_id}/albums/{album_id}")]
struct AlbumRef {
    #[serde(flatten)]
    artist: Artist,
    album_id: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let generated = AlbumRef {
        artist: Artist {
            artist_id: "radiohead".to_string(),
        },
        album_id: "in-rainbows".to_string(),
    }
    .generate()?;
    println!("generated: {generated}");

    let parsed = AlbumRef::parse("//music.example.com/artists/the-smiths/albums/the-queen-is-dead")?;
    println!("parsed: {parsed:?}");

    Ok(())
}
