//
//  AccountSettingsViewModel.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import SwiftProtobuf
import SwiftUI

@Observable
@MainActor
class AccountSettingsViewModel {
  // Data
  var account: Identity_Account?
  var user: Auth_GetSelfResponse?
  var invitations: [Identity_AccountInvitation] = []

  // Loading states
  var isLoading = false
  var errorMessage: String?

  // Form state
  var accountName: String = ""
  var contactPhone: String = ""
  var addressLine1: String = ""
  var addressLine2: String = ""
  var city: String = ""
  var state: String = ""
  var zipCode: String = ""
  var country: String = "USA"

  // Invitation form state
  var invitationEmail: String = ""
  var invitationName: String = ""
  var invitationNote: String = ""

  // Computed properties
  var isAccountAdmin: Bool {
    guard let account = account,
      let user = user,
      user.hasResult,
      !user.result.id.isEmpty
    else {
      return false
    }
    let userID = user.result.id
    return account.members.first { membership in
      membership.hasBelongsToUser && membership.belongsToUser.id == userID
    }?.accountRole == "account_admin"
  }

  var currentUserMembership: Identity_AccountUserMembershipWithUser? {
    guard let account = account,
      let user = user,
      user.hasResult,
      !user.result.id.isEmpty
    else {
      return nil
    }
    let userID = user.result.id
    return account.members.first { membership in
      membership.hasBelongsToUser && membership.belongsToUser.id == userID
    }
  }

  var currentUserID: String {
    guard let user = user,
      user.hasResult,
      !user.result.id.isEmpty
    else {
      return ""
    }
    return user.result.id
  }

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  func loadData() async {
    isLoading = true
    errorMessage = nil

    do {
      // Fetch account, user, and invitations in parallel
      async let accountTask = fetchActiveAccount()
      async let userTask = fetchUser()
      async let invitationsTask = fetchInvitations()

      let (accountResult, userResult, invitationsResult) = try await (
        accountTask, userTask, invitationsTask
      )

      self.account = accountResult
      self.user = userResult
      self.invitations = invitationsResult

      // Initialize form fields from account
      self.accountName = accountResult.name
      self.contactPhone = accountResult.contactPhone
      self.addressLine1 = accountResult.addressLine1
      self.addressLine2 = accountResult.addressLine2
      self.city = accountResult.city
      self.state = accountResult.state
      self.zipCode = accountResult.zipCode
      self.country = accountResult.country.isEmpty ? "USA" : accountResult.country
    } catch {
      errorMessage = "Failed to load data: \(error.localizedDescription)"
      print("❌ Error loading account settings: \(error)")
    }

    isLoading = false
  }

