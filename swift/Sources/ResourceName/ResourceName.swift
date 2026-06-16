// Resourcename — Google AIP-122 style resource-name templates for Swift.
//
// Public surface:
//   - `ResourceTemplate`  — compile a template, then `parse` / `generate`.
//   - `ResourceNamespace` — the object attached by the `@Resource` macro.
//   - `@Resource(_:)`     — attribute macro that attaches a static `resourcename`
//                           namespace to a type (the Swift counterpart of the
//                           Python `@resourcename.resource(...)` decorator).
//   - `ResourceNameError` — errors thrown by the above.

/// Attaches a `static let resourcename` namespace to the annotated type, built
/// from a `{placeholder}` template string.
///
/// The template literal is validated at compile time: an empty template, a
/// template with no placeholders, or duplicate placeholder names are reported as
/// compiler errors (mirroring the runtime checks in ``ResourceTemplate``).
///
/// ```swift
/// import Resourcename
///
/// @Resource("//music.example.com/artists/{artist_id}")
/// struct Artist {}
///
/// try Artist.resourcename.parse("//music.example.com/artists/radiohead")
/// // ["artist_id": "radiohead"]
/// try Artist.resourcename.generate(["artist_id": "bjork"])
/// // "//music.example.com/artists/bjork"
/// ```
@attached(member, names: named(resourcename))
public macro Resource(_ template: String) =
    #externalMacro(module: "ResourcenameMacros", type: "ResourceMacro")
