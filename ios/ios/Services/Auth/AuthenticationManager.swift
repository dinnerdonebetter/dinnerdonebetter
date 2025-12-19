import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2
import GRPCNIOTransportHTTP2TransportServices
import SwiftUI

/// URLSessionDelegate that prevents automatic redirect following
private class NoRedirectDelegate: NSObject, URLSessionTaskDelegate {
  func urlSession(
    _ session: URLSession, task: URLSessionTask,
    willPerformHTTPRedirection response: HTTPURLResponse, newRequest request: URLRequest,
    completionHandler: @escaping (URLRequest?) -> Void
  ) {
    // Return nil to prevent following redirects
    completionHandler(nil)
  }
}

@Observable
// swiftlint:disable:next type_body_length
class AuthenticationManager: AuthenticationManaging {
  var isAuthenticated: Bool = false
  var username: String = ""
  var accessToken: String = ""  // JWT token from login
  var refreshToken: String = ""  // JWT refresh token
  var oauth2AccessToken: String = ""  // OAuth2 access token for API calls
  var oauth2RefreshToken: String = ""  // OAuth2 refresh token
  var oauth2TokenExpiresAt: Date?  // OAuth2 token expiration
  var userID: String = ""
  var accountID: String = ""

  // For iOS Simulator, localhost may not work with TransportServices
  // Use your Mac's IP address instead (e.g., 192.168.1.150)
  // Set this to nil to use localhost, or provide an IP address string
  private let serverHost: String? = "0.0.0.0"

  // Client manager following grpc-swift issue #2211 pattern
  // Reuses a single GRPCClient instance across all service clients
  private var clientManager: ClientManager<HTTP2ClientTransport.TransportServices>?

  // Mock support for UI tests
  private var mockManager: MockAuthenticationManager?
  private var isUsingMock: Bool {
    ProcessInfo.processInfo.arguments.contains("--use-mock-auth")
  }

  init() {
    print("🔧 AuthenticationManager: Initialized")

    // Check if we should use mock behavior (for UI tests)
    if isUsingMock {
      let mock = MockAuthenticationManager()
      let launchArguments = ProcessInfo.processInfo.arguments
      if launchArguments.contains("--mock-requires-totp") {
        mock.configure(behavior: .requiresTOTP)
      } else if launchArguments.contains("--mock-success") {
        mock.configure(behavior: .success)
      } else {
        // Default to requires TOTP for UI tests
        mock.configure(behavior: .requiresTOTP)
      }
      self.mockManager = mock
      print("🎭 AuthenticationManager: Using mock mode")
    } else {
      if let host = serverHost {
        print("🔧 Using custom server host: \(host)")
      } else {
        print("🔧 Using localhost for server connection")
      }
    }
  }

  /// Get or create the client manager, following the grpc-swift issue #2211 pattern.
  /// This ensures we reuse a single GRPCClient instance across all requests.
  func getClientManager() throws -> ClientManager<HTTP2ClientTransport.TransportServices> {
    // If using mock, delegate to mock manager
    if let mock = mockManager {
      return try mock.getClientManager()
    }

    if let existing = clientManager {
      return existing
    }

    let host = serverHost ?? "127.0.0.1"
    print("🔧 Creating ClientManager with HTTP2ClientTransport, host: \(host):8001")
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(host: host, port: 8001)
    clientManager = manager
    return manager
  }