  private func fetchActiveAccount() async throws -> Identity_Account {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    // Get active account from auth service
    let authResponse = try await clientManager.client.auth.getActiveAccount(
      Auth_GetActiveAccountRequest(),
      metadata: metadata,
      options: clientManager.defaultCallOptions
    )

    guard authResponse.hasResult, !authResponse.result.id.isEmpty else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 3,
        userInfo: [NSLocalizedDescriptionKey: "No active account found"])
    }

    let accountID = authResponse.result.id

    // Get full account details from identity service
    var request = Identity_GetAccountRequest()
    request.accountID = accountID

    let identityResponse = try await clientManager.client.identity.getAccount(
      request,
      metadata: metadata,
      options: clientManager.defaultCallOptions
    )

    guard identityResponse.hasResult else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 4,
        userInfo: [NSLocalizedDescriptionKey: "Account not found"])
    }

    return identityResponse.result
  }

  private func fetchUser() async throws -> Auth_GetSelfResponse {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    let response = try await clientManager.client.auth.getSelf(
      Auth_GetSelfRequest(),
      metadata: metadata,
      options: clientManager.defaultCallOptions
    )

    return response
  }

  private func fetchInvitations() async throws -> [Identity_AccountInvitation] {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "AccountSettingsViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    // Create an empty filter (no filtering)
    let filter = Filtering_QueryFilter()
    var request = Identity_GetSentAccountInvitationsRequest()
    request.filter = filter

    let response = try await clientManager.client.identity.getSentAccountInvitations(
      request,
      metadata: metadata,
      options: clientManager.defaultCallOptions
    )

    return response.results
  }

  func updateAccount() async -> Bool {
    guard let account = account else {
      errorMessage = "No account loaded"
      return false
    }

    guard isAccountAdmin else {
      errorMessage = "Only account admins can update account information"
      return false
    }

    isLoading = true
    errorMessage = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "AccountSettingsViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "AccountSettingsViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      var updateInput = Identity_AccountUpdateRequestInput()
      updateInput.name = accountName
      updateInput.contactPhone = contactPhone
      updateInput.addressLine1 = addressLine1
      updateInput.addressLine2 = addressLine2
      updateInput.city = city
      updateInput.state = state
      updateInput.zipCode = zipCode
      updateInput.country = country

      var request = Identity_UpdateAccountRequest()
      request.accountID = account.id
      request.input = updateInput

      _ = try await clientManager.client.identity.updateAccount(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      // Reload data to get updated account
      await loadData()

      isLoading = false
      return true
    } catch {
      errorMessage = "Failed to update account: \(error.localizedDescription)"
      print("❌ Error updating account: \(error)")
      isLoading = false
      return false
    }
  }

  func sendInvitation() async -> Bool {
    guard account != nil else {
      errorMessage = "No account loaded"
      return false
    }

    guard isAccountAdmin else {
      errorMessage = "Only account admins can send invitations"
      return false
    }

    guard !invitationEmail.isEmpty else {
      errorMessage = "Email address is required"
      return false
    }

    // Basic email validation
    let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
    let emailPredicate = NSPredicate(format: "SELF MATCHES %@", emailRegex)
    guard emailPredicate.evaluate(with: invitationEmail) else {
      errorMessage = "Invalid email address"
      return false
    }

    isLoading = true
    errorMessage = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "AccountSettingsViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "AccountSettingsViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      var invitationInput = Identity_AccountInvitationCreationRequestInput()
      invitationInput.toEmail = invitationEmail
      invitationInput.toName = invitationName
      invitationInput.note = invitationNote

      var request = Identity_CreateAccountInvitationRequest()
      request.input = invitationInput

      _ = try await clientManager.client.identity.createAccountInvitation(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      // Clear form and reload invitations
      invitationEmail = ""
      invitationName = ""
      invitationNote = ""
      await loadData()

      isLoading = false
      return true
    } catch {
      errorMessage = "Failed to send invitation: \(error.localizedDescription)"
      print("❌ Error sending invitation: \(error)")
      isLoading = false
      return false
    }
  }

  func updateMemberRole(membershipID: String, newRole: String, reason: String) async -> Bool {
    guard isAccountAdmin else {
      errorMessage = "Only account admins can change member roles"
      return false
    }

    guard !reason.isEmpty else {
      errorMessage = "A reason is required for changing member roles"
      return false
    }

    guard let membership = account?.members.first(where: { $0.id == membershipID }) else {
      errorMessage = "Member not found"
      return false
    }

    guard membership.hasBelongsToUser, !membership.belongsToUser.id.isEmpty else {
      errorMessage = "User ID not found"
      return false
    }
    let userID = membership.belongsToUser.id

    isLoading = true
    errorMessage = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "AccountSettingsViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "AccountSettingsViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      var updateInput = Identity_ModifyUserPermissionsInput()
      updateInput.newRole = newRole
      updateInput.reason = reason

      var request = Identity_UpdateAccountMemberPermissionsRequest()
      request.userID = userID
      request.input = updateInput

      _ = try await clientManager.client.identity.updateAccountMemberPermissions(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      // Reload data to get updated account
      await loadData()

      isLoading = false
      return true
    } catch {
      errorMessage = "Failed to update member role: \(error.localizedDescription)"
      print("❌ Error updating member role: \(error)")
      isLoading = false
      return false
    }
  }

  var accountDataHasChanged: Bool {
    guard let account = account else { return false }
    return account.name != accountName
      || account.contactPhone != contactPhone
      || account.addressLine1 != addressLine1
      || account.addressLine2 != addressLine2
      || account.city != city
      || account.state != state
      || account.zipCode != zipCode
      || account.country != country
  }
}

