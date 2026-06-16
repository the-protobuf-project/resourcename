import Resourcename

// Attribute-macro API — attaches a static `resourcename` namespace to the type.
@Resource("//music.example.com/artists/{artist_id}")
struct Artist {}

@Resource("//music.example.com/artists/{artist_id}/albums/{album_id}")
struct Album {}

print("=== Resource Name Demo ===")

// Template API
let template = try ResourceTemplate("//music.example.com/artists/{artist_id}")
print("\n1. Template API:")
print("   Generated:  \(try template.generate(["artist_id": "bjork"]))")
print("   Parsed:     \(try template.parse("//music.example.com/artists/radiohead"))")

// Macro API
print("\n2. @Resource macro:")
print("   Generated:  \(try Artist.resourcename.generate(["artist_id": "bjork"]))")
print("   Parsed:     \(try Artist.resourcename.parse("//music.example.com/artists/radiohead"))")

// Nested macro
print("\n3. Nested template:")
print(
    "   Generated:  "
        + (try Album.resourcename.generate(["artist_id": "radiohead", "album_id": "in-rainbows"]))
)
print("   Parsed:     \(try Album.resourcename.parse("//music.example.com/artists/the-smiths/albums/the-queen-is-dead"))")