  // swiftlint:disable:next cyclomatic_complexity
  func login(username: String, password: String, totpToken: String? = nil) async -> LoginResult {
    // If using mock, delegate to mock manager
    if let mock = mockManager {
      let result = await mock.login(username: username, password: password, totpToken: totpToken)
      // Sync state from mock to self
      await MainActor.run {
        self.isAuthenticated = mock.isAuthenticated
        self.username = mock.username
        self.accessToken = mock.accessToken
        self.refreshToken = mock.refreshToken
        self.userID = mock.userID
        self.accountID = mock.accountID
      }
      return result
    }

    print("🔐 Login attempt for user: \(username)")

    // Create the login request message
    var loginInput = Auth_UserLoginInput()
    loginInput.username = username
    loginInput.password = password
    if let totpToken = totpToken, !totpToken.isEmpty {
      loginInput.totpToken = totpToken
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

        // Store JWT authentication data
        await MainActor.run {
          self.isAuthenticated = true
          self.username = username
          self.accessToken = tokenResponse.accessToken
          self.refreshToken = tokenResponse.refreshToken
          self.userID = tokenResponse.userID
          self.accountID = tokenResponse.accountID
        }

        // Exchange JWT token for OAuth2 token
        print("🔄 Exchanging JWT token for OAuth2 token...")
        let oauth2Result = await exchangeJWTForOAuth2Token(jwtToken: tokenResponse.accessToken)
        if !oauth2Result.success {
          print("⚠️ Failed to get OAuth2 token: \(oauth2Result.error ?? "Unknown error")")
          // Continue anyway - we can retry later
        } else {
          print("✅ OAuth2 token obtained successfully")
        }

        return LoginResult(success: true, error: nil, requiresTOTP: false)
      } else {
        print("⚠️ Response received but no token result")
        return LoginResult(
          success: false, error: "No token received from server", requiresTOTP: false)
      }
    } catch let error as GRPCCore.RPCError {
      print("❌ RPC error code: \(error.code)")
      print("❌ RPC error message: \(error.message)")

      // Check if TOTP is required
      let requiresTOTP = error.message.contains("TOTP code required")

      // Provide user-friendly error messages
      switch error.code {
      case .deadlineExceeded:
        return LoginResult(
          success: false, error: "Request timed out. Please check your connection.",
          requiresTOTP: false)
      case .unavailable:
        return LoginResult(
          success: false, error: "Server is unavailable. Please try again later.",
          requiresTOTP: false)
      case .unauthenticated:
        if requiresTOTP {
          return LoginResult(
            success: false, error: "Please enter your 2FA code.", requiresTOTP: true)
        }
        return LoginResult(
          success: false, error: "Invalid username or password.", requiresTOTP: false)
      default:
        if requiresTOTP {
          return LoginResult(
            success: false, error: "Please enter your 2FA code.", requiresTOTP: true)
        }
        return LoginResult(
          success: false, error: "Login failed: \(error.message)", requiresTOTP: false)
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
        case 61:  // ECONNREFUSED
          return LoginResult(
            success: false, error: "Connection refused. Is the server running on 127.0.0.1:8001?",
            requiresTOTP: false)
        case 64:  // EHOSTDOWN
          return LoginResult(
            success: false, error: "Host is down. Check that the server is running.",
            requiresTOTP: false)
        default:
          break
        }
      }
    } catch {
      print("❌ Error details: \(String(describing: error))")

      return LoginResult(
        success: false, error: "Login failed: \(error.localizedDescription)", requiresTOTP: false)
    }

    return LoginResult(success: false, error: "Unknown login error", requiresTOTP: false)
  }

  /// Exchange JWT token for OAuth2 access token
  private func exchangeJWTForOAuth2Token(jwtToken: String) async -> (
    success: Bool, error: String?
  ) {
    // Generate state for OAuth2 flow
    let state = UUID().uuidString

    // Build authorization URL
    guard var components = URLComponents(string: APIConfiguration.oauth2AuthorizeURL) else {
      return (false, "Failed to build authorization URL")
    }
    components.queryItems = [
      URLQueryItem(name: "response_type", value: "code"),
      URLQueryItem(name: "client_id", value: APIConfiguration.oauth2ClientID),
      URLQueryItem(name: "redirect_uri", value: APIConfiguration.serverURL),
      URLQueryItem(name: "state", value: state),
      URLQueryItem(name: "code_challenge_method", value: "plain"),
    ]

    guard let authURL = components.url else {
      return (false, "Failed to build authorization URL")
    }

    // Step 1: Get authorization code
    var request = URLRequest(url: authURL)
    request.httpMethod = "GET"
    request.setValue("Bearer \(jwtToken)", forHTTPHeaderField: "Authorization")

    // Configure session to not follow redirects
    let delegate = NoRedirectDelegate()
    let session = URLSession(configuration: .default, delegate: delegate, delegateQueue: nil)

    do {
      let (_, response) = try await session.data(for: request)

      guard let httpResponse = response as? HTTPURLResponse else {
        return (false, "Invalid response type")
      }

      // Check if this is a redirect response (3xx status code)
      guard (300...399).contains(httpResponse.statusCode) else {
        return (false, "Expected redirect response, got status code: \(httpResponse.statusCode)")
      }

      // Get the authorization code from the Location header
      guard let locationHeader = httpResponse.value(forHTTPHeaderField: "Location"),
        let locationURL = URL(string: locationHeader),
        let codeComponents = URLComponents(url: locationURL, resolvingAgainstBaseURL: false),
        let code = codeComponents.queryItems?.first(where: { $0.name == "code" })?.value
      else {
        return (false, "No authorization code in response")
      }

      print("📝 Received OAuth2 authorization code")

      // Step 2: Exchange code for access token
      return await exchangeCodeForToken(code: code)
    } catch {
      print("❌ Error getting authorization code: \(error)")
      return (false, "Failed to get authorization code: \(error.localizedDescription)")
    }
  }

  /// Exchange OAuth2 authorization code for access token
  private func exchangeCodeForToken(code: String) async -> (success: Bool, error: String?) {
    guard let tokenURL = URL(string: APIConfiguration.oauth2TokenURL) else {
      return (false, "Invalid token URL")
    }

    var request = URLRequest(url: tokenURL)
    request.httpMethod = "POST"
    request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")

    // Build form data
    let formData = [
      "grant_type": "authorization_code",
      "code": code,
      "redirect_uri": APIConfiguration.serverURL,
      "client_id": APIConfiguration.oauth2ClientID,
      "client_secret": APIConfiguration.oauth2ClientSecret,
    ]

    // Encode form data
    let formString = formData.map {
      "\($0.key)=\($0.value.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? $0.value)"
    }
    .joined(separator: "&")
    request.httpBody = formString.data(using: .utf8)

    do {
      let (data, response) = try await URLSession.shared.data(for: request)

      guard let httpResponse = response as? HTTPURLResponse else {
        return (false, "Invalid response type")
      }

      guard (200...299).contains(httpResponse.statusCode) else {
        let errorMessage = String(data: data, encoding: .utf8) ?? "Unknown error"
        return (false, "Token exchange failed: \(errorMessage)")
      }

      // Parse JSON response
      guard let json = try JSONSerialization.jsonObject(with: data) as? [String: Any],
        let accessToken = json["access_token"] as? String
      else {
        return (false, "Invalid token response format")
      }

      let refreshToken = json["refresh_token"] as? String ?? ""
      let expiresIn = json["expires_in"] as? Int ?? 86400  // Default to 24 hours

      // Calculate expiration time
      let expiresAt = Date().addingTimeInterval(TimeInterval(expiresIn))

      await MainActor.run {
        self.oauth2AccessToken = accessToken
        self.oauth2RefreshToken = refreshToken
        self.oauth2TokenExpiresAt = expiresAt
      }

      return (true, nil)
    } catch {
      print("❌ Error exchanging code for token: \(error)")
      return (false, "Failed to exchange code: \(error.localizedDescription)")
    }
  }

  /// Refresh OAuth2 access token using refresh token
  func refreshOAuth2Token() async -> Bool {
    // If using mock, delegate to mock manager
    if let mock = mockManager {
      let result = await mock.refreshOAuth2Token()
      // Sync state
      await MainActor.run {
        self.oauth2AccessToken = mock.oauth2AccessToken
        self.oauth2RefreshToken = mock.oauth2RefreshToken
        self.oauth2TokenExpiresAt = mock.oauth2TokenExpiresAt
      }
      return result
    }

    guard !oauth2RefreshToken.isEmpty else {
      print("⚠️ No OAuth2 refresh token available")
      return false
    }

    guard let tokenURL = URL(string: APIConfiguration.oauth2TokenURL) else {
      print("❌ Invalid token URL")
      return false
    }

    var request = URLRequest(url: tokenURL)
    request.httpMethod = "POST"
    request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")

    // Build form data for refresh
    let formData = [
      "grant_type": "refresh_token",
      "refresh_token": oauth2RefreshToken,
      "client_id": APIConfiguration.oauth2ClientID,
      "client_secret": APIConfiguration.oauth2ClientSecret,
    ]

    let formString = formData.map {
      "\($0.key)=\($0.value.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? $0.value)"
    }
    .joined(separator: "&")
    request.httpBody = formString.data(using: .utf8)

    do {
      let (data, response) = try await URLSession.shared.data(for: request)

      guard let httpResponse = response as? HTTPURLResponse,
        (200...299).contains(httpResponse.statusCode)
      else {
        print("❌ Token refresh failed")
        return false
      }

      guard let json = try JSONSerialization.jsonObject(with: data) as? [String: Any],
        let accessToken = json["access_token"] as? String
      else {
        print("❌ Invalid token response format")
        return false
      }

      let refreshToken = json["refresh_token"] as? String ?? oauth2RefreshToken
      let expiresIn = json["expires_in"] as? Int ?? 86400
      let expiresAt = Date().addingTimeInterval(TimeInterval(expiresIn))

      await MainActor.run {
        self.oauth2AccessToken = accessToken
        self.oauth2RefreshToken = refreshToken
        self.oauth2TokenExpiresAt = expiresAt
      }

      print("✅ OAuth2 token refreshed successfully")
      return true
    } catch {
      print("❌ Error refreshing token: \(error)")
      return false
    }
  }

  /// Get OAuth2 access token, refreshing if needed
  func getOAuth2AccessToken() async -> String? {
    // If using mock, delegate to mock manager
    if let mock = mockManager {
      return await mock.getOAuth2AccessToken()
    }

    // Check if token is expired or will expire soon (within 5 minutes)
    if let expiresAt = oauth2TokenExpiresAt,
      expiresAt.timeIntervalSinceNow < 300
    {
      print("🔄 OAuth2 token expired or expiring soon, refreshing...")
      let refreshed = await refreshOAuth2Token()
      if !refreshed {
        // If refresh fails, try to get a new token using JWT
        if !accessToken.isEmpty {
          print("🔄 Attempting to get new OAuth2 token using JWT...")
          let result = await exchangeJWTForOAuth2Token(jwtToken: accessToken)
          if !result.success {
            return nil
          }
        } else {
          return nil
        }
      }
    }

    return oauth2AccessToken.isEmpty ? nil : oauth2AccessToken
  }

  func logout() {
    self.isAuthenticated = false
    self.username = ""
    self.accessToken = ""
    self.refreshToken = ""
    self.oauth2AccessToken = ""
    self.oauth2RefreshToken = ""
    self.oauth2TokenExpiresAt = nil
    self.userID = ""
    self.accountID = ""
  }
}
