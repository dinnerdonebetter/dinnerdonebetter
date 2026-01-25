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
          if viewModel.isLoading {
            ProgressView("Loading recipes...")
              .frame(maxWidth: .infinity, maxHeight: .infinity)
          } else if let errorMessage = viewModel.errorMessage {
            VStack(spacing: 16) {
              Image(systemName: "exclamationmark.triangle")
                .font(.largeTitle)
                .foregroundColor(.orange)
              Text("Error")
                .font(.headline)
              Text(errorMessage)
                .font(.subheadline)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
                .padding(.horizontal)
              Button("Retry") {
                Task {
                  await viewModel.loadRecipes()
                }
              }
              .buttonStyle(.borderedProminent)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
          } else {
            let displayedRecipes = viewModel.displayedRecipes
            let isSearching = viewModel.isSearching
            let hasSearchQuery = !searchQuery.trimmingCharacters(in: .whitespacesAndNewlines)
              .isEmpty

            if isSearching {
              ProgressView("Searching recipes...")
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if displayedRecipes.isEmpty {
              if hasSearchQuery {
                VStack(spacing: 16) {
                  Image(systemName: "magnifyingglass")
                    .font(.largeTitle)
                    .foregroundColor(.secondary)
                  Text("No Results")
                    .font(.headline)
                  Text("No recipes found matching \"\(searchQuery)\"")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
                    .padding(.horizontal)

                  if let searchError = viewModel.searchError {
                    Text(searchError)
                      .font(.caption)
                      .foregroundColor(.red)
                      .padding(.horizontal)
                  }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
              } else {
                VStack(spacing: 16) {
                  Image(systemName: "book.closed")
                    .font(.largeTitle)
                    .foregroundColor(.secondary)
                  Text("No Recipes")
                    .font(.headline)
                  Text("No recipes found. Create some recipes to get started.")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
                    .padding(.horizontal)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
              }
            } else {
              ScrollView {
                LazyVStack(spacing: 12) {
                  ForEach(displayedRecipes, id: \.id) { recipe in
                    NavigationLink(destination: PerformRecipeView(recipeID: recipe.id)) {
                      RecipeCard(recipe: recipe)
                    }
                  }
                }
                .padding()
              }
            }
          }
        } else {
          ProgressView("Initializing...")
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
      }
      .navigationTitle("Recipes")
      .navigationBarTitleDisplayMode(.large)
      .searchable(text: $searchQuery, prompt: "Search recipes...")
      .onChange(of: searchQuery) { oldValue, newValue in
        if let viewModel = viewModel {
          viewModel.searchRecipes(query: newValue)
        }
      }
      .refreshable {
        if let viewModel = viewModel {
          // Clear search when refreshing
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
    VStack(alignment: .leading, spacing: 8) {
      Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
        .font(.headline)
        .foregroundColor(.primary)
        .frame(maxWidth: .infinity, alignment: .leading)

      if !recipe.description_p.isEmpty {
        Text(recipe.description_p)
          .font(.subheadline)
          .foregroundColor(.secondary)
          .lineLimit(2)
          .frame(maxWidth: .infinity, alignment: .leading)
      }

      // Recipe metadata
      HStack(spacing: 12) {
        if !recipe.steps.isEmpty {
          Label(
            "\(recipe.steps.count) step\(recipe.steps.count == 1 ? "" : "s")",
            systemImage: "list.number"
          )
          .font(.caption)
          .foregroundColor(.secondary)
        }

        if recipe.hasEstimatedPortions {
          Label("\(formatPortions(recipe.estimatedPortions))", systemImage: "person.2")
            .font(.caption)
            .foregroundColor(.secondary)
        }
      }
    }
    .padding()
    .frame(maxWidth: .infinity, alignment: .leading)
    .background(Color(.systemBackground))
    .cornerRadius(12)
    .overlay(
      RoundedRectangle(cornerRadius: 12)
        .stroke(Color(.systemGray4), lineWidth: 1)
    )
  }

  private func formatPortions(_ range: Common_Float32RangeWithOptionalMax) -> String {
    if range.hasMax {
      if range.min == range.max {
        return String(format: "%.1f", range.min)
      } else {
        return String(format: "%.1f-%.1f", range.min, range.max)
      }
    } else {
      // min is always present, but max is optional
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
