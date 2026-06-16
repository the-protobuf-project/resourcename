import SwiftCompilerPlugin
import SwiftDiagnostics
import SwiftSyntax
import SwiftSyntaxMacros

/// Implementation of the `@Resource("//...")` attached member macro.
///
/// Generates `static let resourcename = try! Resourcename.ResourceNamespace("...")`
/// on the annotated type and validates the template literal at compile time.
public struct ResourceMacro: MemberMacro {
    public static func expansion(
        of node: AttributeSyntax,
        providingMembersOf _: some DeclGroupSyntax,
        conformingTo _: [TypeSyntax],
        in context: some MacroExpansionContext
    ) throws -> [DeclSyntax] {
        guard
            let arguments = node.arguments?.as(LabeledExprListSyntax.self),
            let literalExpr = arguments.first?.expression.as(StringLiteralExprSyntax.self),
            let template = literalExpr.representedLiteralValue
        else {
            context.diagnose(Diagnostic(node: node, message: ResourceMacroDiagnostic.requiresStringLiteral))
            return []
        }

        if let problem = validate(template) {
            context.diagnose(Diagnostic(node: literalExpr, message: problem))
            return []
        }

        let decl: DeclSyntax =
            """
            static let resourcename = try! Resourcename.ResourceNamespace(\(literal: template))
            """
        return [decl]
    }

    /// Mirrors the runtime checks in `ResourceTemplate.init`.
    private static func validate(_ template: String) -> ResourceMacroDiagnostic? {
        if template.isEmpty {
            return .emptyTemplate
        }
        let placeholders = extractPlaceholders(template)
        if placeholders.isEmpty {
            return .noPlaceholders
        }
        let duplicates = duplicates(in: placeholders)
        if !duplicates.isEmpty {
            return .duplicatePlaceholders(duplicates)
        }
        return nil
    }

    private static func extractPlaceholders(_ template: String) -> [String] {
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

    private static func duplicates(in placeholders: [String]) -> [String] {
        var seen = Set<String>()
        var reported = Set<String>()
        var duplicates: [String] = []
        for placeholder in placeholders where !seen.insert(placeholder).inserted {
            if reported.insert(placeholder).inserted {
                duplicates.append(placeholder)
            }
        }
        return duplicates
    }
}

/// Compile-time diagnostics emitted by ``ResourceMacro``.
enum ResourceMacroDiagnostic: DiagnosticMessage {
    case requiresStringLiteral
    case emptyTemplate
    case noPlaceholders
    case duplicatePlaceholders([String])

    var message: String {
        switch self {
        case .requiresStringLiteral:
            return "@Resource requires a static string literal template"
        case .emptyTemplate:
            return "Template cannot be empty"
        case .noPlaceholders:
            return "Template must contain at least one placeholder"
        case let .duplicatePlaceholders(duplicates):
            return "Template contains duplicate placeholders: \(duplicates)"
        }
    }

    var severity: DiagnosticSeverity { .error }

    var diagnosticID: MessageID {
        let id: String
        switch self {
        case .requiresStringLiteral: id = "requiresStringLiteral"
        case .emptyTemplate: id = "emptyTemplate"
        case .noPlaceholders: id = "noPlaceholders"
        case .duplicatePlaceholders: id = "duplicatePlaceholders"
        }
        return MessageID(domain: "ResourcenameMacros", id: id)
    }
}

@main
struct ResourceNamePlugin: CompilerPlugin {
    let providingMacros: [Macro.Type] = [ResourceMacro.self]
}
