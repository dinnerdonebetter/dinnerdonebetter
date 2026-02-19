//
//  APIConfiguration.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import Foundation

/// Represents the different deployment environments for the app.
enum AppEnvironment: String, CaseIterable {
  case local
  case development
  case production

  var displayName: String {
    switch self {
    case .local: return "Local"
    case .development: return "Development"
    case .production: return "Production"
    }
  }
}

struct APIConfiguration {
  /// Current environment - defaults based on build configuration.
  /// Can be changed at runtime for testing purposes.
  #if DEBUG
    static var currentEnvironment: AppEnvironment = .local
  #else
    static var currentEnvironment: AppEnvironment = .production
  #endif

  /// The base URL for API requests (HTTP/REST endpoints)
  static var serverURL: String {
    switch currentEnvironment {
    case .local:
      return "http://localhost:8000"
    case .development:
      return "https://api.dinnerdonebetter.dev"
    case .production:
      return "https://http-api.dinnerdonebetter.com"
    }
  }

  /// The web URL used for Universal Links
  /// This should match the domain in your entitlements file
  static var webURL: String {
    switch currentEnvironment {
    case .local:
      return "http://localhost:3000"
    case .development:
      return "https://www.dinnerdonebetter.dev"
    case .production:
      return "https://www.dinnerdonebetter.com"
    }
  }

  /// The gRPC host for the current environment
  static var grpcHost: String {
    switch currentEnvironment {
    case .local:
      return "0.0.0.0"
    case .development:
      return "grpc.dinnerdonebetter.dev"
    case .production:
      return "api.dinnerdonebetter.com"
    }
  }

  /// The gRPC port (typically the same across environments)
  static var grpcPort: Int {
    return 8001
  }

  // OAuth2 Configuration
  // swiftlint:disable:next todo
  // FIXME: Find some way to configure these per environment
  static let oauth2ClientID = "AAAAAAAAAAAAAAAA"
  static let oauth2ClientSecret = "AAAAAAAAAAAAAAAA"

  // OAuth2 endpoints
  static var oauth2AuthorizeURL: String {
    return "\(serverURL)/oauth2/authorize"
  }

  static var oauth2TokenURL: String {
    return "\(serverURL)/oauth2/token"
  }
}
