//
//  GroceryListView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct GroceryListView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: GroceryListViewModel
  @State private var showingQuantityInputForItemID: String?
  @State private var quantityInputText: String = ""

  init(
    mealPlan: Mealplanning_MealPlan, items: [Mealplanning_MealPlanGroceryListItem],
    authManager: AuthenticationManager
  ) {
    _viewModel = State(
      initialValue: GroceryListViewModel(
        mealPlan: mealPlan,
        items: items,
        authManager: authManager
      )
    )
  }

  var body: some View {
    ScrollView {
      VStack(alignment: .leading, spacing: 20) {
        // Header
        headerSection

        if viewModel.isLoading {
          ProgressView("Loading grocery list...")
            .frame(maxWidth: .infinity, alignment: .center)
            .padding()
        } else if viewModel.items.isEmpty {
          emptyStateView
        } else {
          // Items by status
          if !viewModel.needsItems.isEmpty {
            itemsSection(
              title: "Needed",
              items: viewModel.needsItems,
              color: .orange
            )
          }

          if !viewModel.alreadyOwnedItems.isEmpty {
            itemsSection(
              title: "Already Owned",
              items: viewModel.alreadyOwnedItems,
              color: .blue
            )
          }

          if !viewModel.acquiredItems.isEmpty {
            itemsSection(
              title: "Acquired",
              items: viewModel.acquiredItems,
              color: .green
            )
          }

          if !viewModel.unavailableItems.isEmpty {
            itemsSection(
              title: "Unavailable",
              items: viewModel.unavailableItems,
              color: .red
            )
          }
        }

        if let errorMessage = viewModel.errorMessage {
          Text(errorMessage)
            .foregroundColor(.red)
            .padding()
        }
      }
      .padding()
    }
    .navigationTitle(viewModel.mealPlan.notes.isEmpty ? "Grocery List" : viewModel.mealPlan.notes)
    .navigationBarTitleDisplayMode(.large)
    .task {
      await viewModel.loadItems()
    }
  }

  private var headerSection: some View {
    VStack(alignment: .leading, spacing: 8) {
      Text(HomeView.formatMealPlanTimeRange(viewModel.mealPlan))
        .font(.subheadline)
        .foregroundColor(.secondary)

      Text("\(viewModel.items.count) item\(viewModel.items.count == 1 ? "" : "s")")
        .font(.caption)
        .foregroundColor(.secondary)
    }
  }

  private var emptyStateView: some View {
    VStack(spacing: 16) {
      Image(systemName: "cart")
        .font(.system(size: 48))
        .foregroundColor(.secondary)

      Text("No grocery items")
        .font(.headline)
        .foregroundColor(.secondary)

      Text("Grocery list items will appear here once the meal plan is finalized.")
        .font(.subheadline)
        .foregroundColor(.secondary)
        .multilineTextAlignment(.center)
    }
    .frame(maxWidth: .infinity)
    .padding(.vertical, 40)
  }

  private func itemsSection(
    title: String,
    items: [Mealplanning_MealPlanGroceryListItem],
    color: Color
  ) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text(title)
        .font(.headline)
        .foregroundColor(color)

      ForEach(items, id: \.id) { item in
        SimpleGroceryItemRow(
          item: item,
          viewModel: viewModel,
          showingQuantityInput: showingQuantityInputForItemID == item.id,
          quantityInputText: $quantityInputText,
          onMarkAsHave: {
            Task {
              await viewModel.markAsAcquired(item)
            }
          },
          onMarkAsNeed: {
            showingQuantityInputForItemID = item.id
            quantityInputText = item.hasQuantityPurchased ? String(item.quantityPurchased) : ""
          },
          onQuantitySubmit: {
            Task {
              await handleQuantitySubmit(for: item)
            }
          },
          onQuantityCancel: {
            showingQuantityInputForItemID = nil
            quantityInputText = ""
          }
        )
      }
    }
  }

  private func handleQuantitySubmit(for item: Mealplanning_MealPlanGroceryListItem) async {
    guard let quantity = Float(quantityInputText), quantity > 0 else {
      // If no valid quantity, just mark as needed
      await viewModel.markAsNeeds(item)
      showingQuantityInputForItemID = nil
      quantityInputText = ""
      return
    }

    // Mark as needed and set quantity purchased
    await viewModel.updateItem(
      item,
      status: .needs,
      quantityNeededMin: nil,
      quantityNeededMax: nil,
      quantityPurchased: quantity
    )

    showingQuantityInputForItemID = nil
    quantityInputText = ""
  }
}

// MARK: - Simple Grocery Item Row

struct SimpleGroceryItemRow: View {
  let item: Mealplanning_MealPlanGroceryListItem
  let viewModel: GroceryListViewModel
  let showingQuantityInput: Bool
  @Binding var quantityInputText: String
  let onMarkAsHave: () -> Void
  let onMarkAsNeed: () -> Void
  let onQuantitySubmit: () -> Void
  let onQuantityCancel: () -> Void

