//
//  AccountSettingsView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftUI

struct AccountSettingsView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: AccountSettingsViewModel?

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          DSContentState(
            isLoading: viewModel.isLoading,
            loadingMessage: "Loading household...",
            error: viewModel.errorMessage,
            onRetry: { await viewModel.loadData() },
            content: {
              ScrollView {
                VStack(spacing: DSTheme.Spacing.xl) {
                  // Members Section
                  if let account = viewModel.account {
                    membersSection(viewModel: viewModel, account: account)
                  }

                  // Send Invitation Section (admins only)
                  if viewModel.isAccountAdmin {
                    sendInvitationSection(viewModel: viewModel)
                  }

                  // Invitations (sent) and their status
                  if !viewModel.invitations.isEmpty {
                    pendingInvitationsSection(viewModel: viewModel)
                  }

                  // Household details (address, contact)
                  if viewModel.account != nil {
                    accountInformationSection(viewModel: viewModel)
                  }
                }
                .dsScreenPadding()
              }
            })
        } else {
          DSInitializingView()
        }
      }
      .navigationTitle("My Household")
      .refreshable {
        if let viewModel = viewModel {
          await viewModel.loadData()
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = AccountSettingsViewModel(authManager: authManager)
        }
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
    }
  }

  // MARK: - Members Section
  private func membersSection(viewModel: AccountSettingsViewModel, account: Identity_Account)
    -> some View
  {
    DSSection("Household Members") {
      if account.members.isEmpty {
        DSSectionEmptyContent(
          "No members yet. Invite someone to join your household.",
          icon: "person.2"
        )
      } else {
        ForEach(account.members, id: \.id) { member in
          MemberCard(
            member: member,
            currentUserID: viewModel.currentUserID,
            isAccountAdmin: viewModel.isAccountAdmin,
            onRoleChange: { newRole, reason in
              Task {
                await viewModel.updateMemberRole(
                  membershipID: member.id, newRole: newRole, reason: reason)
              }
            }
          )
        }
      }
    }
  }

  // MARK: - Account Information Section
  private func accountInformationSection(viewModel: AccountSettingsViewModel) -> some View {
    DSSection(
      "Household Details",
      subtitle: viewModel.isAccountAdmin ? nil : "Only household admins can edit household details"
    ) {
      VStack(spacing: DSTheme.Spacing.lg) {
        DSTextField(
          "Household Name",
          text: Binding(get: { viewModel.accountName }, set: { viewModel.accountName = $0 }),
          isDisabled: !viewModel.isAccountAdmin
        )

        // DSTextField(
        //   "Contact Phone",
        //   text: Binding(get: { viewModel.contactPhone }, set: { viewModel.contactPhone = $0 }),
        //   type: .phone,
        //   isDisabled: !viewModel.isAccountAdmin
        // )
        //
        //        DSTextField(
        //          "Address Line 1",
        //          text: Binding(get: { viewModel.addressLine1 }, set: { viewModel.addressLine1 = $0 }),
        //          isDisabled: !viewModel.isAccountAdmin
        //        )
        //
        //        DSTextField(
        //          "Address Line 2",
        //          text: Binding(get: { viewModel.addressLine2 }, set: { viewModel.addressLine2 = $0 }),
        //          isDisabled: !viewModel.isAccountAdmin
        //        )
        //
        //        HStack(spacing: DSTheme.Spacing.md) {
        //          DSTextField(
        //            "City",
        //            text: Binding(get: { viewModel.city }, set: { viewModel.city = $0 }),
        //            isDisabled: !viewModel.isAccountAdmin
        //          )
        //
        //          DSTextField(
        //            "State",
        //            text: Binding(get: { viewModel.state }, set: { viewModel.state = $0 }),
        //            isDisabled: !viewModel.isAccountAdmin
        //          )
        //
        //          DSTextField(
        //            "Zip Code",
        //            text: Binding(get: { viewModel.zipCode }, set: { viewModel.zipCode = $0 }),
        //            type: .number,
        //            isDisabled: !viewModel.isAccountAdmin
        //          )
        //        }
        //
        //        DSTextField(
        //          "Country",
        //          text: Binding(get: { viewModel.country }, set: { viewModel.country = $0 }),
        //          isDisabled: !viewModel.isAccountAdmin
        //        )

        if viewModel.isAccountAdmin && viewModel.accountDataHasChanged {
          DSButton("Update Household", icon: "checkmark", fullWidth: true) {
            Task {
              await viewModel.updateAccount()
            }
          }
        }
      }
    }
  }

  // MARK: - Pending Invitations Section
  private func pendingInvitationsSection(viewModel: AccountSettingsViewModel) -> some View {
    DSSection(
      "Invitations",
      subtitle: "Invitations you've sent for this household and their status."
    ) {
      ForEach(viewModel.invitations, id: \.id) { invitation in
        InvitationCard(
          invitation: invitation,
          isAccountAdmin: viewModel.isAccountAdmin,
          onCancel: {
            _ = await viewModel.cancelInvitation(invitationID: invitation.id)
          }
        )
      }
    }
  }

  // MARK: - Send Invitation Section
  private func sendInvitationSection(viewModel: AccountSettingsViewModel) -> some View {
    DSSection(
      "Add Someone to Your Household",
      subtitle: "Send an invitation by email. They can join once they have an account."
    ) {
      VStack(spacing: DSTheme.Spacing.lg) {
        DSTextField(
          "Email Address",
          text: Binding(
            get: { viewModel.invitationEmail }, set: { viewModel.invitationEmail = $0 }),
          type: .email
        )

        DSTextField(
          "Name (Optional)",
          text: Binding(get: { viewModel.invitationName }, set: { viewModel.invitationName = $0 })
        )

        DSTextField(
          "Note (Optional)",
          text: Binding(get: { viewModel.invitationNote }, set: { viewModel.invitationNote = $0 }),
          type: .multiline
        )

        DSButton(
          "Send Invitation",
          icon: "envelope",
          fullWidth: true,
          isDisabled: viewModel.invitationEmail.isEmpty
        ) {
          Task {
            await viewModel.sendInvitation()
          }
        }
      }
    }
  }
}

