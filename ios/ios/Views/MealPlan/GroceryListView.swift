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
  @State private var showingEditQuantityNeededForItemID: String?
  @State private var quantityNeededMinText: String = ""
  @State private var quantityNeededMaxText: String = ""
  @State private var showingStatusMenuForItemID: String?

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
        EnhancedGroceryItemRow(
          item: item,
          viewModel: viewModel,
          showingQuantityInput: showingQuantityInputForItemID == item.id,
          quantityInputText: $quantityInputText,
          showingEditQuantityNeeded: showingEditQuantityNeededForItemID == item.id,
          quantityNeededMinText: $quantityNeededMinText,
          quantityNeededMaxText: $quantityNeededMaxText,
          showingStatusMenu: showingStatusMenuForItemID == item.id,
          onStatusMenuToggle: {
            if showingStatusMenuForItemID == item.id {
              showingStatusMenuForItemID = nil
            } else {
              showingStatusMenuForItemID = item.id
            }
          },
          onMarkAsAcquired: {
            Task {
              await viewModel.markAsAcquired(item)
              showingStatusMenuForItemID = nil
            }
          },
          onMarkAsAlreadyOwned: {
            Task {
              await viewModel.markAsAlreadyOwned(item)
              showingStatusMenuForItemID = nil
            }
          },
          onMarkAsNeeds: {
            showingQuantityInputForItemID = item.id
            quantityInputText = item.hasQuantityPurchased ? String(item.quantityPurchased) : ""
            showingStatusMenuForItemID = nil
          },
          onMarkAsUnavailable: {
            Task {
              await viewModel.markAsUnavailable(item)
              showingStatusMenuForItemID = nil
            }
          },
          onEditQuantityNeeded: {
            showingEditQuantityNeededForItemID = item.id
            quantityNeededMinText = item.hasQuantityNeeded ? String(item.quantityNeeded.min) : ""
            quantityNeededMaxText =
              item.hasQuantityNeeded && item.quantityNeeded.hasMax
              ? String(item.quantityNeeded.max) : ""
            showingStatusMenuForItemID = nil
          },
          onQuantityPurchasedSubmit: {
            Task {
              await handleQuantityPurchasedSubmit(for: item)
            }
          },
          onQuantityPurchasedCancel: {
            showingQuantityInputForItemID = nil
            quantityInputText = ""
          },
          onQuantityNeededSubmit: {
            Task {
              await handleQuantityNeededSubmit(for: item)
            }
          },
          onQuantityNeededCancel: {
            showingEditQuantityNeededForItemID = nil
            quantityNeededMinText = ""
            quantityNeededMaxText = ""
          }
        )
      }
    }
  }

  private func handleQuantityPurchasedSubmit(for item: Mealplanning_MealPlanGroceryListItem) async {
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

  private func handleQuantityNeededSubmit(for item: Mealplanning_MealPlanGroceryListItem) async {
    guard let min = Float(quantityNeededMinText), min > 0 else {
      // Invalid input, cancel
      showingEditQuantityNeededForItemID = nil
      quantityNeededMinText = ""
      quantityNeededMaxText = ""
      return
    }

    let max: Float? = quantityNeededMaxText.isEmpty ? nil : Float(quantityNeededMaxText)

    await viewModel.updateQuantityNeeded(item, min: min, max: max)

    showingEditQuantityNeededForItemID = nil
    quantityNeededMinText = ""
    quantityNeededMaxText = ""
  }
}

// MARK: - Enhanced Grocery Item Row

struct EnhancedGroceryItemRow: View {
  let item: Mealplanning_MealPlanGroceryListItem
  let viewModel: GroceryListViewModel
  let showingQuantityInput: Bool
  @Binding var quantityInputText: String
  let showingEditQuantityNeeded: Bool
  @Binding var quantityNeededMinText: String
  @Binding var quantityNeededMaxText: String
  let showingStatusMenu: Bool
  let onStatusMenuToggle: () -> Void
  let onMarkAsAcquired: () -> Void
  let onMarkAsAlreadyOwned: () -> Void
  let onMarkAsNeeds: () -> Void
  let onMarkAsUnavailable: () -> Void
  let onEditQuantityNeeded: () -> Void
  let onQuantityPurchasedSubmit: () -> Void
  let onQuantityPurchasedCancel: () -> Void
  let onQuantityNeededSubmit: () -> Void
  let onQuantityNeededCancel: () -> Void

