//
//  HouseholdInstrumentsView.swift
//  ios
//

import SwiftUI

struct HouseholdInstrumentsView: View {
  let viewModel: AccountSettingsViewModel
  @State private var instrumentSearchQuery: String = ""

  var body: some View {
    DSContentState(
      isLoading: viewModel.isLoading,
      loadingMessage: "Loading household...",
      error: viewModel.errorMessage,
      errorTitle: viewModel.errorTitle,
      errorIcon: viewModel.errorIcon,
      errorIconColor: viewModel.errorIconColor,
      onRetry: { await viewModel.loadData() },
      content: { instrumentsContent }
    )
    .navigationTitle("Kitchen Instruments")
    .refreshable {
      await viewModel.loadData()
    }
  }

  @ViewBuilder
  private var instrumentsContent: some View {
    if viewModel.account != nil {
      ScrollView {
        VStack(spacing: DSTheme.Spacing.xl) {
          addInstrumentSection
          instrumentsSection
        }
        .dsScreenPadding()
      }
    }
  }

  private var instrumentsSection: some View {
    DSSection("Kitchen Instruments") {
      if viewModel.instrumentOwnerships.isEmpty {
        DSSectionEmptyContent(
          "No kitchen instruments added yet. Add the tools and appliances your household has.",
          icon: "frying.pan"
        )
      } else {
        ForEach(viewModel.instrumentOwnerships, id: \.id) { ownership in
          InstrumentOwnershipCard(
            ownership: ownership,
            onEdit: { quantity, notes in
              Task {
                await viewModel.updateInstrumentOwnership(
                  ownershipID: ownership.id,
                  quantity: quantity,
                  notes: notes
                )
              }
            },
            onRemove: {
              Task {
                await viewModel.archiveInstrumentOwnership(ownershipID: ownership.id)
              }
            }
          )
        }
      }
    }
  }

  private var filteredInstruments: [Mealplanning_ValidInstrument] {
    let query = instrumentSearchQuery.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
    guard !query.isEmpty else {
      return viewModel.validInstruments
    }
    return viewModel.validInstruments.filter { instrument in
      instrument.name.lowercased().contains(query)
        || instrument.pluralName.lowercased().contains(query)
        || (!instrument.description_p.isEmpty && instrument.description_p.lowercased().contains(query))
    }
  }

  private var selectedInstrument: Mealplanning_ValidInstrument? {
    guard !viewModel.newInstrumentValidInstrumentID.isEmpty else { return nil }
    return viewModel.validInstruments.first { $0.id == viewModel.newInstrumentValidInstrumentID }
  }

  private var addInstrumentSection: some View {
    DSSection(
      "Add Instrument",
      subtitle: "Add a tool or appliance your household owns"
    ) {
      VStack(spacing: DSTheme.Spacing.lg) {
        DSTextField(
          "Search instruments...",
          text: $instrumentSearchQuery
        )
        .autocorrectionDisabled()

        if let selected = selectedInstrument {
          HStack {
            Text("Selected: \(selected.name)")
              .font(DSTheme.Typography.body)
              .foregroundColor(DSTheme.Colors.textPrimary)
            Spacer()
            DSButton("Clear", style: .ghost, size: .small) {
              viewModel.newInstrumentValidInstrumentID = ""
            }
          }
          .padding(.vertical, DSTheme.Spacing.xs)
        }

        if filteredInstruments.isEmpty {
          Text("No instruments found")
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .frame(maxWidth: .infinity, alignment: .leading)
        } else {
          ScrollView {
            LazyVStack(spacing: DSTheme.Spacing.xs) {
              ForEach(filteredInstruments, id: \.id) { instrument in
                Button {
                  viewModel.newInstrumentValidInstrumentID = instrument.id
                } label: {
                  HStack {
                    Text(instrument.name)
                      .font(DSTheme.Typography.body)
                      .foregroundColor(DSTheme.Colors.textPrimary)
                    if viewModel.newInstrumentValidInstrumentID == instrument.id {
                      Image(systemName: "checkmark.circle.fill")
                        .foregroundColor(DSTheme.Colors.primary)
                    }
                    Spacer()
                  }
                  .padding(.vertical, DSTheme.Spacing.sm)
                  .padding(.horizontal, DSTheme.Spacing.md)
                  .background(
                    viewModel.newInstrumentValidInstrumentID == instrument.id
                      ? DSTheme.Colors.primary.opacity(0.1)
                      : Color.clear
                  )
                  .contentShape(Rectangle())
                }
                .buttonStyle(.plain)
              }
            }
          }
          .frame(maxHeight: 200)
        }

        HStack {
          Text("Quantity")
            .font(DSTheme.Typography.body)
            .foregroundColor(DSTheme.Colors.textPrimary)
          Spacer()
          Stepper(
            value: Binding(
              get: { Int(viewModel.newInstrumentQuantity) },
              set: { viewModel.newInstrumentQuantity = UInt32(max(1, min(999, $0))) }
            ),
            in: 1...999
          ) {
            Text("\(viewModel.newInstrumentQuantity)")
              .font(DSTheme.Typography.label)
              .foregroundColor(DSTheme.Colors.textPrimary)
              .frame(minWidth: 32, alignment: .trailing)
          }
        }

        DSTextField(
          "Notes (Optional)",
          text: Binding(
            get: { viewModel.newInstrumentNotes },
            set: { viewModel.newInstrumentNotes = $0 }
          ),
          type: .multiline
        )

        DSButton(
          "Add Instrument",
          icon: "plus",
          fullWidth: true,
          isDisabled: viewModel.newInstrumentValidInstrumentID.isEmpty
        ) {
          Task {
            await viewModel.createInstrumentOwnership()
          }
        }
      }
    }
  }
}

