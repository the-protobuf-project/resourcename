use resourcename::Resource;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Resource)]
#[resource_name(template = "//a/{id}")]
#[resource_name(template = "//b/{id}")]
struct ArtistKey {
    id: String,
}

fn main() {
}
