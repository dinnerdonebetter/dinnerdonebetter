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
  @State private var showingSwipeReview = false
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
    menuView
      .navigationTitle(viewModel.mealPlan.notes.isEmpty ? "Grocery List" : viewModel.mealPlan.notes)
      .navigationBarTitleDisplayMode(.large)
      .toolbar {
        ToolbarItem(placement: .primaryAction) {
          Button {
            showingSwipeReview = true
          } label: {
            Label("Swipe review", systemImage: "hand.draw")
          }
          .disabled(viewModel.items.isEmpty)
        }
      }
      .fullScreenCover(isPresented: $showingSwipeReview) {
        SwipeReviewSheet(
          viewModel: viewModel,
          onDismiss: { showingSwipeReview = false }
        )
      }
      .task {
        await viewModel.loadItems()
      }
  }

  private var menuView: some View {
    ScrollView {
      LazyVStack(alignment: .leading, spacing: DSTheme.Spacing.xl, pinnedViews: [.sectionHeaders]) {
        headerSection
          .padding(.horizontal, DSTheme.Spacing.md)

        if viewModel.isLoading {
          DSLoadingView("Loading grocery list...")
            .padding(.horizontal, DSTheme.Spacing.md)
        } else if viewModel.items.isEmpty {
          emptyStateView
            .padding(.horizontal, DSTheme.Spacing.md)
        } else {
          if !viewModel.needsItems.isEmpty {
            Section {
              itemRows(items: viewModel.needsItems)
            } header: {
              sectionHeader(title: "Needed", color: .orange)
            }
          }

          if !viewModel.alreadyOwnedItems.isEmpty {
            Section {
              itemRows(items: viewModel.alreadyOwnedItems)
            } header: {
              sectionHeader(title: "Already Owned", color: .blue)
            }
          }

          if !viewModel.acquiredItems.isEmpty {
            Section {
              itemRows(items: viewModel.acquiredItems)
            } header: {
              sectionHeader(title: "Acquired", color: .green)
            }
          }

          if !viewModel.unavailableItems.isEmpty {
            Section {
              itemRows(items: viewModel.unavailableItems)
            } header: {
              sectionHeader(title: "Unavailable", color: .red)
            }
          }
        }

        if let errorMessage = viewModel.errorMessage {
          Text(errorMessage)
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.error)
            .padding()
        }
      }
      .padding(.vertical, DSTheme.Spacing.md)
    }
  }

  private func sectionHeader(title: String, color: Color) -> some View {
    Text(title)
      .font(DSTheme.Typography.label)
      .foregroundColor(color)
      .frame(maxWidth: .infinity, alignment: .leading)
      .padding(.horizontal, DSTheme.Spacing.md)
      .padding(.vertical, DSTheme.Spacing.sm)
      .background(.regularMaterial)
  }

  private var headerSection: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      Text(HomeView.formatMealPlanTimeRange(viewModel.mealPlan))
        .font(DSTheme.Typography.body)
        .foregroundColor(DSTheme.Colors.textSecondary)

      Text("\(viewModel.items.count) item\(viewModel.items.count == 1 ? "" : "s")")
        .font(DSTheme.Typography.caption)
        .foregroundColor(DSTheme.Colors.textSecondary)
    }
  }

  private var emptyStateView: some View {
    DSEmptyState(
      icon: "cart",
      title: "No grocery items",
      message: "Grocery list items will appear here once the meal plan is finalized."
    )
  }

  private func itemRows(
    items: [Mealplanning_MealPlanGroceryListItem]
  ) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
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
    .padding(.horizontal, DSTheme.Spacing.md)
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

// MARK: - Swipe Review Sheet

struct SwipeReviewSheet: View {
  var viewModel: GroceryListViewModel
  let onDismiss: () -> Void

