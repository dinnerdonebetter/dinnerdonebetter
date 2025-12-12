//
//  APIConfiguration.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import Foundation

struct APIConfiguration {
    // swiftlint:disable:next todo
    // FIXME: Update this with your actual server URL
    static let serverURL = "http://localhost:8000"
    
    // OAuth2 Configuration
    // swiftlint:disable:next todo
    // FIXME: Find some way to configure these
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
