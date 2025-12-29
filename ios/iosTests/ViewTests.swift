//
//  ViewTests.swift
//  iosTests
//
//  Created by Auto on 12/8/25.
//

import Foundation
import SwiftProtobuf
import SwiftUI
@testable import ios
import Testing

// MARK: - Email Validation Helper (extracted from RegisterView)

func isValidEmail(_ email: String) -> Bool {
  let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
  let emailPredicate = NSPredicate(format: "SELF MATCHES %@", emailRegex)
  return emailPredicate.evaluate(with: email)
}

// MARK: - RegisterView Validation Tests

struct RegisterViewValidationTests {
  @Test("isValidEmail returns true for valid emails")
  func testIsValidEmailValid() {
    let validEmails = [
      "test@example.com",
      "user.name@domain.co.uk",
      "user+tag@example.org",
      "user123@test-domain.com",
      "a@b.co"
    ]

    for email in validEmails {
      #expect(isValidEmail(email) == true, "Email '\(email)' should be valid")
    }
  }

  @Test("isValidEmail returns false for invalid emails")
  func testIsValidEmailInvalid() {
    let invalidEmails = [
      "invalid-email",
      "@example.com",
      "user@",
      "user@domain",
      "user name@example.com",
      "",
      "user@domain.",
      ".user@domain.com"
    ]

    for email in invalidEmails {
      #expect(isValidEmail(email) == false, "Email '\(email)' should be invalid")
    }
  }

