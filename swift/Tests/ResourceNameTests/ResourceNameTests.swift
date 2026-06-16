import XCTest

@testable import Resourcename

final class ResourceTemplateTests: XCTestCase {
    func testParseExtractsPlaceholders() throws {
        let template = try ResourceTemplate("//music.example.com/artists/{artist_id}")
        XCTAssertEqual(
            try template.parse("//music.example.com/artists/radiohead"),
            ["artist_id": "radiohead"]
        )
    }

    func testGenerateSubstitutesValues() throws {
        let template = try ResourceTemplate("//music.example.com/artists/{artist_id}")
        XCTAssertEqual(
            try template.generate(["artist_id": "bjork"]),
            "//music.example.com/artists/bjork"
        )
    }

    func testNestedRoundTrip() throws {
        let template = try ResourceTemplate(
            "//music.example.com/artists/{artist_id}/albums/{album_id}")
        let name = try template.generate(["artist_id": "radiohead", "album_id": "in-rainbows"])
        XCTAssertEqual(name, "//music.example.com/artists/radiohead/albums/in-rainbows")
        XCTAssertEqual(
            try template.parse(name),
            ["artist_id": "radiohead", "album_id": "in-rainbows"]
        )
    }

    func testMetadataIsExposed() throws {
        let template = try ResourceTemplate(
            "//music.example.com/artists/{artist_id}/albums/{album_id}")
        XCTAssertEqual(template.template, "//music.example.com/artists/{artist_id}/albums/{album_id}")
        XCTAssertEqual(template.placeholders, ["artist_id", "album_id"])
    }

    func testRejectsEmptyTemplate() {
        XCTAssertThrowsError(try ResourceTemplate("")) {
            XCTAssertEqual($0 as? ResourceNameError, .emptyTemplate)
        }
    }

    func testRejectsTemplateWithoutPlaceholders() {
        XCTAssertThrowsError(try ResourceTemplate("//music.example.com/artists")) {
            XCTAssertEqual($0 as? ResourceNameError, .noPlaceholders)
        }
    }

    func testRejectsDuplicatePlaceholders() {
        XCTAssertThrowsError(try ResourceTemplate("//x/{id}/y/{id}")) {
            XCTAssertEqual($0 as? ResourceNameError, .duplicatePlaceholders(["id"]))
        }
    }

    func testParseRejectsMismatch() throws {
        let template = try ResourceTemplate("//music.example.com/artists/{artist_id}")
        XCTAssertThrowsError(try template.parse("//other.example.com/things/x"))
    }

    func testParseRejectsSlashInSegment() throws {
        // A single placeholder segment is `[^/]+`, so extra slashes must not match.
        let template = try ResourceTemplate("//music.example.com/artists/{artist_id}")
        XCTAssertThrowsError(try template.parse("//music.example.com/artists/label/radiohead"))
    }

    func testGenerateRejectsMissingValues() throws {
        let template = try ResourceTemplate(
            "//music.example.com/artists/{artist_id}/albums/{album_id}")
        XCTAssertThrowsError(try template.generate(["artist_id": "radiohead"]))
    }

    func testGenerateRejectsExtraValues() throws {
        let template = try ResourceTemplate("//music.example.com/artists/{artist_id}")
        XCTAssertThrowsError(try template.generate(["artist_id": "radiohead", "other": "x"]))
    }

    func testGenerateRejectsSlashInValue() throws {
        let template = try ResourceTemplate("//music.example.com/artists/{artist_id}")
        XCTAssertThrowsError(try template.generate(["artist_id": "label/radiohead"]))
    }
}

// MARK: - @Resource macro

@Resource("//music.example.com/artists/{artist_id}")
private struct Artist {}

@Resource("//music.example.com/artists/{artist_id}/albums/{album_id}")
private struct Album {}

final class ResourceMacroTests: XCTestCase {
    func testMacroAttachesParse() throws {
        XCTAssertEqual(
            try Artist.resourcename.parse("//music.example.com/artists/radiohead"),
            ["artist_id": "radiohead"]
        )
    }

    func testMacroAttachesGenerate() throws {
        XCTAssertEqual(
            try Artist.resourcename.generate(["artist_id": "bjork"]),
            "//music.example.com/artists/bjork"
        )
    }

    func testMacroSupportsNestedTemplates() throws {
        let name = try Album.resourcename.generate([
            "artist_id": "radiohead", "album_id": "in-rainbows",
        ])
        XCTAssertEqual(name, "//music.example.com/artists/radiohead/albums/in-rainbows")
        XCTAssertEqual(Album.resourcename.placeholders, ["artist_id", "album_id"])
    }
}
