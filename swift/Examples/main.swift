import Resourcename

// MARK: - Models
//
// Like the Go example's tagged structs, each stored property maps to a template
// placeholder of the same name. Encodable/Decodable drives the typed
// `generate(from:)` / `parse(_:as:)` API.

/// A basic resource: `//music.example.com/artists/{id}/{name}`.
struct Artist: Codable {
    let id: String
    let name: String
}

/// A multi-segment resource:
/// `//music.example.com/labels/{label}/artists/{artist}/albums/{album}`.
struct Release: Codable {
    let label: String
    let artist: String
    let album: String
}

/// Scalar (non-`String`) fields are stringified when generating.
struct TrackInfo: Encodable {
    let album: String
    let track: String
    let number: Int
    let explicit: Bool
}

/// The `@Resource` macro attaches a static `resourcename` namespace to the type.
@Resource("//music.example.com/stations/{station_id}")
struct Station {}

// MARK: - Demos

func demoBasic() throws {
    print("\n1. Basic struct (typed generate / parse):")
    let template = try ResourceTemplate("//music.example.com/artists/{id}/{name}")

    let rn = try template.generate(from: Artist(id: "ar-42", name: "Radiohead"))
    print("   Generated:  \(rn)")

    let artist = try template.parse(rn, as: Artist.self)
    print("   Parsed:     id=\(artist.id), name=\(artist.name)")
}

func demoNested() throws {
    print("\n2. Multi-segment struct (label / artist / album):")
    let template = try ResourceTemplate(
        "//music.example.com/labels/{label}/artists/{artist}/albums/{album}")

    let release = Release(label: "xl-recordings", artist: "radiohead", album: "in-rainbows")
    let rn = try template.generate(from: release)
    print("   Generated:  \(rn)")

    let parsed = try template.parse(rn, as: Release.self)
    print("   Parsed:     label=\(parsed.label), artist=\(parsed.artist), album=\(parsed.album)")
}

func demoScalarFields() throws {
    print("\n3. Scalar struct fields (Int / Bool stringified on generate):")
    let template = try ResourceTemplate(
        "//music.example.com/albums/{album}/tracks/{track}/{number}/{explicit}")

    let info = TrackInfo(album: "in-rainbows", track: "15-step", number: 1, explicit: false)
    print("   Generated:  \(try template.generate(from: info))")
}

func demoMacro() throws {
    print("\n4. @Resource macro (attached resourcename namespace):")
    print("   Generated:  \(try Station.resourcename.generate(["station_id": "kexp"]))")

    let parsed = try Station.resourcename.parse("//music.example.com/stations/kexp")
    print("   Parsed:     \(parsed)")
}

// MARK: - Entry point

print("=== Resourcename Demo ===")
do {
    try demoBasic()
    try demoNested()
    try demoScalarFields()
    try demoMacro()
} catch {
    print("error: \(error)")
}
