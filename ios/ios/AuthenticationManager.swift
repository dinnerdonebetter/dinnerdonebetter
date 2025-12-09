//
//  AuthenticationManager.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import Foundation
import SwiftUI

@Observable
class AuthenticationManager {
    var isAuthenticated: Bool = false
    var username: String = ""
    
    func login(username: String, password: String) -> Bool {
        // TODO: Replace with actual authentication logic
        // For now, accept any non-empty username and password
        if !username.isEmpty && !password.isEmpty {
            self.isAuthenticated = true
            self.username = username
            return true
        }
        return false
    }
    
    func logout() {
        self.isAuthenticated = false
        self.username = ""
    }
}
