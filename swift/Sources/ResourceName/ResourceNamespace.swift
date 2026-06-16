/// A namespace exposing `parse` / `generate` plus template metadata.
///
/// This is the Swift counterpart of the Python `ResourceNamespace` attached to a
/// class as `.resourcename`. The ``Resource(_:)`` macro generates a
/// `static let resourcename` of this type on the annotated type.
///
/// ```swift
/// @Resource("//music.example.com/artists/{artist_id}")
/// struct Artist {}
///
/// try Artist.resourcename.parse("//music.example.com/artists/radiohead")
/// try Artist.resourcename.generate(["artist_id": "bjork"])
/// ```
public struct ResourceNamespace: Sendable {
    private let resourceTemplate: ResourceTemplate

    /// Wraps an already-compiled template.
    public init(_ template: ResourceTemplate) {
        self.resourceTemplate = template
    }

    /// Compiles `template` and wraps it.
    ///
    /// - Throws: ``ResourceNameError`` when the template is invalid.
    public init(_ template: String) throws {
        self.resourceTemplate = try ResourceTemplate(template)
    }

    /// Parses a resource name into its component values.
    public func parse(_ resourceName: String) throws -> [String: String] {
        try resourceTemplate.parse(resourceName)
    }

    /// Generates a resource name from component values.
    public func generate(_ values: [String: String]) throws -> String {
        try resourceTemplate.generate(values)
    }

    /// The original template string.
    public var template: String { resourceTemplate.template }

    /// Placeholder names in declaration order.
    public var placeholders: [String] { resourceTemplate.placeholders }
}