  var body: some View {
    VStack(spacing: 0) {
      // Main row
      HStack(spacing: 12) {
        // Item info
        VStack(alignment: .leading, spacing: 4) {
          Text(item.ingredient.name)
            .font(.body)
            .fontWeight(.medium)

          // Quantity needed (read-only)
          if item.hasQuantityNeeded {
            Text(formatQuantityNeeded(item.quantityNeeded))
              .font(.caption)
              .foregroundColor(.secondary)
          }

          // Show quantity purchased if item is marked as have
          if (item.status == .acquired || item.status == .alreadyOwned) && item.hasQuantityPurchased
          {
            Text("Have: \(formatQuantity(item.quantityPurchased, unit: item.measurementUnit.name))")
              .font(.caption)
              .foregroundColor(.green)
          }
        }

        Spacer()

        // Action buttons
        HStack(spacing: 8) {
          // "I have this" button
          Button {
            onMarkAsHave()
          } label: {
            Text("I have this")
              .font(.subheadline)
              .fontWeight(.medium)
              .padding(.horizontal, 16)
              .padding(.vertical, 8)
              .background(
                (item.status == .acquired || item.status == .alreadyOwned)
                  ? Color.green.opacity(0.2) : Color(.systemGray5)
              )
              .foregroundColor(
                (item.status == .acquired || item.status == .alreadyOwned) ? .green : .primary
              )
              .cornerRadius(8)
          }
          .buttonStyle(.plain)

          // "I need this" button
          Button {
            onMarkAsNeed()
          } label: {
            Text("I need this")
              .font(.subheadline)
              .fontWeight(.medium)
              .padding(.horizontal, 16)
              .padding(.vertical, 8)
              .background(
                item.status == .needs ? Color.orange.opacity(0.2) : Color(.systemGray5)
              )
              .foregroundColor(item.status == .needs ? .orange : .primary)
              .cornerRadius(8)
          }
          .buttonStyle(.plain)
        }
      }
      .padding()
      .background(Color(.systemGray6))
      .cornerRadius(8)

      // Quantity input (shown when "I need this" is clicked)
      if showingQuantityInput {
        VStack(spacing: 12) {
          Divider()

          VStack(alignment: .leading, spacing: 8) {
            Text("How much do you need?")
              .font(.subheadline)
              .fontWeight(.medium)

            HStack(spacing: 12) {
              TextField("0", text: $quantityInputText)
                .keyboardType(.decimalPad)
                .textFieldStyle(.roundedBorder)
                .onChange(of: quantityInputText) { _, newValue in
                  // Filter to only allow numbers and a single decimal point
                  var filtered = newValue.filter { $0.isNumber || $0 == "." }
                  // Ensure only one decimal point
                  let parts = filtered.split(separator: ".", omittingEmptySubsequences: false)
                  if parts.count > 2 {
                    filtered = parts[0] + "." + parts.dropFirst().joined()
                  }
                  if filtered != newValue {
                    quantityInputText = filtered
                  }
                }

              if !item.measurementUnit.name.isEmpty {
                Text(item.measurementUnit.name)
                  .font(.subheadline)
                  .foregroundColor(.secondary)
              }
            }

            HStack(spacing: 12) {
              Button {
                onQuantityCancel()
              } label: {
                Text("Cancel")
                  .frame(maxWidth: .infinity)
                  .padding(.vertical, 8)
                  .background(Color(.systemGray5))
                  .foregroundColor(.primary)
                  .cornerRadius(8)
              }
              .buttonStyle(.plain)

              Button {
                onQuantitySubmit()
              } label: {
                Text("Save")
                  .frame(maxWidth: .infinity)
                  .padding(.vertical, 8)
                  .background(Color.orange)
                  .foregroundColor(.white)
                  .cornerRadius(8)
              }
              .buttonStyle(.plain)
            }
          }
          .padding()
        }
        .background(Color(.systemGray5))
        .cornerRadius(8)
      }
    }
    .animation(.easeInOut(duration: 0.2), value: showingQuantityInput)
  }

  private func formatQuantityNeeded(_ quantity: Common_Float32RangeWithOptionalMax) -> String {
    let unit = item.measurementUnit.name.isEmpty ? "" : " \(item.measurementUnit.name)"
    if quantity.hasMax {
      if quantity.min == quantity.max {
        return "\(quantity.min)\(unit) needed"
      } else {
        return "\(quantity.min) - \(quantity.max)\(unit) needed"
      }
    } else {
      return "\(quantity.min)+\(unit) needed"
    }
  }

  private func formatQuantity(_ quantity: Float, unit: String) -> String {
    let unitText = unit.isEmpty ? "" : " \(unit)"
    return "\(quantity)\(unitText)"
  }
}

// Make MealPlanGroceryListItem Identifiable
extension Mealplanning_MealPlanGroceryListItem: Identifiable {
  // Already has id property, so this extension just makes it conform to Identifiable
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  // Create a sample meal plan and items
  var mealPlan = Mealplanning_MealPlan()
  mealPlan.id = "mealplan123"
  mealPlan.notes = "Sample Meal Plan"

  var item = Mealplanning_MealPlanGroceryListItem()
  item.id = "item123"
  var ingredient = Mealplanning_ValidIngredient()
  ingredient.name = "Chicken Breast"
  item.ingredient = ingredient
  var quantity = Common_Float32RangeWithOptionalMax()
  quantity.min = 2.0
  quantity.max = 4.0
  item.quantityNeeded = quantity
  item.status = .needs

  return NavigationView {
    GroceryListView(
      mealPlan: mealPlan,
      items: [item],
      authManager: authManager
    )
  }
  .environment(authManager)
}
