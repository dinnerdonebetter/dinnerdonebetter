//
//  MealListView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct MealListView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: MealListViewModel?
  @State private var searchQuery: String = ""

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          DSContentState(
            isLoading: viewModel.isLoading,
            loadingMessage: "Loading meals...",
            error: viewModel.errorMessage,
            onRetry: { await viewModel.loadMeals() },
            content: {
              let displayedMeals = viewModel.displayedMeals
              let isSearching = viewModel.isSearching
              let hasSearchQuery = !searchQuery.trimmingCharacters(in: .whitespacesAndNewlines)
                .isEmpty

              if isSearching {
                DSLoadingView("Searching meals...")
              } else if displayedMeals.isEmpty {
                if hasSearchQuery {
                  VStack(spacing: DSTheme.Spacing.lg) {
                    DSEmptyState(
                      icon: "magnifyingglass",
                      title: "No Results",
                      message: "No meals found matching \"\(searchQuery)\""
                    )

                    if let searchError = viewModel.searchError {
                      Text(searchError)
                        .font(DSTheme.Typography.caption)
                        .foregroundColor(DSTheme.Colors.error)
                        .padding(.horizontal)
                    }
                  }
                } else {
                  DSEmptyState(
                    icon: "fork.knife",
                    title: "No Meals",
                    message: "No meals found. Create some meals to get started.",
                    size: .large
                  )
                }
              } else {
                ScrollView {
                  LazyVStack(spacing: DSTheme.Spacing.md) {
                    ForEach(displayedMeals, id: \.id) { meal in
                      NavigationLink(destination: MealDetailView(mealID: meal.id)) {
                        MealCard(meal: meal)
                      }
                    }
                  }
                  .dsScreenPadding()
                }
              }
            })
        } else {
          DSInitializingView()
        }
      }
      .navigationTitle("Meals")
      .navigationBarTitleDisplayMode(.large)
      .searchable(text: $searchQuery, prompt: "Search meals...")
      .onChange(of: searchQuery) { _, newValue in
        if let viewModel = viewModel {
          viewModel.searchMeals(query: newValue)
        }
      }
      .refreshable {
        if let viewModel = viewModel {
          searchQuery = ""
          viewModel.searchResults = []
          await viewModel.loadMeals()
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = MealListViewModel(authManager: authManager)
        }
        if let viewModel = viewModel {
          Task {
            await viewModel.loadMeals()
          }
        }
      }
    }
  }
}

// MARK: - Meal Card

struct MealCard: View {
  let meal: Mealplanning_Meal

  var body: some View {
    DSCard(style: .outlined) {
      VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
        Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
          .font(DSTheme.Typography.label)
          .foregroundColor(DSTheme.Colors.textPrimary)
          .frame(maxWidth: .infinity, alignment: .leading)

        if !meal.description_p.isEmpty {
          Text(meal.description_p)
            .font(DSTheme.Typography.body)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .lineLimit(2)
            .frame(maxWidth: .infinity, alignment: .leading)
        }

        // Meal metadata
        HStack(spacing: DSTheme.Spacing.md) {
          if !meal.components.isEmpty {
            let recipeNames = meal.components.compactMap { component -> String? in
              component.recipe.name.isEmpty ? nil : component.recipe.name
            }
            if !recipeNames.isEmpty {
              Label(
                recipeNames.count == 1 ? recipeNames[0] : "\(recipeNames.count) recipes",
                systemImage: "book.closed"
              )
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
              .lineLimit(1)
            }
          }

          if meal.hasEstimatedPortions {
            Label("\(formatPortions(meal.estimatedPortions))", systemImage: "person.2")
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
        }
      }
      .frame(maxWidth: .infinity, alignment: .leading)
    }
  }

  private func formatPortions(_ range: Common_Float32RangeWithOptionalMax) -> String {
    if range.hasMax {
      if range.min == range.max {
        return String(format: "%.1f", range.min)
      } else {
        return String(format: "%.1f-%.1f", range.min, range.max)
      }
    } else {
      return String(format: "%.1f+", range.min)
    }
  }
}

// MARK: - Preview

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "Test User"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return MealListView()
    .environment(authManager)
}
