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
          if viewModel.isLoading {
            ProgressView("Loading...")
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
                  await viewModel.loadData()
                }
              }
              .buttonStyle(.borderedProminent)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
          } else {
            ScrollView {
              VStack(spacing: 24) {
                // Members Section
                if let account = viewModel.account, !account.members.isEmpty {
                  membersSection(viewModel: viewModel, account: account)
                }

                // Account Information Section
                if viewModel.account != nil {
                  accountInformationSection(viewModel: viewModel)
                }

                // Pending Invitations Section
                if !viewModel.invitations.isEmpty {
                  pendingInvitationsSection(viewModel: viewModel)
                }

                // Send Invitation Section
                if viewModel.isAccountAdmin {
                  sendInvitationSection(viewModel: viewModel)
                }
              }
              .padding()
            }
          }
        } else {
          ProgressView("Initializing...")
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
      }
      .navigationTitle("Account Settings")
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
    VStack(alignment: .leading, spacing: 12) {
      Text("Members")
        .font(.title2)
        .fontWeight(.bold)
        .padding(.horizontal, 4)

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

  // MARK: - Account Information Section
  private func accountInformationSection(viewModel: AccountSettingsViewModel) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Information")
        .font(.title2)
        .fontWeight(.bold)
        .padding(.horizontal, 4)

      if !viewModel.isAccountAdmin {
        Text("Only account admins can edit account information")
          .font(.subheadline)
          .foregroundColor(.secondary)
          .padding(.vertical, 8)
      }

      VStack(spacing: 16) {
        TextField(
          "Account Name",
          text: Binding(get: { viewModel.accountName }, set: { viewModel.accountName = $0 })
        )
        .textFieldStyle(.roundedBorder)
        .disabled(!viewModel.isAccountAdmin)

        TextField(
          "Contact Phone",
          text: Binding(get: { viewModel.contactPhone }, set: { viewModel.contactPhone = $0 })
        )
        .textFieldStyle(.roundedBorder)
        .keyboardType(.phonePad)
        .disabled(!viewModel.isAccountAdmin)

        TextField(
          "Address Line 1",
          text: Binding(get: { viewModel.addressLine1 }, set: { viewModel.addressLine1 = $0 })
        )
        .textFieldStyle(.roundedBorder)
        .disabled(!viewModel.isAccountAdmin)

        TextField(
          "Address Line 2",
          text: Binding(get: { viewModel.addressLine2 }, set: { viewModel.addressLine2 = $0 })
        )
        .textFieldStyle(.roundedBorder)
        .disabled(!viewModel.isAccountAdmin)

        HStack(spacing: 12) {
          TextField("City", text: Binding(get: { viewModel.city }, set: { viewModel.city = $0 }))
            .textFieldStyle(.roundedBorder)
            .disabled(!viewModel.isAccountAdmin)

          TextField("State", text: Binding(get: { viewModel.state }, set: { viewModel.state = $0 }))
            .textFieldStyle(.roundedBorder)
            .disabled(!viewModel.isAccountAdmin)

          TextField(
            "Zip Code",
            text: Binding(get: { viewModel.zipCode }, set: { viewModel.zipCode = $0 })
          )
          .textFieldStyle(.roundedBorder)
          .keyboardType(.numberPad)
          .disabled(!viewModel.isAccountAdmin)
        }

        TextField(
          "Country", text: Binding(get: { viewModel.country }, set: { viewModel.country = $0 })
        )
        .textFieldStyle(.roundedBorder)
        .disabled(!viewModel.isAccountAdmin)

        if viewModel.isAccountAdmin && viewModel.accountDataHasChanged {
          Button("Update Account") {
            Task {
              await viewModel.updateAccount()
            }
          }
          .buttonStyle(.borderedProminent)
          .frame(maxWidth: .infinity)
        }
      }
    }
  }

  // MARK: - Pending Invitations Section
  private func pendingInvitationsSection(viewModel: AccountSettingsViewModel) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Awaiting Invites")
        .font(.title2)
        .fontWeight(.bold)
        .padding(.horizontal, 4)

      ForEach(viewModel.invitations, id: \.id) { invitation in
        InvitationCard(invitation: invitation)
      }
    }
  }

  // MARK: - Send Invitation Section
  private func sendInvitationSection(viewModel: AccountSettingsViewModel) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Send Invite")
        .font(.title2)
        .fontWeight(.bold)
        .padding(.horizontal, 4)

      VStack(spacing: 16) {
        TextField(
          "Email Address",
          text: Binding(get: { viewModel.invitationEmail }, set: { viewModel.invitationEmail = $0 })
        )
        .textFieldStyle(.roundedBorder)
        .keyboardType(.emailAddress)
        .autocapitalization(.none)

        TextField(
          "Name (Optional)",
          text: Binding(get: { viewModel.invitationName }, set: { viewModel.invitationName = $0 })
        )
        .textFieldStyle(.roundedBorder)

        TextField(
          "Note (Optional)",
          text: Binding(get: { viewModel.invitationNote }, set: { viewModel.invitationNote = $0 }),
          axis: .vertical
        )
        .textFieldStyle(.roundedBorder)
        .lineLimit(3...6)

        Button("Send Invitation") {
          Task {
            await viewModel.sendInvitation()
          }
        }
        .buttonStyle(.borderedProminent)
        .frame(maxWidth: .infinity)
        .disabled(viewModel.invitationEmail.isEmpty)
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
    HStack(spacing: 12) {
      // Avatar placeholder
      Circle()
        .fill(Color.gray.opacity(0.3))
        .frame(width: 50, height: 50)
        .overlay {
          Text(displayName.prefix(1).uppercased())
            .font(.headline)
            .foregroundColor(.secondary)
        }

      VStack(alignment: .leading, spacing: 4) {
        Text(displayName)
          .font(.headline)

        if member.hasBelongsToUser && member.belongsToUser.id == currentUserID {
          Text("(You)")
            .font(.caption)
            .foregroundColor(.secondary)
        }
      }

      Spacer()

      if isAccountAdmin {
        Picker(
          "Role",
          selection: Binding(
            get: { selectedRole },
            set: { newValue in
              // Don't update selectedRole yet - show the confirmation sheet first
              let newRole = newValue == "Admin" ? "account_admin" : "account_member"
              let currentRole = member.accountRole
              if newRole != currentRole {
                // Store what they selected
                pendingNewRole = newRole
                reasonText = ""
                showReasonAlert = true
              } else {
                // If they selected the same role, just update it (no change needed)
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
        Text(member.accountRole == "account_admin" ? "Admin" : "Member")
          .font(.subheadline)
          .foregroundColor(.secondary)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
    .onChange(of: member.accountRole) { _, newRole in
      // When member role changes (e.g., after successful update), sync the UI
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
              .font(.subheadline)
              .foregroundColor(.secondary)
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
              // Don't update selectedRole - it's already at the original value
              pendingNewRole = ""
              reasonText = ""
              showReasonAlert = false
            }
          }
          ToolbarItem(placement: .confirmationAction) {
            Button("Confirm") {
              if !reasonText.isEmpty {
                // Update selectedRole to reflect the change
                selectedRole = pendingNewRole == "account_admin" ? "Admin" : "Member"
                onRoleChange(pendingNewRole, reasonText)
                // originalRole will be updated via onChange when member.accountRole changes
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
      // When the sheet is dismissed (isPresented becomes false) without confirmation
      if !isPresented && !pendingNewRole.isEmpty {
        // Sheet was dismissed without confirmation - clear the pending change
        // selectedRole is already at the correct value (wasn't updated)
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

  var body: some View {
    HStack {
      VStack(alignment: .leading, spacing: 4) {
        Text(invitation.toEmail)
          .font(.headline)

        if !invitation.toName.isEmpty {
          Text(invitation.toName)
            .font(.subheadline)
            .foregroundColor(.secondary)
        }

        Text("Status: \(invitation.status)")
          .font(.caption)
          .foregroundColor(.secondary)
      }

      Spacer()
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
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
