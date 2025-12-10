import Foundation
import SwiftUI
import GRPCCore
import GRPCNIOTransportHTTP2TransportServices
import GRPCNIOTransportHTTP2

@Observable
class AuthenticationManager {
    var isAuthenticated: Bool = false
    var username: String = ""
    var accessToken: String = ""
    var refreshToken: String = ""
    var userID: String = ""
    var accountID: String = ""
    
    // For iOS Simulator, localhost may not work with TransportServices
    // Use your Mac's IP address instead (e.g., 192.168.1.150)
    // Set this to nil to use localhost, or provide an IP address string
    private let serverHost: String? = "0.0.0.0"
    
    // Client manager following grpc-swift issue #2211 pattern
    // Reuses a single GRPCClient instance across all service clients
    private var clientManager: ClientManager<HTTP2ClientTransport.TransportServices>?
    
    init() {
        print("🔧 AuthenticationManager: Initialized")
        if let host = serverHost {
            print("🔧 Using custom server host: \(host)")
        } else {
            print("🔧 Using localhost for server connection")
        }
    }
    
    /// Get or create the client manager, following the grpc-swift issue #2211 pattern.
    /// This ensures we reuse a single GRPCClient instance across all requests.
    private func getClientManager() throws -> ClientManager<HTTP2ClientTransport.TransportServices> {
        if let existing = clientManager {
            return existing
        }
        
        let host = serverHost ?? "127.0.0.1"
        print("🔧 Creating ClientManager with HTTP2ClientTransport, host: \(host):8001")
        let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(host: host, port: 8001)
        clientManager = manager
        return manager
    }
    
    func login(username: String, password: String, totpToken: String? = nil) async -> (success: Bool, error: String?, requiresTOTP: Bool) {
        print("🔐 Login attempt for user: \(username)")
        
        // Create the login request message
        var loginInput = Auth_UserLoginInput()
        loginInput.username = username
        loginInput.password = password
        if let totpToken = totpToken, !totpToken.isEmpty {
            loginInput.totptoken = totpToken
        }
        
        var requestMessage = Auth_LoginForTokenRequest()
        requestMessage.input = loginInput
        
        print("📤 Creating gRPC client and sending login request...")
        
        do {
            // Check for cancellation before starting
            try Task.checkCancellation()
            
            // Get or create the client manager (follows grpc-swift issue #2211 pattern)
            // This reuses a single GRPCClient instance across all requests
            let manager = try getClientManager()
            
            // Use the auth service client from the unified Client
            // The ClientManager automatically manages connection lifecycle and provides default call options
            // Default timeout (5 seconds) is set at the ClientManager level
            let response = try await manager.client.auth.loginForToken(
                requestMessage,
                options: manager.defaultCallOptions
            )
            
            // Extract the token response
            if response.hasResult {
                let tokenResponse = response.result
                
                print("✅ Login successful, storing tokens")
                
                // Store authentication data
                await MainActor.run {
                    self.isAuthenticated = true
                    self.username = username
                    self.accessToken = tokenResponse.accessToken
                    self.refreshToken = tokenResponse.refreshToken
                    self.userID = tokenResponse.userID
                    self.accountID = tokenResponse.accountID
                }
                
                return (true, nil, false)
            } else {
                print("⚠️ Response received but no token result")
                return (false, "No token received from server", false)
            }
                
        } catch let error as GRPCCore.RPCError {
            print("❌ RPC error code: \(error.code)")
            print("❌ RPC error message: \(error.message)")
            
            // Check if TOTP is required
            let requiresTOTP = error.message.contains("TOTP code required")
            
            // Provide user-friendly error messages
            switch error.code {
            case .deadlineExceeded:
                return (false, "Request timed out. Please check your connection.", false)
            case .unavailable:
                return (false, "Server is unavailable. Please try again later.", false)
            case .unauthenticated:
                if requiresTOTP {
                    return (false, "Please enter your 2FA code.", true)
                }
                return (false, "Invalid username or password.", false)
            default:
                if requiresTOTP {
                    return (false, "Please enter your 2FA code.", true)
                }
                return (false, "Login failed: \(error.message)", false)
            }
        } catch let error as CancellationError {
            print("❌ CancellationError details: \(String(describing: error))")
            print("⏱️ Error occurred at: \(Date())")
        
            
            // Try to get underlying error information
            let nsError = error as NSError
            print("❌ NSError: \(nsError)")
            
            // Check for connection-related error codes
            if nsError.domain == NSPOSIXErrorDomain {
                switch nsError.code {
                case 61: // ECONNREFUSED
                    return (false, "Connection refused. Is the server running on 127.0.0.1:8001?", false)
                case 64: // EHOSTDOWN
                    return (false, "Host is down. Check that the server is running.", false)
                default:
                    break
                }
            }
        } catch {
            print("❌ Error details: \(String(describing: error))")
            
            return (false, "Login failed: \(error.localizedDescription)", false)
        }
        
        return (false, "Unknown login error", false)
    }
    
    func logout() {
        self.isAuthenticated = false
        self.username = ""
        self.accessToken = ""
        self.refreshToken = ""
        self.userID = ""
        self.accountID = ""
    }
}
