//
//  HomeView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

struct HomeView: View {
    @Environment(AuthenticationManager.self) private var authManager
    
    var body: some View {
        NavigationStack {
            VStack(spacing: 20) {
                Spacer()
                
                Text("Welcome, \(authManager.username)!")
                    .font(.largeTitle)
                    .fontWeight(.bold)
                
                Text("You are now logged in")
                    .font(.title3)
                    .foregroundColor(.secondary)
                
                Spacer()
                
                Button(action: authManager.logout) {
                    Text("Sign Out")
                        .fontWeight(.semibold)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.red)
                        .foregroundColor(.white)
                        .cornerRadius(10)
                }
                .padding(.horizontal, 32)
                .padding(.bottom, 32)
            }
            .navigationTitle("Home")
        }
    }
}

#Preview {
    let authManager = AuthenticationManager()
    authManager.isAuthenticated = true
    authManager.username = "John Doe"
    
    return HomeView()
        .environment(authManager)
}
