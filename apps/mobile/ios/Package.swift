// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "TchatApp",
    platforms: [
        .iOS(.v16),
        .macOS(.v13)
    ],
    products: [
        .library(
            name: "TchatApp",
            targets: ["TchatApp"]
        ),
    ],
    dependencies: [
        // SwiftUI is part of iOS SDK - no external dependency needed
        // Combine is part of iOS SDK - no external dependency needed

        // Additional dependencies for enhanced functionality
        .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.8.0"),
        .package(url: "https://github.com/onevcat/Kingfisher.git", from: "7.9.0"),
        .package(url: "https://github.com/apple/swift-log.git", from: "1.5.0"),

        // Core networking and utilities only
        // Removed swift-navigation to use native SwiftUI navigation instead
    ],
    targets: [
        .target(
            name: "TchatApp",
            dependencies: [
                "Alamofire",
                "Kingfisher",
                .product(name: "Logging", package: "swift-log"),
                // Using native SwiftUI navigation instead of external dependencies
            ],
            path: "Sources"
        ),
        .testTarget(
            name: "TchatAppTests",
            dependencies: ["TchatApp"],
            path: "Tests"
        ),
    ]
)