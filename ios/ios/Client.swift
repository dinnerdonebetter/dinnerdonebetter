//
//  Client.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2TransportServices

/// A unified gRPC client that provides access to all service clients.
/// This is the Swift analog of the Go client in backend/pkg/client/client.go
@available(macOS 15.0, iOS 18.0, watchOS 11.0, tvOS 18.0, visionOS 2.0, *)
internal struct Client<Transport> where Transport: GRPCCore.ClientTransport {
    /// Auth service client
    internal let auth: Auth_AuthService.Client<Transport>
    
    /// Identity service client
    internal let identity: Identity_IdentityService.Client<Transport>
    
    /// Audit service client
    internal let audit: Audit_AuditService.Client<Transport>
    
    /// Data privacy service client
    internal let dataPrivacy: Dataprivacy_DataPrivacyService.Client<Transport>
    
    /// Internal operations service client
    internal let internalOps: Internalops_InternalOperations.Client<Transport>
    
    /// Meal planning service client
    internal let mealPlanning: Mealplanning_MealPlanningService.Client<Transport>
    
    /// Notifications service client
    internal let notifications: Notifications_UserNotificationsService.Client<Transport>
    
    /// OAuth service client
    internal let oauth: Oauth_OAuthService.Client<Transport>
    
    /// Settings service client
    internal let settings: Settings_SettingsService.Client<Transport>
    
    /// Webhooks service client
    internal let webhooks: Webhooks_WebhooksService.Client<Transport>
    
    /// Internal gRPC client
    private let grpcClient: GRPCCore.GRPCClient<Transport>
    
    /// Initialize a new client with a gRPC client.
    ///
    /// - Parameter grpcClient: The underlying gRPC client to use for all service clients
    internal init(grpcClient: GRPCCore.GRPCClient<Transport>) {
        self.grpcClient = grpcClient
        
        // Initialize all service clients with the same underlying gRPC client
        self.auth = Auth_AuthService.Client(wrapping: grpcClient)
        self.identity = Identity_IdentityService.Client(wrapping: grpcClient)
        self.audit = Audit_AuditService.Client(wrapping: grpcClient)
        self.dataPrivacy = Dataprivacy_DataPrivacyService.Client(wrapping: grpcClient)
        self.internalOps = Internalops_InternalOperations.Client(wrapping: grpcClient)
        self.mealPlanning = Mealplanning_MealPlanningService.Client(wrapping: grpcClient)
        self.notifications = Notifications_UserNotificationsService.Client(wrapping: grpcClient)
        self.oauth = Oauth_OAuthService.Client(wrapping: grpcClient)
        self.settings = Settings_SettingsService.Client(wrapping: grpcClient)
        self.webhooks = Webhooks_WebhooksService.Client(wrapping: grpcClient)
    }
}

// MARK: - Factory Methods

@available(macOS 15.0, iOS 18.0, watchOS 11.0, tvOS 18.0, visionOS 2.0, *)
internal func buildClient<Transport: GRPCCore.ClientTransport>(transport: Transport) -> Client<Transport> {
    let grpcClient = GRPCCore.GRPCClient(transport: transport)
    return Client(grpcClient: grpcClient)
}

/// Build an unauthenticated gRPC client using TransportServices with plaintext security.
///
/// - Parameters:
///   - host: The server host (e.g., "127.0.0.1" or "localhost")
///   - port: The server port (default: 8001)
/// - Returns: A new client instance
/// - Throws: An error if the transport cannot be created
@available(macOS 15.0, iOS 18.0, watchOS 11.0, tvOS 18.0, visionOS 2.0, *)
internal func buildUnauthenticatedClient(
    host: String = "127.0.0.1",
    port: Int = 8001
) throws -> Client<HTTP2ClientTransport.TransportServices> {
    let transport = try HTTP2ClientTransport.TransportServices(
        target: .dns(host: host, port: port),
        transportSecurity: .plaintext
    )
    return buildClient(transport: transport)
}

/// Build an unauthenticated gRPC client using TransportServices with plaintext security.
/// Attempts multiple connection strategies (IPv4, IPv6, DNS) for better compatibility.
///
/// - Parameters:
///   - host: The server host (e.g., "127.0.0.1", "localhost", or nil for auto-detection)
///   - port: The server port (default: 8001)
/// - Returns: A new client instance
/// - Throws: An error if all transport creation attempts fail
@available(macOS 15.0, iOS 18.0, watchOS 11.0, tvOS 18.0, visionOS 2.0, *)
internal func buildUnauthenticatedClientWithFallback(
    host: String? = nil,
    port: Int = 8001
) throws -> Client<HTTP2ClientTransport.TransportServices> {
    // If a specific host is provided, use it directly
    if let host = host {
        return try buildUnauthenticatedClient(host: host, port: port)
    }
    
    // Try IPv4 first (most reliable on iOS Simulator)
    do {
        return try buildUnauthenticatedClient(host: "127.0.0.1", port: port)
    } catch {
        // Try IPv6 next
        do {
            let transport = try HTTP2ClientTransport.TransportServices(
                target: .ipv6(address: "::1", port: port),
                transportSecurity: .plaintext
            )
            return buildClient(transport: transport)
        } catch {
            // Fallback to DNS resolution
            let transport = try HTTP2ClientTransport.TransportServices(
                target: .dns(host: "localhost", port: port),
                transportSecurity: .plaintext
            )
            return buildClient(transport: transport)
        }
    }
}
