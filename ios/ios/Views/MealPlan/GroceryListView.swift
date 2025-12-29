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
  @State private var editingItem: Mealplanning_MealPlanGroceryListItem?
  @State private var showingQuantityEditor = false
  @State private var quantityMinText = ""
  @State private var quantityMaxText = ""
  @State private var quantityPurchasedText = ""

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
    .sheet(item: $editingItem) { item in
      GroceryItemEditSheet(
        item: item,
        onSave: { status, quantityMin, quantityMax, quantityPurchased in
          await handleItemUpdate(
            item: item,
            status: status,
            quantityMin: quantityMin,
            quantityMax: quantityMax,
            quantityPurchased: quantityPurchased
          )
        }
      )
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
        EditableGroceryItemRow(
          item: item,
          onTap: {
            editingItem = item
          },
          onCheck: {
            Task {
              await viewModel.markAsAcquired(item)
            }
          },
          onMarkOwned: {
            Task {
              await viewModel.markAsAlreadyOwned(item)
            }
          }
        )
      }
    }
  }

  private func handleItemUpdate(
    item: Mealplanning_MealPlanGroceryListItem,
    status: Mealplanning_MealPlanGroceryListItemStatus?,
    quantityMin: Float?,
    quantityMax: Float?,
    quantityPurchased: Float?
  ) async {
    await viewModel.updateItem(
      item,
      status: status,
      quantityNeededMin: quantityMin,
      quantityNeededMax: quantityMax,
      quantityPurchased: quantityPurchased
    )
    editingItem = nil
  }
}

// MARK: - Editable Grocery Item Row

struct EditableGroceryItemRow: View {
  let item: Mealplanning_MealPlanGroceryListItem
  let onTap: () -> Void
  let onCheck: () -> Void
  let onMarkOwned: () -> Void

