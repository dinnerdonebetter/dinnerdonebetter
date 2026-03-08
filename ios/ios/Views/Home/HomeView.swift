//
//  HomeView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import Combine
import SwiftUI

struct HomeView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(UserSettingsService.self) private var userSettingsService
  @State private var viewModel: HomeViewModel?

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          DSContentState(
            isLoading: viewModel.isLoading,
            loadingMessage: "Loading...",
            error: viewModel.errorMessage,
            errorTitle: viewModel.errorTitle,
            errorIcon: viewModel.errorIcon,
            errorIconColor: viewModel.errorIconColor,
            onRetry: { await viewModel.loadData() },
            content: {
              VStack(spacing: 0) {
                // Sticky header: Greeting
                greetingSection(viewModel: viewModel)
                  .dsScreenPadding()
                  .padding(.bottom, DSTheme.Spacing.md)

                MealPlanningHomeContent(viewModel: viewModel)
              }
            }
          )
        } else {
          DSInitializingView()
        }
      }
      .navigationTitle("")
      .navigationBarTitleDisplayMode(.inline)
      .toolbar {
        ToolbarItem(placement: .navigationBarTrailing) {
          NavigationLink(destination: AccountSettingsView()) {
            DSAvatar(
              name: viewModel?.currentUserDisplayName ?? authManager.username,
              size: .sm,
              imageURL: viewModel.flatMap { homeViewModel in
                guard let user = homeViewModel.currentUser,
                  user.hasAvatar,
                  !user.avatar.storagePath.isEmpty
                else { return nil }
                return APIConfiguration.mediaURL(
                  forStoragePath: user.avatar.storagePath, bucket: "avatars")
              }
            )
          }
        }
      }
      .refreshable {
        if let viewModel = viewModel {
          await viewModel.loadData()
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = HomeViewModel(authManager: authManager)
        }
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
      .onReceive(NotificationCenter.default.publisher(for: .mealPlanCreated)) { _ in
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
      .onReceive(NotificationCenter.default.publisher(for: .mealPlanArchived)) { _ in
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
      .onReceive(NotificationCenter.default.publisher(for: .mealPlanEventsUpdated)) { _ in
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
    }
  }

  // MARK: - Greeting
  private func greetingSection(viewModel: HomeViewModel) -> some View {
    Text("\(greeting), \(viewModel.currentUserDisplayName)!")
      .font(DSTheme.Typography.largeTitle)
      .foregroundColor(DSTheme.Colors.textPrimary)
      .frame(maxWidth: .infinity, alignment: .leading)
  }

  private var greeting: String {
    let hour = Calendar.current.component(.hour, from: Date())
    switch hour {
    case 0..<12:
      return "Good morning"
    case 12..<17:
      return "Good afternoon"
    default:
      return "Good evening"
    }
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return HomeView()
    .environment(authManager)
}