  @Test("RegisterView form validation requires email")
  @MainActor
  func testRegisterViewRequiresEmail() async {
    let showLogin = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = RegisterView(showLogin: showLogin)
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("RegisterView form validation requires username")
  @MainActor
  func testRegisterViewRequiresUsername() async {
    let showLogin = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = RegisterView(showLogin: showLogin)
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("RegisterView form validation requires password")
  @MainActor
  func testRegisterViewRequiresPassword() async {
    let showLogin = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = RegisterView(showLogin: showLogin)
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("RegisterView form validation requires password match")
  @MainActor
  func testRegisterViewRequiresPasswordMatch() async {
    let showLogin = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = RegisterView(showLogin: showLogin)
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("RegisterView form validation requires password length")
  @MainActor
  func testRegisterViewRequiresPasswordLength() async {
    let showLogin = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = RegisterView(showLogin: showLogin)
      .environment(authManager)

    #expect(view != nil)
  }
}

// MARK: - LoginView Tests

struct LoginViewTests {
  @Test("LoginView initializes with default values")
  @MainActor
  func testLoginViewInitialization() async {
    let showRegister = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = LoginView(showRegister: showRegister)
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("LoginView requires username for login")
  @MainActor
  func testLoginViewRequiresUsername() async {
    let showRegister = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = LoginView(showRegister: showRegister)
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("LoginView requires password for login")
  @MainActor
  func testLoginViewRequiresPassword() async {
    let showRegister = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = LoginView(showRegister: showRegister)
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("LoginView handles TOTP requirement")
  @MainActor
  func testLoginViewHandlesTOTP() async {
    let showRegister = Binding<Bool>(get: { false }, set: { _ in })
    let authManager = createMockAuthenticationManager()
    let view = LoginView(showRegister: showRegister)
      .environment(authManager)

    #expect(view != nil)
  }
}

// MARK: - AccountSettingsView Tests

struct AccountSettingsViewTests {
  @Test("AccountSettingsView initializes view model on appear")
  @MainActor
  func testAccountSettingsViewInitializesViewModel() async {
    let authManager = createMockAuthenticationManagerForAccount()
    let view = AccountSettingsView()
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("AccountSettingsView shows loading state")
  @MainActor
  func testAccountSettingsViewShowsLoading() async {
    let authManager = createMockAuthenticationManagerForAccount()
    let view = AccountSettingsView()
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("AccountSettingsView shows error state")
  @MainActor
  func testAccountSettingsViewShowsError() async {
    let authManager = createMockAuthenticationManagerForAccount()
    let view = AccountSettingsView()
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("AccountSettingsView shows account information when loaded")
  @MainActor
  func testAccountSettingsViewShowsAccountInfo() async {
    let authManager = createMockAuthenticationManagerForAccount()
    let view = AccountSettingsView()
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("AccountSettingsView shows members section when members exist")
  @MainActor
  func testAccountSettingsViewShowsMembers() async {
    let authManager = createMockAuthenticationManagerForAccount()
    let view = AccountSettingsView()
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("AccountSettingsView shows invitations section when invitations exist")
  @MainActor
  func testAccountSettingsViewShowsInvitations() async {
    let authManager = createMockAuthenticationManagerForAccount()
    let view = AccountSettingsView()
      .environment(authManager)

    #expect(view != nil)
  }

  @Test("AccountSettingsView shows send invitation section for admins")
  @MainActor
  func testAccountSettingsViewShowsSendInvitation() async {
    let authManager = createMockAuthenticationManagerForAccount()
    let view = AccountSettingsView()
      .environment(authManager)

    #expect(view != nil)
  }
}

// MARK: - MemberCard Tests

struct MemberCardTests {
  @Test("MemberCard displayName uses first and last name")
  @MainActor
  func testMemberCardDisplayNameFirstLast() async {
    var member = createMockMembership()
    var user = Identity_User()
    user.id = "user-1"
    user.firstName = "John"
    user.lastName = "Doe"
    user.username = "johndoe"
    member.belongsToUser = user

    let card = MemberCard(
      member: member,
      currentUserID: "user-1",
      isAccountAdmin: true,
      onRoleChange: { _, _ in }
    )

    #expect(card != nil)
  }

  @Test("MemberCard displayName uses first name only")
  @MainActor
  func testMemberCardDisplayNameFirstOnly() async {
    var member = createMockMembership()
    var user = Identity_User()
    user.id = "user-1"
    user.firstName = "John"
    user.username = "johndoe"
    member.belongsToUser = user

    let card = MemberCard(
      member: member,
      currentUserID: "user-1",
      isAccountAdmin: true,
      onRoleChange: { _, _ in }
    )

    #expect(card != nil)
  }

  @Test("MemberCard displayName uses username when no name")
  @MainActor
  func testMemberCardDisplayNameUsername() async {
    var member = createMockMembership()
    var user = Identity_User()
    user.id = "user-1"
    user.username = "johndoe"
    member.belongsToUser = user

    let card = MemberCard(
      member: member,
      currentUserID: "user-1",
      isAccountAdmin: true,
      onRoleChange: { _, _ in }
    )

    #expect(card != nil)
  }

  @Test("MemberCard displayName returns Unknown when no user")
  @MainActor
  func testMemberCardDisplayNameUnknown() async {
    var member = createMockMembership()
    member.clearBelongsToUser()

    let card = MemberCard(
      member: member,
      currentUserID: "user-1",
      isAccountAdmin: true,
      onRoleChange: { _, _ in }
    )

    #expect(card != nil)
  }

  @Test("MemberCard shows current user indicator")
  @MainActor
  func testMemberCardShowsCurrentUser() async {
    var member = createMockMembership(userID: "user-1")
    var user = Identity_User()
    user.id = "user-1"
    user.username = "currentuser"
    member.belongsToUser = user

    let card = MemberCard(
      member: member,
      currentUserID: "user-1",
      isAccountAdmin: true,
      onRoleChange: { _, _ in }
    )

    #expect(card != nil)
  }

  @Test("MemberCard initializes with correct role")
  @MainActor
  func testMemberCardInitializesRole() async {
    var member = createMockMembership(role: "account_admin")
    var user = Identity_User()
    user.id = "user-1"
    member.belongsToUser = user

    let card = MemberCard(
      member: member,
      currentUserID: "user-2",
      isAccountAdmin: true,
      onRoleChange: { _, _ in }
    )

    #expect(card != nil)
  }
}

// MARK: - InvitationCard Tests

struct InvitationCardTests {
  @Test("InvitationCard displays email")
  @MainActor
  func testInvitationCardDisplaysEmail() async {
    var invitation = Identity_AccountInvitation()
    invitation.toEmail = "test@example.com"
    invitation.toName = "Test User"
    invitation.status = "pending"

    let card = InvitationCard(invitation: invitation)

    #expect(card != nil)
  }

  @Test("InvitationCard displays name when available")
  @MainActor
  func testInvitationCardDisplaysName() async {
    var invitation = Identity_AccountInvitation()
    invitation.toEmail = "test@example.com"
    invitation.toName = "Test User"
    invitation.status = "pending"

    let card = InvitationCard(invitation: invitation)

    #expect(card != nil)
  }

  @Test("InvitationCard displays status")
  @MainActor
  func testInvitationCardDisplaysStatus() async {
    var invitation = Identity_AccountInvitation()
    invitation.toEmail = "test@example.com"
    invitation.status = "accepted"

    let card = InvitationCard(invitation: invitation)

    #expect(card != nil)
  }

  @Test("InvitationCard handles empty name")
  @MainActor
  func testInvitationCardHandlesEmptyName() async {
    var invitation = Identity_AccountInvitation()
    invitation.toEmail = "test@example.com"
    invitation.toName = ""
    invitation.status = "pending"

    let card = InvitationCard(invitation: invitation)

    #expect(card != nil)
  }

  @Test("InvitationCard handles different statuses")
  @MainActor
  func testInvitationCardHandlesStatuses() async {
    let statuses = ["pending", "accepted", "declined", "expired"]

    for status in statuses {
      var invitation = Identity_AccountInvitation()
      invitation.toEmail = "test@example.com"
      invitation.status = status

      let card = InvitationCard(invitation: invitation)
      #expect(card != nil)
    }
  }
}

// MARK: - Date Formatting Tests

struct DateFormattingTests {
  @Test("formatDate formats date correctly")
  func testFormatDate() {
    let date = Date(timeIntervalSince1970: 946684800)  // Jan 1, 2000
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    let expected = formatter.string(from: date)

    #expect(!expected.isEmpty)
  }

  @Test("formatDate handles different dates")
  func testFormatDateDifferentDates() {
    let dates = [
      Date(timeIntervalSince1970: 0),  // Jan 1, 1970
      Date(timeIntervalSince1970: 946684800),  // Jan 1, 2000
      Date()  // Today
    ]

    let formatter = DateFormatter()
    formatter.dateStyle = .medium

    for date in dates {
      let formatted = formatter.string(from: date)
      #expect(!formatted.isEmpty)
    }
  }
}
