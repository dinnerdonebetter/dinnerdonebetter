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
    
    init() {
        print("🔧 AuthenticationManager: Initialized")
        if let host = serverHost {
            print("🔧 Using custom server host: \(host)")
        } else {
            print("🔧 Using localhost for server connection")
        }
    }
    
    func login(username: String, password: String) async -> (success: Bool, error: String?) {
        print("🔐 Login attempt for user: \(username)")
        
        // Create the login request message
        var loginInput = Auth_UserLoginInput()
        loginInput.username = username
        loginInput.password = password
        
        var requestMessage = Auth_LoginForTokenRequest()
        requestMessage.input = loginInput
        
        print("📤 Creating gRPC client and sending login request...")
        
        do {
            // Check for cancellation before starting
            try Task.checkCancellation()
            
            // Use the unified Client factory function which handles transport creation
            // with automatic fallback strategies (IPv4 -> IPv6 -> DNS)
            let client: Client<HTTP2ClientTransport.TransportServices>
            
            
            
            if let customHost = serverHost {
                // Use custom host (IP address) - more reliable on iOS Simulator
                print("🔧 Using custom host: \(customHost):8001")
                client = try buildUnauthenticatedClient(host: customHost, port: 8001)
            } else {
                // Use fallback strategy (tries IPv4, then IPv6, then DNS)
                print("🔧 Using fallback strategy for localhost connection...")
                client = try buildUnauthenticatedClientWithFallback(host: nil, port: 8001)
            }
            
            // Configure timeout for the request (15 seconds to give more time)
            // Also set a longer deadline to ensure connection has time to establish
            var options = GRPCCore.CallOptions.defaults
            options.timeout = .seconds(5)
            
            // Ensure we're not cancelled before making the call
            try Task.checkCancellation()
            
            // Make the call - all variables (transport, grpcClient, authClient) must stay alive
            // during this async operation. They're all in the same scope, so they should.
            let callStartTime = Date()
            print("🔧 Starting loginForToken RPC call at: \(callStartTime)")
            
            let response: Auth_LoginForTokenResponse
            do {
                // Use the auth service client from the unified Client
                
                try await withGRPCClient(
                  transport: .http2NIOPosix(
                    target: .dns(host: "127.0.0.1", port: 8001),
                    transportSecurity: .plaintext
                  )
                ) { client in
                    let greeter = Client(grpcClient: client)
                    let response = try await greeter.auth.loginForToken(
                        requestMessage,
                        options: options
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
                        
                        return (true, "")
                    } else {
                        print("⚠️ Response received but no token result")
                        return (false, "No token received from server")
                    }
                }
                
                let callDuration = abs(callStartTime.timeIntervalSinceNow)
                print("⏱️ RPC call took: \(String(format: "%.2f", callDuration)) seconds")
            } catch let callError {
                let callDuration = abs(callStartTime.timeIntervalSinceNow)
                print("❌ Error during loginForToken call")
                print("⏱️ RPC call duration: \(String(format: "%.2f", callDuration)) seconds")
                print("❌ Error: \(callError)")
                print("❌ Error type: \(type(of: callError))")
                print("❌ Error at: \(Date())")
                
                // Check if it took close to the timeout
                if callDuration >= 14.5 && callDuration <= 15.5 {
                    print("⚠️ Call timed out after ~15 seconds (matches timeout setting)")
                } else if callDuration >= 9.5 && callDuration <= 10.5 {
                    print("⚠️ Call timed out after ~10 seconds")
                }
                
                // Re-throw to be caught by outer catch blocks
                throw callError
            }
            
            print("📥 Received response from server")
            
        } catch let error as GRPCCore.RPCError {
            print("❌ RPC error code: \(error.code)")
            print("❌ RPC error message: \(error.message)")
            
            // Provide user-friendly error messages
            switch error.code {
            case .deadlineExceeded:
                return (false, "Request timed out. Please check your connection.")
            case .unavailable:
                return (false, "Server is unavailable. Please try again later.")
            case .unauthenticated:
                return (false, "Invalid username or password.")
            default:
                return (false, "Login failed: \(error.message)")
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
                    return (false, "Connection refused. Is the server running on 127.0.0.1:8001?")
                case 64: // EHOSTDOWN
                    return (false, "Host is down. Check that the server is running.")
                default:
                    break
                }
            }
        } catch {
            print("❌ Error details: \(String(describing: error))")
            
            return (false, "Login failed: \(error.localizedDescription)")
        }
        
        return (false, "Unknown login error")
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
