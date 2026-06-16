# Resourcename (Swift)

Google [AIP-122](https://aip.dev/122)-style **resource name** templates with
`{placeholder}` segments: compile a template, **parse** a full name into its
components, and **generate** a name from components. A `@Resource` attribute macro
attaches the API to a type — the Swift counterpart of the Python
`@resourcename.resource(...)` decorator.

Each `{placeholder}` matches exactly one path segment (`[^/]+`), so generated
values may not contain `/`.

## Usage

```swift
import Resourcename

// Template API
let t = try ResourceTemplate("//music.example.com/artists/{artist_id}")
try t.parse("//music.example.com/artists/radiohead")  // ["artist_id": "radiohead"]
try t.generate(["artist_id": "bjork"])                // "//music.example.com/artists/bjork"

// Attribute macro — attaches a static `resourcename` namespace
@Resource("//music.example.com/artists/{artist_id}")
struct Artist {}

try Artist.resourcename.parse("//music.example.com/artists/radiohead")
try Artist.resourcename.generate(["artist_id": "bjork"])
```

The `@Resource` template literal is validated at **compile time**: an empty
template, a template with no placeholders, or duplicate placeholder names are
reported as compiler errors.

## Build & test

The package manifest lives at the repository root (`Package.swift`); the Swift
sources are under `swift/`.

```bash
swift build
swift run ResourcenameExamples   # struct-field demos under swift/Examples/
swift test                       # requires full Xcode (XCTest); CLT-only toolchains can't run it
```

## API

- `ResourceTemplate` — `init(_:)`, `parse(_:)`, `generate(_:)`, `template`, `placeholders`.
- `ResourceNamespace` — the object attached as `static let resourcename` by `@Resource`.
- `@Resource(_:)` — attribute macro (see [Swift attributes](https://docs.swift.org/swift-book/documentation/the-swift-programming-language/attributes/)).
- `ResourceNameError` — thrown for invalid templates, mismatches, and bad values.