// MARK: - Instrument Ownership Card

struct InstrumentOwnershipCard: View {
  let ownership: Mealplanning_AccountInstrumentOwnership
  let onEdit: (UInt32?, String?) -> Void
  let onRemove: () -> Void

  @State private var showEditSheet = false
  @State private var showRemoveConfirmation = false
  @State private var editQuantity: UInt32 = 1
  @State private var editNotes: String = ""

  var body: some View {
    DSCard {
      HStack(spacing: DSTheme.Spacing.md) {
        Image(systemName: "frying.pan")
          .font(.title2)
          .foregroundColor(DSTheme.Colors.textSecondary)

        VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
          Text(instrumentName)
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          HStack(spacing: DSTheme.Spacing.sm) {
            Text(quantityLabel)
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)

            if !ownership.notes.isEmpty {
              Text("•")
                .font(DSTheme.Typography.caption)
                .foregroundColor(DSTheme.Colors.textSecondary)
              Text(ownership.notes)
                .font(DSTheme.Typography.caption)
                .foregroundColor(DSTheme.Colors.textSecondary)
                .lineLimit(1)
            }
          }
        }

        Spacer()

        HStack(spacing: DSTheme.Spacing.sm) {
          DSButton("Edit", style: .ghost, size: .small) {
            editQuantity = ownership.quantity
            editNotes = ownership.notes
            showEditSheet = true
          }
          DSButton("Remove", style: .ghost, size: .small) {
            showRemoveConfirmation = true
          }
        }
      }
    }
    .sheet(isPresented: $showEditSheet) {
      NavigationStack {
        Form {
          Section {
            HStack {
              Text("Quantity")
              Spacer()
              Stepper(
                value: Binding(
                  get: { Int(editQuantity) },
                  set: { editQuantity = UInt32(max(1, min(999, $0))) }
                ),
                in: 1...999
              ) {
                Text("\(editQuantity)")
                  .frame(minWidth: 32, alignment: .trailing)
              }
            }
            TextField("Notes (Optional)", text: $editNotes, axis: .vertical)
              .lineLimit(2...4)
          }
        }
        .navigationTitle("Edit Instrument")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
          ToolbarItem(placement: .cancellationAction) {
            Button("Cancel") {
              showEditSheet = false
            }
          }
          ToolbarItem(placement: .confirmationAction) {
            Button("Save") {
              onEdit(editQuantity, editNotes)
              showEditSheet = false
            }
          }
        }
      }
    }
    .alert("Remove Instrument", isPresented: $showRemoveConfirmation) {
      Button("Cancel", role: .cancel) {
        showRemoveConfirmation = false
      }
      Button("Remove", role: .destructive) {
        onRemove()
        showRemoveConfirmation = false
      }
    } message: {
      Text("Remove this instrument from your household?")
    }
  }

  private var instrumentName: String {
    ownership.hasInstrument ? ownership.instrument.name : "Unknown"
  }

  private var quantityLabel: String {
    let qty = ownership.quantity
    if ownership.hasInstrument, !ownership.instrument.pluralName.isEmpty, qty > 1 {
      return "\(qty) \(ownership.instrument.pluralName)"
    }
    return "\(qty) \(instrumentName)"
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return NavigationStack {
    HouseholdInstrumentsView(viewModel: AccountSettingsViewModel(authManager: authManager))
      .environment(authManager)
  }
}