  var body: some View {
    HStack(spacing: 12) {
      // Checkbox/Status indicator
      Button(action: onCheck) {
        Image(systemName: item.status == .acquired ? "checkmark.circle.fill" : "circle")
          .foregroundColor(item.status == .acquired ? .green : .gray)
          .font(.title3)
      }
      .buttonStyle(.plain)

      // Item details
      VStack(alignment: .leading, spacing: 4) {
        Text(item.ingredient.name)
          .font(.body)
          .fontWeight(.medium)

        // Quantity info
        HStack(spacing: 8) {
          if item.hasQuantityNeeded {
            Text(formatQuantityNeeded(item.quantityNeeded))
              .font(.caption)
              .foregroundColor(.secondary)
          }

          if item.hasQuantityPurchased {
            Text(
              "• Purchased: \(formatQuantity(item.quantityPurchased, unit: item.measurementUnit.name))"
            )
            .font(.caption)
            .foregroundColor(.secondary)
          }
        }

        // Status badge
        if item.status != .needs {
          Text(statusText(item.status))
            .font(.caption2)
            .padding(.horizontal, 6)
            .padding(.vertical, 2)
            .background(statusColor(item.status).opacity(0.2))
            .foregroundColor(statusColor(item.status))
            .cornerRadius(4)
        }
      }

      Spacer()

      // Action buttons
      Menu {
        Button(action: onMarkOwned) {
          Label("Mark as Already Owned", systemImage: "house.fill")
        }

        Button(action: onTap) {
          Label("Edit Quantity", systemImage: "pencil")
        }
      } label: {
        Image(systemName: "ellipsis")
          .foregroundColor(.secondary)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(8)
    .contentShape(Rectangle())
    .onTapGesture {
      onTap()
    }
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

  private func statusText(_ status: Mealplanning_MealPlanGroceryListItemStatus) -> String {
    switch status {
    case .alreadyOwned:
      return "Already Owned"
    case .acquired:
      return "Acquired"
    case .unavailable:
      return "Unavailable"
    case .needs:
      return "Needed"
    default:
      return ""
    }
  }

  private func statusColor(_ status: Mealplanning_MealPlanGroceryListItemStatus) -> Color {
    switch status {
    case .alreadyOwned:
      return .blue
    case .acquired:
      return .green
    case .unavailable:
      return .red
    case .needs:
      return .orange
    default:
      return .gray
    }
  }
}

// MARK: - Grocery Item Edit Sheet

struct GroceryItemEditSheet: View {
  let item: Mealplanning_MealPlanGroceryListItem
  let onSave:
    (
      Mealplanning_MealPlanGroceryListItemStatus?,
      Float?,
      Float?,
      Float?
    ) async -> Void

  @Environment(\.dismiss) private var dismiss
  @State private var selectedStatus: Mealplanning_MealPlanGroceryListItemStatus?
  @State private var quantityMinText: String
  @State private var quantityMaxText: String
  @State private var quantityPurchasedText: String
  @State private var isSaving = false

  init(
    item: Mealplanning_MealPlanGroceryListItem,
    onSave:
      @escaping (Mealplanning_MealPlanGroceryListItemStatus?, Float?, Float?, Float?) async -> Void
  ) {
    self.item = item
    self.onSave = onSave
    _selectedStatus = State(initialValue: item.status)
    _quantityMinText = State(
      initialValue: item.hasQuantityNeeded ? String(item.quantityNeeded.min) : "")
    _quantityMaxText = State(
      initialValue: item.hasQuantityNeeded && item.quantityNeeded.hasMax
        ? String(item.quantityNeeded.max) : "")
    _quantityPurchasedText = State(
      initialValue: item.hasQuantityPurchased ? String(item.quantityPurchased) : "")
  }

  var body: some View {
    NavigationView {
      Form {
        Section("Status") {
          Picker("Status", selection: $selectedStatus) {
            Text("Needed").tag(
              Mealplanning_MealPlanGroceryListItemStatus.needs
                as Mealplanning_MealPlanGroceryListItemStatus?)
            Text("Already Owned").tag(
              Mealplanning_MealPlanGroceryListItemStatus.alreadyOwned
                as Mealplanning_MealPlanGroceryListItemStatus?)
            Text("Acquired").tag(
              Mealplanning_MealPlanGroceryListItemStatus.acquired
                as Mealplanning_MealPlanGroceryListItemStatus?)
            Text("Unavailable").tag(
              Mealplanning_MealPlanGroceryListItemStatus.unavailable
                as Mealplanning_MealPlanGroceryListItemStatus?)
          }
        }

        Section("Quantity Needed") {
          HStack {
            Text("Min")
            TextField("Min", text: $quantityMinText)
              .keyboardType(.decimalPad)
          }

          HStack {
            Text("Max (optional)")
            TextField("Max", text: $quantityMaxText)
              .keyboardType(.decimalPad)
          }
        }

        Section("Quantity Purchased") {
          TextField("Quantity", text: $quantityPurchasedText)
            .keyboardType(.decimalPad)
        }

        Section {
          Text("Ingredient: \(item.ingredient.name)")
            .font(.caption)
            .foregroundColor(.secondary)

          if !item.measurementUnit.name.isEmpty {
            Text("Unit: \(item.measurementUnit.name)")
              .font(.caption)
              .foregroundColor(.secondary)
          }
        }
      }
      .navigationTitle("Edit Item")
      .navigationBarTitleDisplayMode(.inline)
      .toolbar {
        ToolbarItem(placement: .cancellationAction) {
          Button("Cancel") {
            dismiss()
          }
        }

        ToolbarItem(placement: .confirmationAction) {
          Button("Save") {
            Task {
              await save()
            }
          }
          .disabled(isSaving)
        }
      }
    }
  }

  private func save() async {
    isSaving = true

    let quantityMin = Float(quantityMinText)
    let quantityMax = quantityMaxText.isEmpty ? nil : Float(quantityMaxText)
    let quantityPurchased = quantityPurchasedText.isEmpty ? nil : Float(quantityPurchasedText)

    await onSave(selectedStatus, quantityMin, quantityMax, quantityPurchased)
    isSaving = false
    dismiss()
  }
}

// Make MealPlanGroceryListItem Identifiable for sheet
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
