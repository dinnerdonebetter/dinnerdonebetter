//
//  AuthenticationManaging.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2TransportServices

/// Protocol defining the authentication interface
/// Allows both AuthenticationManager and MockAuthenticationManager to be used interchangeably
protocol AuthenticationManaging: AnyObject {
    var isAuthenticated: Bool { get set }
    var username: String { get set }
    var accessToken: String { get set }
    var refreshToken: String { get set }
    var oauth2AccessToken: String { get set }
    var oauth2RefreshToken: String { get set }
    var oauth2TokenExpiresAt: Date? { get set }
    var userID: String { get set }
    var accountID: String { get set }
    
    func login(username: String, password: String, totpToken: String?) async -> (success: Bool, error: String?, requiresTOTP: Bool)
    func getClientManager() throws -> ClientManager<HTTP2ClientTransport.TransportServices>
    func getOAuth2AccessToken() async -> String?
    func refreshOAuth2Token() async -> Bool
    func logout()
}
