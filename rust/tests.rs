use crate::{ResourceNameError, ResourceTemplate};
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, PartialEq, Eq)]
struct ArtistKey {
    #[serde(rename = "artist_id")]
    id: String,
}

#[test]
fn parse_and_generate_map_roundtrip() {
    let template = ResourceTemplate::new("//music.example.com/artists/{artist_id}")
        .expect("template should compile");
    let parsed = template
        .parse("//music.example.com/artists/radiohead")
        .expect("resource name should parse");
    assert_eq!(
        parsed.get("artist_id"),
        Some(&"radiohead".to_string()),
        "placeholder should be extracted",
    );
    let generated = template.generate(&parsed).expect("values should generate");
    assert_eq!(generated, "//music.example.com/artists/radiohead");
}

#[test]
fn serde_renamed_fields_work() {
    let template = ResourceTemplate::new("//music.example.com/artists/{artist_id}")
        .expect("template should compile");
    let parsed: ArtistKey = template
        .parse_into("//music.example.com/artists/bjork")
        .expect("typed parse should succeed");
    assert_eq!(
        parsed,
        ArtistKey {
            id: "bjork".to_string()
        }
    );
    let generated = template
        .generate_from(&ArtistKey {
            id: "the-cure".to_string(),
        })
        .expect("typed generate should succeed");
    assert_eq!(generated, "//music.example.com/artists/the-cure");
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
    let err = ResourceTemplate::new("//x/{artist-id}")
        .expect_err("invalid placeholder name should error");
    assert!(matches!(err, ResourceNameError::InvalidPlaceholder(_)));
}

#[derive(Debug, Serialize)]
struct ArtistStatus {
    artist_id: String,
    touring: bool,
    rank: u8,
}

#[test]
fn typed_generate_accepts_bool_and_number_scalars() {
    let template = ResourceTemplate::new("//x/{artist_id}/{touring}/{rank}")
        .expect("template should compile");
    let path = template
        .generate_from(&ArtistStatus {
            artist_id: "the-cure".to_string(),
            touring: true,
            rank: 3,
        })
        .expect("typed generation should succeed");
    assert_eq!(path, "//x/the-cure/true/3");
}