// MARK: - Member Card
struct MemberCard: View {
  let member: Identity_AccountUserMembershipWithUser
  let currentUserID: String
  let isAccountAdmin: Bool
  let onRoleChange: (String, String) -> Void

  @State private var selectedRole: String
  @State private var showReasonAlert = false
  @State private var pendingNewRole: String = ""
  @State private var reasonText: String = ""
  @State private var originalRole: String

  init(
    member: Identity_AccountUserMembershipWithUser,
    currentUserID: String,
    isAccountAdmin: Bool,
    onRoleChange: @escaping (String, String) -> Void
  ) {
    self.member = member
    self.currentUserID = currentUserID
    self.isAccountAdmin = isAccountAdmin
    self.onRoleChange = onRoleChange
    let initialRole = member.accountRole == "account_admin" ? "Admin" : "Member"
    _selectedRole = State(initialValue: initialRole)
    _originalRole = State(initialValue: initialRole)
  }

  var body: some View {
    DSCard {
      HStack(spacing: DSTheme.Spacing.md) {
        // Avatar
        DSAvatarView(name: displayName, size: .lg)

        VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
          Text(displayName)
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          if member.hasBelongsToUser && member.belongsToUser.id == currentUserID {
            Text("(You)")
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
        }

        Spacer()

        if isAccountAdmin {
          Picker(
            "Role",
            selection: Binding(
              get: { selectedRole },
              set: { newValue in
                let newRole = newValue == "Admin" ? "account_admin" : "account_member"
                let currentRole = member.accountRole
                if newRole != currentRole {
                  pendingNewRole = newRole
                  reasonText = ""
                  showReasonAlert = true
                } else {
                  selectedRole = newValue
                }
              }
            )
          ) {
            Text("Member").tag("Member")
            Text("Admin").tag("Admin")
          }
          .pickerStyle(.menu)
        } else {
          DSStatusBadge(
            .custom(
              member.accountRole == "account_admin" ? "Admin" : "Member",
              DSTheme.Colors.textSecondary
            ),
            style: .minimal
          )
        }
      }
    }
    .onChange(of: member.accountRole) { _, newRole in
      let newRoleString = newRole == "account_admin" ? "Admin" : "Member"
      originalRole = newRoleString
      selectedRole = newRoleString
    }
    .sheet(isPresented: $showReasonAlert) {
      NavigationStack {
        Form {
          Section {
            let roleName = pendingNewRole == "account_admin" ? "Admin" : "Member"
            Text("Please provide a reason for changing this user's role to \(roleName).")
              .font(DSTheme.Typography.body)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
          Section {
            TextField("Reason", text: $reasonText, axis: .vertical)
              .lineLimit(3...6)
          }
        }
        .navigationTitle("Change Role")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
          ToolbarItem(placement: .cancellationAction) {
            Button("Cancel") {
              pendingNewRole = ""
              reasonText = ""
              showReasonAlert = false
            }
          }
          ToolbarItem(placement: .confirmationAction) {
            Button("Confirm") {
              if !reasonText.isEmpty {
                selectedRole = pendingNewRole == "account_admin" ? "Admin" : "Member"
                onRoleChange(pendingNewRole, reasonText)
                pendingNewRole = ""
                reasonText = ""
                showReasonAlert = false
              }
            }
            .disabled(reasonText.isEmpty)
          }
        }
      }
      .presentationDetents([.medium])
    }
    .onChange(of: showReasonAlert) { _, isPresented in
      if !isPresented && !pendingNewRole.isEmpty {
        pendingNewRole = ""
        reasonText = ""
      }
    }
  }

  private var displayName: String {
    guard member.hasBelongsToUser else {
      return "Unknown User"
    }
    let user = member.belongsToUser
    if !user.firstName.isEmpty {
      if !user.lastName.isEmpty {
        return "\(user.firstName) \(user.lastName)"
      }
      return user.firstName
    }
    return user.username.isEmpty ? "Unknown User" : user.username
  }
}

// MARK: - Invitation Card
struct InvitationCard: View {
  let invitation: Identity_AccountInvitation
  let isAccountAdmin: Bool
  let onCancel: (() async -> Void)?

  var body: some View {
    DSCard {
      HStack {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
          Text(invitation.toEmail)
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          if !invitation.toName.isEmpty {
            Text(invitation.toName)
              .font(DSTheme.Typography.body)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }

          DSStatusBadge(status: invitation.status)
        }

        Spacer()

        if isAccountAdmin && invitation.status.lowercased() == "pending" {
          DSButton("Cancel", style: .ghost, size: .small) {
            Task {
              await onCancel?()
            }
          }
        }
      }
    }
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return AccountSettingsView()
    .environment(authManager)
}
