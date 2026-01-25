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
          if viewModel.isLoading {
            ProgressView("Loading meals...")
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
                  await viewModel.loadMeals()
                }
              }
              .buttonStyle(.borderedProminent)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
          } else {
            let displayedMeals = viewModel.displayedMeals
            let isSearching = viewModel.isSearching
            let hasSearchQuery = !searchQuery.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty

            if isSearching {
              ProgressView("Searching meals...")
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if displayedMeals.isEmpty {
              if hasSearchQuery {
                VStack(spacing: 16) {
                  Image(systemName: "magnifyingglass")
                    .font(.largeTitle)
                    .foregroundColor(.secondary)
                  Text("No Results")
                    .font(.headline)
                  Text("No meals found matching \"\(searchQuery)\"")
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
                  Image(systemName: "fork.knife")
                    .font(.largeTitle)
                    .foregroundColor(.secondary)
                  Text("No Meals")
                    .font(.headline)
                  Text("No meals found. Create some meals to get started.")
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
                  ForEach(displayedMeals, id: \.id) { meal in
                    NavigationLink(destination: MealDetailView(mealID: meal.id)) {
                      MealCard(meal: meal)
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
      .navigationTitle("Meals")
      .navigationBarTitleDisplayMode(.large)
      .searchable(text: $searchQuery, prompt: "Search meals...")
      .onChange(of: searchQuery) { oldValue, newValue in
        if let viewModel = viewModel {
          viewModel.searchMeals(query: newValue)
        }
      }
      .refreshable {
        if let viewModel = viewModel {
          // Clear search when refreshing
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
    VStack(alignment: .leading, spacing: 8) {
      Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
        .font(.headline)
        .foregroundColor(.primary)
        .frame(maxWidth: .infinity, alignment: .leading)

      if !meal.description_p.isEmpty {
        Text(meal.description_p)
          .font(.subheadline)
          .foregroundColor(.secondary)
          .lineLimit(2)
          .frame(maxWidth: .infinity, alignment: .leading)
      }

      // Meal metadata
      HStack(spacing: 12) {
        // Show recipe names from components
        if !meal.components.isEmpty {
          let recipeNames = meal.components.compactMap { component -> String? in
            component.recipe.name.isEmpty ? nil : component.recipe.name
          }
          if !recipeNames.isEmpty {
            Label(
              recipeNames.count == 1 ? recipeNames[0] : "\(recipeNames.count) recipes",
              systemImage: "book.closed"
            )
            .font(.caption)
            .foregroundColor(.secondary)
            .lineLimit(1)
          }
        }

        if meal.hasEstimatedPortions {
          Label("\(formatPortions(meal.estimatedPortions))", systemImage: "person.2")
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

  return MealListView()
    .environment(authManager)
}
