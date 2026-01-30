//
//  RecipeListView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct RecipeListView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: RecipeListViewModel?
  @State private var searchQuery: String = ""

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          DSContentState(
            isLoading: viewModel.isLoading,
            loadingMessage: "Loading recipes...",
            error: viewModel.errorMessage,
            onRetry: { await viewModel.loadRecipes() },
            content: {
              let displayedRecipes = viewModel.displayedRecipes
              let isSearching = viewModel.isSearching
              let hasSearchQuery = !searchQuery.trimmingCharacters(in: .whitespacesAndNewlines)
                .isEmpty

              if isSearching {
                DSLoadingView("Searching recipes...")
              } else if displayedRecipes.isEmpty {
                if hasSearchQuery {
                  VStack(spacing: DSTheme.Spacing.lg) {
                    DSEmptyState(
                      icon: "magnifyingglass",
                      title: "No Results",
                      message: "No recipes found matching \"\(searchQuery)\""
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
                    icon: "book.closed",
                    title: "No Recipes",
                    message: "No recipes found. Create some recipes to get started.",
                    size: .large
                  )
                }
              } else {
                ScrollView {
                  LazyVStack(spacing: DSTheme.Spacing.md) {
                    ForEach(displayedRecipes, id: \.id) { recipe in
                      NavigationLink(destination: PerformRecipeView(recipeID: recipe.id)) {
                        RecipeCard(recipe: recipe)
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
      .navigationTitle("Recipes")
      .navigationBarTitleDisplayMode(.large)
      .searchable(text: $searchQuery, prompt: "Search recipes...")
      .onChange(of: searchQuery) { _, newValue in
        if let viewModel = viewModel {
          viewModel.searchRecipes(query: newValue)
        }
      }
      .refreshable {
        if let viewModel = viewModel {
          searchQuery = ""
          viewModel.searchResults = []
          await viewModel.loadRecipes()
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = RecipeListViewModel(authManager: authManager)
        }
        if let viewModel = viewModel {
          Task {
            await viewModel.loadRecipes()
          }
        }
      }
    }
  }
}

// MARK: - Recipe Card

struct RecipeCard: View {
  let recipe: Mealplanning_Recipe

  var body: some View {
    DSCard(style: .outlined) {
      VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
        Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
          .font(DSTheme.Typography.label)
          .foregroundColor(DSTheme.Colors.textPrimary)
          .frame(maxWidth: .infinity, alignment: .leading)

        if !recipe.description_p.isEmpty {
          Text(recipe.description_p)
            .font(DSTheme.Typography.body)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .lineLimit(2)
            .frame(maxWidth: .infinity, alignment: .leading)
        }

        // Recipe metadata
        HStack(spacing: DSTheme.Spacing.md) {
          if !recipe.steps.isEmpty {
            Label(
              "\(recipe.steps.count) step\(recipe.steps.count == 1 ? "" : "s")",
              systemImage: "list.number"
            )
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
          }

          if recipe.hasEstimatedPortions {
            Label("\(formatPortions(recipe.estimatedPortions))", systemImage: "person.2")
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

  return RecipeListView()
    .environment(authManager)
}
