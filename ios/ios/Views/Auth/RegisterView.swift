//
//  RegisterView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftUI

struct RegisterView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Binding var showLogin: Bool

  @State private var emailAddress: String = ""
  @State private var username: String = ""
  @State private var password: String = ""
  @State private var repeatedPassword: String = ""
  @State private var accountName: String = ""
  @State private var firstName: String = ""
  @State private var lastName: String = ""
  @State private var birthday: Date?
  @State private var showBirthdayPicker: Bool = false
  @State private var errorMessage: String = ""
  @State private var isLoading: Bool = false
  @State private var registrationTask: Task<Void, Never>?

  // Optional invitation fields
  @State private var invitationToken: String = ""
  @State private var invitationID: String = ""

  var body: some View {
    ScrollView {
      VStack(spacing: 20) {
        Spacer()
          .frame(height: 20)

        // App Title
        Text("Dinner Done Better")
          .font(.largeTitle)
          .fontWeight(.bold)

        Text("Create your account")
          .font(.subheadline)
          .foregroundColor(.secondary)

        // Registration Form
        VStack(spacing: 16) {
          // Required fields
          TextField("Email Address", text: $emailAddress)
            .textFieldStyle(.roundedBorder)
            .textInputAutocapitalization(.never)
            .autocorrectionDisabled()
            .keyboardType(.emailAddress)
            .disabled(isLoading)
            .accessibilityIdentifier("registrationEmailAddressInput")

          TextField("Username", text: $username)
            .textFieldStyle(.roundedBorder)
            .textInputAutocapitalization(.never)
            .autocorrectionDisabled()
            .disabled(isLoading)
            .accessibilityIdentifier("registrationUsernameInput")

          SecureField("Password", text: $password)
            .textFieldStyle(.roundedBorder)
            .disabled(isLoading)
            .accessibilityIdentifier("registrationPasswordInput")

          SecureField("Password (again)", text: $repeatedPassword)
            .textFieldStyle(.roundedBorder)
            .disabled(isLoading)
            .accessibilityIdentifier("registrationPasswordConfirmInput")

          // Divider for optional fields
          HStack {
            Rectangle()
              .fill(Color.gray.opacity(0.3))
              .frame(height: 1)
            Text("optional fields")
              .font(.caption)
              .foregroundColor(.secondary)
            Rectangle()
              .fill(Color.gray.opacity(0.3))
              .frame(height: 1)
          }
          .padding(.vertical, 8)

          // Optional fields
          TextField("Account Name", text: $accountName)
            .textFieldStyle(.roundedBorder)
            .disabled(isLoading)
            .accessibilityIdentifier("registrationAccountNameInput")

          // Birthday picker
          Button(
            action: {
              showBirthdayPicker.toggle()
            },
            label: {
              HStack {
                if let birthday = birthday {
                  Text(formatDate(birthday))
                    .foregroundColor(.primary)
                } else {
                  Text("Birthday (optional)")
                    .foregroundColor(.secondary)
                }
                Spacer()
                Image(systemName: "calendar")
                  .foregroundColor(.secondary)
              }
              .padding()
              .background(Color(uiColor: .systemGray6))
              .cornerRadius(8)
            }
          )
          .disabled(isLoading)
          .accessibilityIdentifier("registrationBirthdayInput")

          if showBirthdayPicker {
            let maxDate = Calendar.current.date(byAdding: .year, value: -13, to: Date()) ?? Date()
            DatePicker(
              "Birthday",
              selection: Binding(
                get: { birthday ?? maxDate },
                set: { birthday = $0 }
              ),
              in: ...maxDate,
              displayedComponents: .date
            )
            .datePickerStyle(.compact)
            .labelsHidden()
            .onChange(of: birthday) { _, _ in
              showBirthdayPicker = false
            }
          }

          HStack(spacing: 12) {
            TextField("First Name", text: $firstName)
              .textFieldStyle(.roundedBorder)
              .disabled(isLoading)
              .accessibilityIdentifier("registrationFirstNameInput")

            TextField("Last Name", text: $lastName)
              .textFieldStyle(.roundedBorder)
              .disabled(isLoading)
              .accessibilityIdentifier("registrationLastNameInput")
          }

          if !errorMessage.isEmpty {
            Text(errorMessage)
              .font(.caption)
              .foregroundColor(.red)
              .multilineTextAlignment(.center)
              .accessibilityIdentifier("errorMessage")
          }

          Button(
            action: {
              registrationTask?.cancel()
              registrationTask = Task { await handleRegistration() }
            },
            label: {
              HStack {
                if isLoading {
                  ProgressView()
                    .progressViewStyle(CircularProgressViewStyle(tint: .white))
                    .scaleEffect(0.8)
                }
                Text(isLoading ? "Registering..." : "Register")
                  .fontWeight(.semibold)
              }
              .frame(maxWidth: .infinity)
              .padding()
              .background(isLoading ? Color.gray : Color.accentColor)
              .foregroundColor(.white)
              .cornerRadius(10)
            }
          )
          .disabled(isLoading || !isFormValid)
          .accessibilityIdentifier("registrationButton")
          .padding(.top, 8)

          // Navigation to login
          Button(
            action: {
              showLogin = true
            },
            label: {
              Text("Already have an account? Sign in")
                .font(.caption)
                .foregroundColor(.accentColor)
            }
          )
          .padding(.top, 8)
        }
        .padding(.horizontal, 32)

        Spacer()
          .frame(height: 20)
      }
    }
    .padding()
  }

  // MARK: - Validation

  private var isFormValid: Bool {
    !emailAddress.trimmingCharacters(in: .whitespaces).isEmpty
      && !username.trimmingCharacters(in: .whitespaces).isEmpty
      && !password.isEmpty
      && !repeatedPassword.isEmpty
      && password == repeatedPassword
      && password.count >= 8
      && isValidEmail(emailAddress)
  }

  private func isValidEmail(_ email: String) -> Bool {
    let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
    let emailPredicate = NSPredicate(format: "SELF MATCHES %@", emailRegex)
    return emailPredicate.evaluate(with: email)
  }

  private func formatDate(_ date: Date) -> String {
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    return formatter.string(from: date)
  }

  // MARK: - Registration

  private func handleRegistration() async {
    guard !Task.isCancelled else { return }

    await MainActor.run {
      errorMessage = ""
      isLoading = true
    }

    // Validate passwords match
    if password != repeatedPassword {
      await MainActor.run {
        errorMessage = "Passwords do not match"
        isLoading = false
      }
      return
    }

    // Validate password length
    if password.count < 8 {
      await MainActor.run {
        errorMessage = "Password must have at least 8 characters"
        isLoading = false
      }
      return
    }

    // Validate email
    if !isValidEmail(emailAddress) {
      await MainActor.run {
        errorMessage = "Invalid email address"
        isLoading = false
      }
      return
    }

    let registrationInput = RegistrationInput(
      emailAddress: emailAddress.trimmingCharacters(in: .whitespaces),
      username: username.trimmingCharacters(in: .whitespaces),
      password: password,
      accountName: accountName.trimmingCharacters(in: .whitespaces),
      firstName: firstName.trimmingCharacters(in: .whitespaces),
      lastName: lastName.trimmingCharacters(in: .whitespaces),
      birthday: birthday,
      invitationToken: invitationToken.trimmingCharacters(in: .whitespaces),
      invitationID: invitationID.trimmingCharacters(in: .whitespaces)
    )
    let result = await authManager.register(input: registrationInput)

    guard !Task.isCancelled else { return }

    await MainActor.run {
      isLoading = false

      if result.success {
        // Registration successful, navigate to login
        clearForm()
        showLogin = true
      } else {
        errorMessage = result.error ?? "Unknown error occurred"
      }
    }
  }

  private func clearForm() {
    emailAddress = ""
    username = ""
    password = ""
    repeatedPassword = ""
    accountName = ""
    firstName = ""
    lastName = ""
    birthday = nil
    invitationToken = ""
    invitationID = ""
  }
}

#Preview {
  RegisterView(showLogin: .constant(false))
    .environment(AuthenticationManager())
}
