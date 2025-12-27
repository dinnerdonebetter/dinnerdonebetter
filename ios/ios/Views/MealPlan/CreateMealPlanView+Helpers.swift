//
//  CreateMealPlanView+Helpers.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - CreateMealPlanView Helpers

extension CreateMealPlanView {
  // MARK: - Meal Plan Details Section

  func mealPlanDetailsSection(viewModel: CreateMealPlanViewModel) -> some View {
    let bindableViewModel = Bindable(viewModel)
    
    return VStack(alignment: .leading, spacing: 16) {
      Text("Meal Plan Details")
        .font(.title2)
        .fontWeight(.bold)

      VStack(alignment: .leading, spacing: 12) {
        Text("Name")
          .font(.headline)
        TextField("Enter meal plan name", text: bindableViewModel.mealPlanName)
          .textFieldStyle(.roundedBorder)
          .focused($focusedField, equals: .mealPlanName)
      }

      VStack(alignment: .leading, spacing: 12) {
        Text("Voting Deadline")
          .font(.headline)
        DatePicker(
          "",
          selection: bindableViewModel.votingDeadline,
          displayedComponents: [.date, .hourAndMinute]
        )
        .datePickerStyle(.compact)
        .labelsHidden()
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }
  
  // MARK: - Meal Plan Details Section (Horizontal for iPad)
  
  func mealPlanDetailsSectionHorizontal(viewModel: CreateMealPlanViewModel) -> some View {
    let bindableViewModel = Bindable(viewModel)
    
    return VStack(alignment: .leading, spacing: 16) {
      Text("Meal Plan Details")
        .font(.title2)
        .fontWeight(.bold)

      HStack(alignment: .top, spacing: 24) {
        VStack(alignment: .leading, spacing: 12) {
          Text("Name")
            .font(.headline)
          TextField("Enter meal plan name", text: bindableViewModel.mealPlanName)
            .textFieldStyle(.roundedBorder)
            .focused($focusedField, equals: .mealPlanName)
        }
        .frame(maxWidth: .infinity)

        VStack(alignment: .leading, spacing: 12) {
          Text("Voting Deadline")
            .font(.headline)
          DatePicker(
            "",
            selection: bindableViewModel.votingDeadline,
            displayedComponents: [.date, .hourAndMinute]
          )
          .datePickerStyle(.compact)
          .labelsHidden()
          .frame(maxWidth: .infinity, alignment: .leading)
        }
        .frame(maxWidth: .infinity)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Events Section

  func eventsSection(viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 16) {
      HStack {
        Text("Events")
          .font(.title2)
          .fontWeight(.bold)
        Spacer()
        Button(action: {
          viewModel.addEvent()
        }, label: {
          HStack {
            Image(systemName: "plus.circle.fill")
            Text("Add Event")
          }
          .font(.subheadline)
          .foregroundColor(.blue)
        })
      }

      // Always use vertical stack for events
      ForEach(viewModel.events) { event in
        eventCard(event: event, viewModel: viewModel)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Error Message

  func errorMessage(_ message: String) -> some View {
    HStack {
      Image(systemName: "exclamationmark.triangle")
        .foregroundColor(.red)
      Text(message)
        .font(.subheadline)
        .foregroundColor(.red)
    }
    .padding()
    .background(Color.red.opacity(0.1))
    .cornerRadius(8)
  }

  // MARK: - Create Button

  func createButton(viewModel: CreateMealPlanViewModel) -> some View {
    _ = Bindable(viewModel)
    let hasSelectedMeals = viewModel.events.contains { !$0.selectedMeals.isEmpty }
    
    return Button(action: {
      Task {
        let success = await viewModel.createMealPlan()
        if success {
          // Post notification to refresh home view
          NotificationCenter.default.post(name: .mealPlanCreated, object: nil)
          dismiss()
        }
      }
    }, label: {
      HStack {
        if viewModel.isCreating {
          ProgressView()
            .progressViewStyle(CircularProgressViewStyle(tint: .white))
        }
        Text(viewModel.isCreating ? "Creating..." : "Create Meal Plan")
          .fontWeight(.semibold)
      }
      .frame(maxWidth: .infinity)
      .padding()
      .background(
        viewModel.isCreating || !hasSelectedMeals
          ? Color.gray : Color.blue
      )
      .foregroundColor(.white)
      .cornerRadius(10)
    })
    .disabled(viewModel.isCreating || !hasSelectedMeals)
  }
}

