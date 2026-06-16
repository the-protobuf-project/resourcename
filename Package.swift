// swift-tools-version: 6.3
// The swift-tools-version declares the minimum version of Swift required to build this package.

import CompilerPluginSupport
import PackageDescription

let package = Package(
    name: "Resourcename",
    platforms: [
        .macOS(.v13),
    ],
    products: [
        .library(
            name: "Resourcename",
            targets: ["Resourcename"]
        ),
    ],
    dependencies: [
        .package(url: "https://github.com/swiftlang/swift-syntax.git", "600.0.0" ..< "700.0.0"),
    ],
    targets: [
        // Macro implementation (compiler plugin) that powers `@Resource`.
        .macro(
            name: "ResourcenameMacros",
            dependencies: [
                .product(name: "SwiftSyntaxMacros", package: "swift-syntax"),
                .product(name: "SwiftCompilerPlugin", package: "swift-syntax"),
            ],
            path: "swift/Sources/ResourcenameMacros"
        ),
        // Public library: ResourceTemplate, ResourceNamespace, ResourceNameError, @Resource.
        .target(
            name: "Resourcename",
            dependencies: ["ResourcenameMacros"],
            path: "swift/Sources/Resourcename"
        ),
        // Runnable demos with struct fields: `swift run ResourcenameExamples`.
        .executableTarget(
            name: "ResourcenameExamples",
            dependencies: ["Resourcename"],
            path: "swift/Examples"
        ),
        .testTarget(
            name: "ResourcenameTests",
            dependencies: ["Resourcename"],
            path: "swift/Tests/ResourcenameTests"
        ),
    ],
    swiftLanguageModes: [.v6]
)