  private var reviewItems: [Mealplanning_MealPlanGroceryListItem] {
    func sortPriority(_ status: Mealplanning_MealPlanGroceryListItemStatus) -> Int {
      switch status {
      case .needs: return 0
      case .unknown: return 1
      case .unavailable: return 2
      case .acquired, .alreadyOwned: return 3
      default: return 1
      }
    }
    return viewModel.items.sorted { item1, item2 in
      let priority1 = sortPriority(item1.status)
      let priority2 = sortPriority(item2.status)
      if priority1 != priority2 { return priority1 < priority2 }
      return item1.ingredient.name < item2.ingredient.name
    }
  }

  var body: some View {
    NavigationStack {
      List {
        Section {
          HStack {
            HStack(spacing: 4) {
              Image(systemName: "chevron.left")
                .font(.caption2)
              Text("Have")
                .font(.caption2)
            }
            .foregroundColor(.green)
            .frame(maxWidth: .infinity, alignment: .leading)

            HStack(spacing: 4) {
              Text("Need")
                .font(.caption2)
              Image(systemName: "chevron.right")
                .font(.caption2)
            }
            .foregroundColor(.orange)
            .frame(maxWidth: .infinity, alignment: .trailing)
          }
          .listRowBackground(Color.clear)
          .listRowInsets(EdgeInsets())
        }

        Section {
          ForEach(reviewItems, id: \.id) { item in
            ResistantSwipeReviewRow(
              item: item,
              onMarkAsHave: { Task { await viewModel.markAsAlreadyOwned(item) } },
              onMarkAsNeed: { Task { await viewModel.markAsNeeds(item) } }
            )
            .listRowInsets(EdgeInsets(top: 2, leading: 0, bottom: 2, trailing: 0))
          }
        }
      }
      .listStyle(.insetGrouped)
      .navigationTitle("Swipe review")
      .navigationBarTitleDisplayMode(.inline)
      .toolbar {
        ToolbarItem(placement: .confirmationAction) {
          Button("Done") {
            onDismiss()
          }
        }
      }
    }
  }
}

// MARK: - Resistant Swipe Review Row

private struct ResistantSwipeReviewRow: View {
  let item: Mealplanning_MealPlanGroceryListItem
  let onMarkAsHave: () -> Void
  let onMarkAsNeed: () -> Void

  private let swipeThreshold: CGFloat = 60
  private let resistanceFactor: CGFloat = 0.25

  @State private var dragOffset: CGFloat = 0
  @State private var hasTriggeredResistanceHaptic = false
  @State private var isHorizontalDrag: Bool?

  private var isAlreadyHave: Bool {
    item.status == .acquired || item.status == .alreadyOwned
  }

  private var isAlreadyNeed: Bool {
    item.status == .needs
  }

  private var appliedOffset: CGFloat {
    let raw = dragOffset
    if raw > 0 {
      if isAlreadyNeed {
        return raw * resistanceFactor
      }
      return raw
    } else {
      if isAlreadyHave {
        return raw * resistanceFactor
      }
      return raw
    }
  }

  var body: some View {
    GeometryReader { geometry in
      ZStack(alignment: .leading) {
        HStack {
          Color.orange
            .frame(width: swipeThreshold + 20)
            .overlay {
              HStack {
                Image(systemName: "cart.fill")
                Text("Need")
                  .font(.caption)
                  .fontWeight(.medium)
              }
              .foregroundColor(.white)
            }
          Spacer()
          Color.green
            .frame(width: swipeThreshold + 20)
            .overlay {
              HStack {
                Text("Have")
                  .font(.caption)
                  .fontWeight(.medium)
                Image(systemName: "checkmark.circle.fill")
              }
              .foregroundColor(.white)
            }
        }

        ZStack {
          Color(.secondarySystemGroupedBackground)
          swipeReviewRowContent
            .padding(.horizontal, 16)
        }
        .frame(width: geometry.size.width, height: geometry.size.height)
        .offset(x: appliedOffset)
        .simultaneousGesture(
          DragGesture(minimumDistance: 10)
            .onChanged { value in
              if isHorizontalDrag == nil {
                let horizontal = abs(value.translation.width)
                let vertical = abs(value.translation.height)
                isHorizontalDrag = horizontal > vertical
              }

              guard isHorizontalDrag == true else { return }

              let translation = value.translation.width
              dragOffset = translation

              if translation > 20 && isAlreadyNeed, !hasTriggeredResistanceHaptic {
                hasTriggeredResistanceHaptic = true
                let generator = UIImpactFeedbackGenerator(style: .light)
                generator.impactOccurred()
              } else if translation < -20 && isAlreadyHave, !hasTriggeredResistanceHaptic {
                hasTriggeredResistanceHaptic = true
                let generator = UIImpactFeedbackGenerator(style: .light)
                generator.impactOccurred()
              }
            }
            .onEnded { value in
              defer {
                isHorizontalDrag = nil
                hasTriggeredResistanceHaptic = false
              }

              guard isHorizontalDrag == true else { return }

              let translation = value.translation.width

              if translation > swipeThreshold, !isAlreadyNeed {
                onMarkAsNeed()
              } else if translation < -swipeThreshold, !isAlreadyHave {
                onMarkAsHave()
              }

              withAnimation(.spring(response: 0.35, dampingFraction: 0.7)) {
                dragOffset = 0
              }
            }
        )
      }
    }
    .frame(height: 64)
    .clipped()
  }