  var body: some View {
    VStack(spacing: 0) {
      // Main row
      HStack(spacing: 12) {
        // Item info
        VStack(alignment: .leading, spacing: 4) {
          Text(item.ingredient.name)
            .font(.body)
            .fontWeight(.medium)

          // Quantity needed (editable)
          if item.hasQuantityNeeded {
            Button {
              onEditQuantityNeeded()
            } label: {
              HStack(spacing: 4) {
                Text(formatQuantityNeeded(item.quantityNeeded))
                  .font(.caption)
                  .foregroundColor(.secondary)
                Image(systemName: "pencil")
                  .font(.caption2)
                  .foregroundColor(.secondary)
              }
            }
            .buttonStyle(.plain)
          }

          // Show quantity purchased if item is marked as have
          if (item.status == .acquired || item.status == .alreadyOwned) && item.hasQuantityPurchased
          {
            Text("Have: \(formatQuantity(item.quantityPurchased, unit: item.measurementUnit.name))")
              .font(.caption)
              .foregroundColor(.green)
          }

          // Status badge
          statusBadge
        }

        Spacer()

        // Status menu button
        Menu {
          Button {
            onMarkAsAcquired()
          } label: {
            Label(
              "Mark as Acquired",
              systemImage: item.status == .acquired ? "checkmark.circle.fill" : "checkmark.circle")
          }

          Button {
            onMarkAsAlreadyOwned()
          } label: {
            Label(
              "Mark as Already Owned",
              systemImage: item.status == .alreadyOwned ? "checkmark.circle.fill" : "house.fill")
          }

          Button {
            onMarkAsNeeds()
          } label: {
            Label(
              "Mark as Needs",
              systemImage: item.status == .needs ? "checkmark.circle.fill" : "cart.fill")
          }

          Button {
            onMarkAsUnavailable()
          } label: {
            Label(
              "Mark as Unavailable",
              systemImage: item.status == .unavailable ? "checkmark.circle.fill" : "xmark.circle")
          }

          Divider()

          Button {
            onEditQuantityNeeded()
          } label: {
            Label("Edit Quantity Needed", systemImage: "pencil")
          }
        } label: {
          Image(systemName: "ellipsis.circle")
            .font(.title3)
            .foregroundColor(.secondary)
        }
      }
      .padding()
      .background(Color(.systemGray6))
      .cornerRadius(8)

      // Quantity purchased input (shown when "Mark as Needs" is clicked)
      if showingQuantityInput {
        quantityPurchasedInputSection
      }

      // Quantity needed edit (shown when "Edit Quantity Needed" is clicked)
      if showingEditQuantityNeeded {
        quantityNeededEditSection
      }
    }
    .animation(.easeInOut(duration: 0.2), value: showingQuantityInput)
    .animation(.easeInOut(duration: 0.2), value: showingEditQuantityNeeded)
  }

  private var statusBadge: some View {
    HStack(spacing: 4) {
      Image(systemName: statusIcon)
        .font(.caption2)
      Text(statusText)
        .font(.caption2)
    }
    .padding(.horizontal, 6)
    .padding(.vertical, 2)
    .background(statusColor.opacity(0.2))
    .foregroundColor(statusColor)
    .cornerRadius(4)
  }

  private var statusIcon: String {
    switch item.status {
    case .acquired:
      return "checkmark.circle.fill"
    case .alreadyOwned:
      return "house.fill"
    case .needs:
      return "cart.fill"
    case .unavailable:
      return "xmark.circle.fill"
    default:
      return "questionmark.circle"
    }
  }

  private var statusText: String {
    switch item.status {
    case .acquired:
      return "Acquired"
    case .alreadyOwned:
      return "Already Owned"
    case .needs:
      return "Needs"
    case .unavailable:
      return "Unavailable"
    default:
      return "Unknown"
    }
  }

  private var statusColor: Color {
    switch item.status {
    case .acquired:
      return .green
    case .alreadyOwned:
      return .blue
    case .needs:
      return .orange
    case .unavailable:
      return .red
    default:
      return .gray
    }
  }

  private var quantityPurchasedInputSection: some View {
    VStack(spacing: 12) {
      Divider()

      VStack(alignment: .leading, spacing: 8) {
        Text("Quantity Purchased")
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
            onQuantityPurchasedCancel()
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
            onQuantityPurchasedSubmit()
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

  private var quantityNeededEditSection: some View {
    VStack(spacing: 12) {
      Divider()

      VStack(alignment: .leading, spacing: 8) {
        Text("Edit Quantity Needed")
          .font(.subheadline)
          .fontWeight(.medium)

        VStack(alignment: .leading, spacing: 8) {
          HStack(spacing: 8) {
            Text("Min:")
              .font(.caption)
              .foregroundColor(.secondary)
              .frame(width: 40, alignment: .leading)

            TextField("0", text: $quantityNeededMinText)
              .keyboardType(.decimalPad)
              .textFieldStyle(.roundedBorder)
              .onChange(of: quantityNeededMinText) { _, newValue in
                var filtered = newValue.filter { $0.isNumber || $0 == "." }
                let parts = filtered.split(separator: ".", omittingEmptySubsequences: false)
                if parts.count > 2 {
                  filtered = parts[0] + "." + parts.dropFirst().joined()
                }
                if filtered != newValue {
                  quantityNeededMinText = filtered
                }
              }

            if !item.measurementUnit.name.isEmpty {
              Text(item.measurementUnit.name)
                .font(.caption)
                .foregroundColor(.secondary)
            }
          }

          HStack(spacing: 8) {
            Text("Max:")
              .font(.caption)
              .foregroundColor(.secondary)
              .frame(width: 40, alignment: .leading)

            TextField("Optional", text: $quantityNeededMaxText)
              .keyboardType(.decimalPad)
              .textFieldStyle(.roundedBorder)
              .onChange(of: quantityNeededMaxText) { _, newValue in
                var filtered = newValue.filter { $0.isNumber || $0 == "." }
                let parts = filtered.split(separator: ".", omittingEmptySubsequences: false)
                if parts.count > 2 {
                  filtered = parts[0] + "." + parts.dropFirst().joined()
                }
                if filtered != newValue {
                  quantityNeededMaxText = filtered
                }
              }

            if !item.measurementUnit.name.isEmpty {
              Text(item.measurementUnit.name)
                .font(.caption)
                .foregroundColor(.secondary)
            }
          }
        }

        HStack(spacing: 12) {
          Button {
            onQuantityNeededCancel()
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
            onQuantityNeededSubmit()
          } label: {
            Text("Save")
              .frame(maxWidth: .infinity)
              .padding(.vertical, 8)
              .background(Color.blue)
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
