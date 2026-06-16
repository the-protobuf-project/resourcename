/// Errors thrown while validating templates, parsing resource names, or generating values.
///
/// Mirrors the `ValueError` cases raised by the Python reference implementation.
public enum ResourceNameError: Error, Equatable, Sendable, CustomStringConvertible {
    /// The template string was empty.
    case emptyTemplate
    /// The template contained no `{placeholder}` segments.
    case noPlaceholders
    /// The same placeholder name appeared more than once.
    case duplicatePlaceholders([String])
    /// A resource name did not match the compiled template.
    case templateMismatch(name: String, template: String)
    /// Required placeholder values were missing while generating.
    case missingValues(missing: [String], required: [String], provided: [String])
    /// Values were provided for keys that are not placeholders in the template.
    case unexpectedValues(unexpected: [String], expected: [String])
    /// One or more values contained `'/'`, which is not allowed in a single segment.
    case slashInValue([String: String])

    public var description: String {
        switch self {
        case .emptyTemplate:
            return "Template cannot be empty"
        case .noPlaceholders:
            return "Template must contain at least one placeholder"
        case let .duplicatePlaceholders(duplicates):
            return "Template contains duplicate placeholders: \(duplicates)"
        case let .templateMismatch(name, template):
            return "Resource name '\(name)' does not match template '\(template)'"
        case let .missingValues(missing, required, provided):
            return "Missing values for placeholders: \(missing.sorted()). "
                + "Required: \(required), Provided: \(provided)"
        case let .unexpectedValues(unexpected, expected):
            return "Unexpected values provided: \(unexpected.sorted()). Expected only: \(expected)"
        case let .slashInValue(invalid):
            return "Values contain invalid character '/': \(invalid)"
        }
    }
}
