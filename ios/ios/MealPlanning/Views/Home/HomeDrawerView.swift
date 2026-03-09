//
//  HomeDrawerView.swift
//  ios
//
//  Side drawer for home screen navigation (Browse, Create Meal Plan, Account).
//

import SwiftUI

struct HomeDrawerView: View {
  @Binding var isPresented: Bool

  var displayName: String = ""
  var avatarURL: URL?
  var acceptedOccupiedDates: Set<Date> = []
  var proposedOccupiedDates: Set<Date> = []

  private let drawerWidth: CGFloat = 280

  var body: some View {
    ZStack(alignment: .trailing) {
      // Backdrop - tap to dismiss
      Color.black.opacity(0.3)
        .ignoresSafeArea()
        .onTapGesture {
          isPresented = false
        }

      // Drawer panel (slides in from right)
      VStack(alignment: .leading, spacing: 0) {
        // Close button (visible when drawer is open)
        HStack {
          Spacer()
          Button {
            isPresented = false
          } label: {
            Image(systemName: "xmark")
              .font(.system(size: 18, weight: .bold))
              .foregroundColor(.red)
          }
          .padding(DSTheme.Spacing.md)
          .offset(x: -15)
        }

        // Account section
        VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
          Text("Account")
            .font(DSTheme.Typography.caption)
            .fontWeight(.semibold)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .padding(.horizontal, DSTheme.Spacing.md)
            .padding(.top, DSTheme.Spacing.lg)

          NavigationLink(destination: AccountSettingsView()) {
            HStack(spacing: DSTheme.Spacing.md) {
              DSAvatar(
                name: displayName,
                size: .sm,
                imageURL: avatarURL
              )
              Text(displayName.isEmpty ? "Account" : displayName)
                .font(DSTheme.Typography.bodyLarge)
                .foregroundColor(DSTheme.Colors.textPrimary)
            }
          }
          .buttonStyle(.plain)
          .simultaneousGesture(TapGesture().onEnded { isPresented = false })
          .padding(DSTheme.Spacing.md)
          .padding(.horizontal, DSTheme.Spacing.md)
        }

        // Create Meal Plan CTA
        VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
          Text("Meal Planning")
            .font(DSTheme.Typography.caption)
            .fontWeight(.semibold)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .padding(.horizontal, DSTheme.Spacing.md)
            .padding(.top, DSTheme.Spacing.lg)

          NavigationLink(
            destination: CreateMealPlanWizardView(
              acceptedOccupiedDates: acceptedOccupiedDates,
              proposedOccupiedDates: proposedOccupiedDates
            )
          ) {
            Label("Create Meal Plan", systemImage: "plus.circle.fill")
              .font(DSTheme.Typography.bodyLarge)
              .foregroundColor(DSTheme.Colors.primary)
          }
          .buttonStyle(.plain)
          .simultaneousGesture(TapGesture().onEnded { isPresented = false })
          .padding(DSTheme.Spacing.md)
          .padding(.horizontal, DSTheme.Spacing.md)
        }

        // Browse section
        VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
          Text("Browse")
            .font(DSTheme.Typography.caption)
            .fontWeight(.semibold)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .padding(.horizontal, DSTheme.Spacing.md)
            .padding(.top, DSTheme.Spacing.lg)

          NavigationLink(destination: MealListView()) {
            Label("Meals", systemImage: "fork.knife")
              .font(DSTheme.Typography.bodyLarge)
              .foregroundColor(DSTheme.Colors.textPrimary)
          }
          .buttonStyle(.plain)
          .simultaneousGesture(TapGesture().onEnded { isPresented = false })
          .padding(DSTheme.Spacing.md)
          .padding(.horizontal, DSTheme.Spacing.md)

          NavigationLink(destination: RecipeListView()) {
            Label("Recipes", systemImage: "book.closed.fill")
              .font(DSTheme.Typography.bodyLarge)
              .foregroundColor(DSTheme.Colors.textPrimary)
          }
          .buttonStyle(.plain)
          .simultaneousGesture(TapGesture().onEnded { isPresented = false })
          .padding(DSTheme.Spacing.md)
          .padding(.horizontal, DSTheme.Spacing.md)
        }

        Spacer()
      }
      .frame(width: drawerWidth)
      .background(Color(.systemBackground))
      .offset(x: isPresented ? 0 : drawerWidth)
    }
    .animation(.easeInOut(duration: 0.25), value: isPresented)
  }
}
