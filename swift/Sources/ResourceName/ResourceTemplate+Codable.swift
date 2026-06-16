import Foundation

/// Typed parse/generate that maps a struct's fields to template placeholders by name —
/// the Swift counterpart of Go's tagged-struct marshaling and Rust's serde flow.
public extension ResourceTemplate {
    /// Generates a resource name from an `Encodable` value whose property names match
    /// the template placeholders. Scalar fields (`String`, `Bool`, integer, floating
    /// point) are stringified; nested objects are not supported.
    ///
    /// ```swift
    /// struct Artist: Encodable { let id: String; let name: String }
    /// let t = try ResourceTemplate("//music.example.com/artists/{id}/{name}")
    /// try t.generate(from: Artist(id: "ar-42", name: "Radiohead"))
    /// // "//music.example.com/artists/ar-42/Radiohead"
    /// ```
    ///
    /// - Throws: ``ResourceNameError/generateInputMustBeObject`` if `value` does not
    ///   encode to an object, ``ResourceNameError/nonStringValue(field:)`` for a
    ///   non-scalar field, or any error from ``generate(_:)``.
    func generate<Value: Encodable>(from value: Value) throws -> String {
        let data = try JSONEncoder().encode(value)
        let object = try JSONSerialization.jsonObject(with: data)
        guard let dictionary = object as? [String: Any] else {
            throw ResourceNameError.generateInputMustBeObject
        }
        var values: [String: String] = [:]
        values.reserveCapacity(dictionary.count)
        for (key, raw) in dictionary {
            guard let string = Self.scalarString(raw) else {
                throw ResourceNameError.nonStringValue(field: key)
            }
            values[key] = string
        }
        return try generate(values)
    }

    /// Parses a resource name into a `Decodable` value whose property names match the
    /// template placeholders. Each captured segment is a `String`, so the target's
    /// fields should be `String` (mirrors the typed parse in the Rust implementation).
    ///
    /// ```swift
    /// struct Artist: Decodable { let id: String; let name: String }
    /// let t = try ResourceTemplate("//music.example.com/artists/{id}/{name}")
    /// try t.parse("//music.example.com/artists/ar-42/Radiohead", as: Artist.self)
    /// // Artist(id: "ar-42", name: "Radiohead")
    /// ```
    ///
    /// - Throws: ``ResourceNameError/templateMismatch(name:template:)`` on mismatch, or
    ///   a decoding error if the captured segments do not fit `Value`.
    func parse<Value: Decodable>(_ resourceName: String, as _: Value.Type) throws -> Value {
        let values = try parse(resourceName)
        let data = try JSONSerialization.data(withJSONObject: values)
        return try JSONDecoder().decode(Value.self, from: data)
    }

    private static func scalarString(_ value: Any) -> String? {
        switch value {
        case let string as String:
            return string
        case let number as NSNumber:
            // JSON booleans surface as NSNumber; distinguish them from integers.
            if CFGetTypeID(number) == CFBooleanGetTypeID() {
                return number.boolValue ? "true" : "false"
            }
            return number.stringValue
        default:
            return nil
        }
    }
}