  private var swipeReviewRowContent: some View {
    HStack {
      VStack(alignment: .leading, spacing: 2) {
        Text(item.ingredient.name)
          .font(.body)
          .fontWeight(.medium)

        if item.hasQuantityNeeded {
          Text(formatQuantityNeeded(item.quantityNeeded))
            .font(.caption)
            .foregroundColor(.secondary)
        }
      }

      Spacer()

      if let label = statusLabel {
        Text(label)
          .font(.caption)
          .fontWeight(.medium)
          .foregroundColor(statusColor)
      } else {
        Text("—")
          .font(.caption)
          .foregroundColor(.secondary)
      }
    }
    .padding(.vertical, 12)
  }

  private var statusLabel: String? {
    switch item.status {
    case .acquired, .alreadyOwned: return "Have"
    case .needs: return "Need"
    case .unavailable: return "Unavailable"
    case .unknown: return nil
    default: return nil
    }
  }

  private var statusColor: Color {
    switch item.status {
    case .acquired, .alreadyOwned: return .green
    case .needs: return .orange
    case .unavailable: return .red
    default: return .gray
    }
  }

  private func formatQuantityNeeded(_ quantity: Common_Float32RangeWithOptionalMax) -> String {
    let unitName = MeasurementUnitFormatter.displayName(
      for: quantity.min, unit: item.measurementUnit)
    let unit = unitName.isEmpty ? "" : " \(unitName)"
    if quantity.hasMax && quantity.min != quantity.max {
      return "\(formatNumber(quantity.min))–\(formatNumber(quantity.max))\(unit)"
    } else {
      return "\(formatNumber(quantity.min))\(unit)"
    }
  }

  private func formatNumber(_ value: Float) -> String {
    value.truncatingRemainder(dividingBy: 1) == 0
      ? String(format: "%.0f", value) : String(format: "%g", value)
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
            Text(
              "Have: \(formatQuantity(item.quantityPurchased, unit: MeasurementUnitFormatter.displayName(for: item.quantityPurchased, unit: item.measurementUnit)))"
            )
            .font(.caption)
            .foregroundColor(.green)
          }
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
    let unitName = MeasurementUnitFormatter.displayName(
      for: quantity.min, unit: item.measurementUnit)
    let unit = unitName.isEmpty ? "" : " \(unitName)"
    if quantity.hasMax && quantity.min != quantity.max {
      return "\(formatNumber(quantity.min)) - \(formatNumber(quantity.max))\(unit) needed"
    } else {
      return "\(formatNumber(quantity.min))\(unit) needed"
    }
  }

  private func formatNumber(_ value: Float) -> String {
    value.truncatingRemainder(dividingBy: 1) == 0
      ? String(format: "%.0f", value) : String(format: "%g", value)
  }

  private func formatQuantity(_ quantity: Float, unit: String) -> String {
    let unitText = unit.isEmpty ? "" : " \(unit)"
    return "\(formatNumber(quantity))\(unitText)"
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
