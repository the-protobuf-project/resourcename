import Foundation

/// An immutable resource-name template with `{placeholder}` segments.
///
/// Each `{placeholder}` matches exactly one path segment (`[^/]+`), so generated
/// values may not contain `'/'`. This is a direct port of the Python
/// `ResourceTemplate`: parse a full resource name into its components, or generate
/// a name from component values.
///
/// ```swift
/// let t = try ResourceTemplate("//music.example.com/artists/{artist_id}")
/// try t.parse("//music.example.com/artists/radiohead")  // ["artist_id": "radiohead"]
/// try t.generate(["artist_id": "bjork"])                // "//music.example.com/artists/bjork"
/// ```
///
/// `NSRegularExpression` is immutable and thread-safe, so this value type is
/// `@unchecked Sendable`.
public struct ResourceTemplate: @unchecked Sendable {
    /// The original template string.
    public let template: String
    /// Placeholder names in declaration order.
    public let placeholders: [String]

    private let regex: NSRegularExpression

    /// Compiles a template string.
    ///
    /// - Throws: ``ResourceNameError/emptyTemplate``,
    ///   ``ResourceNameError/noPlaceholders``, or
    ///   ``ResourceNameError/duplicatePlaceholders(_:)`` when invalid.
    public init(_ template: String) throws {
        if template.isEmpty {
            throw ResourceNameError.emptyTemplate
        }
        let placeholders = Self.extractPlaceholders(template)
        if placeholders.isEmpty {
            throw ResourceNameError.noPlaceholders
        }
        let duplicates = Self.duplicates(in: placeholders)
        if !duplicates.isEmpty {
            throw ResourceNameError.duplicatePlaceholders(duplicates)
        }
        self.template = template
        self.placeholders = placeholders
        self.regex = try Self.compileRegex(template, placeholders: placeholders)
    }

    /// Parses a resource name into a `placeholder -> value` dictionary.
    ///
    /// - Throws: ``ResourceNameError/templateMismatch(name:template:)`` when the
    ///   name does not match this template.
    public func parse(_ resourceName: String) throws -> [String: String] {
        let fullRange = NSRange(resourceName.startIndex..., in: resourceName)
        guard
            let match = regex.firstMatch(in: resourceName, range: fullRange),
            match.range == fullRange
        else {
            throw ResourceNameError.templateMismatch(name: resourceName, template: template)
        }
        var values: [String: String] = [:]
        for (index, placeholder) in placeholders.enumerated() {
            let captureRange = match.range(at: index + 1)
            if let range = Range(captureRange, in: resourceName) {
                values[placeholder] = String(resourceName[range])
            }
        }
        return values
    }

    /// Generates a resource name by substituting `values` into the template.
    ///
    /// - Throws: ``ResourceNameError/missingValues(missing:required:provided:)``,
    ///   ``ResourceNameError/unexpectedValues(unexpected:expected:)``, or
    ///   ``ResourceNameError/slashInValue(_:)``.
    public func generate(_ values: [String: String]) throws -> String {
        let expected = Set(placeholders)
        let provided = Set(values.keys)

        let missing = expected.subtracting(provided)
        if !missing.isEmpty {
            throw ResourceNameError.missingValues(
                missing: Array(missing),
                required: placeholders,
                provided: Array(values.keys)
            )
        }

        let extra = provided.subtracting(expected)
        if !extra.isEmpty {
            throw ResourceNameError.unexpectedValues(
                unexpected: Array(extra),
                expected: placeholders
            )
        }

        let invalid = values.filter { $0.value.contains("/") }
        if !invalid.isEmpty {
            throw ResourceNameError.slashInValue(invalid)
        }

        var result = template
        for placeholder in placeholders {
            result = result.replacingOccurrences(of: "{\(placeholder)}", with: values[placeholder]!)
        }
        return result
    }

    // MARK: - Internals

    /// Returns the placeholder names found inside `{...}`, in order.
    static func extractPlaceholders(_ template: String) -> [String] {
        var result: [String] = []
        var rest = Substring(template)
        while let open = rest.firstIndex(of: "{") {
            let afterOpen = rest.index(after: open)
            guard let close = rest[afterOpen...].firstIndex(of: "}") else { break }
            let name = String(rest[afterOpen..<close])
            if !name.isEmpty {
                result.append(name)
            }
            rest = rest[rest.index(after: close)...]
        }
        return result
    }

    /// Returns each placeholder that appears more than once (first occurrence order).
    static func duplicates(in placeholders: [String]) -> [String] {
        var seen = Set<String>()
        var reported = Set<String>()
        var duplicates: [String] = []
        for placeholder in placeholders {
            if seen.contains(placeholder) {
                if reported.insert(placeholder).inserted {
                    duplicates.append(placeholder)
                }
            } else {
                seen.insert(placeholder)
            }
        }
        return duplicates
    }

    /// Builds an anchored regex where each placeholder becomes a `([^/]+)` capture group.
    static func compileRegex(_ template: String, placeholders: [String]) throws -> NSRegularExpression {
        let marker = "\u{0}"
        var pattern = template
        for placeholder in placeholders {
            pattern = pattern.replacingOccurrences(of: "{\(placeholder)}", with: marker)
        }
        pattern = NSRegularExpression.escapedPattern(for: pattern)
        pattern = pattern.replacingOccurrences(of: marker, with: "([^/]+)")
        return try NSRegularExpression(pattern: "^\(pattern)$")
    }
}
