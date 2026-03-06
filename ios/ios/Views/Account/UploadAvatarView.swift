//
//  UploadAvatarView.swift
//  ios
//

import PhotosUI
import SwiftUI
import UIKit

struct UploadAvatarView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: UploadAvatarViewModel?
  @State private var selectedItem: PhotosPickerItem?
  @Environment(\.dismiss) private var dismiss

  var body: some View {
    ScrollView {
      VStack(spacing: DSTheme.Spacing.xl) {
        if let viewModel = viewModel {
          DSSection("Profile Photo", subtitle: "Select a photo to set as your profile picture") {
            VStack(spacing: DSTheme.Spacing.lg) {
              PhotosPicker(
                selection: $selectedItem,
                matching: .images,
                photoLibrary: .shared()
              ) {
                HStack(spacing: DSTheme.Spacing.md) {
                  Image(systemName: "person.crop.circle.badge.plus")
                    .font(.system(size: DSTheme.IconSize.lg))
                    .foregroundColor(DSTheme.Colors.primary)
                  Text("Select Photo")
                    .font(DSTheme.Typography.label)
                    .foregroundColor(DSTheme.Colors.textPrimary)
                }
                .frame(maxWidth: .infinity)
                .padding(DSTheme.Spacing.lg)
                .background(DSTheme.Colors.cardBackground)
                .cornerRadius(DSTheme.Radius.md)
              }
              .disabled(viewModel.isUploading)

              if viewModel.isUploading {
                ProgressView("Uploading...")
                  .padding()
              }

              if let error = viewModel.errorMessage {
                Text(error)
                  .font(DSTheme.Typography.caption)
                  .foregroundColor(DSTheme.Colors.error)
                  .multilineTextAlignment(.leading)
              }

              if viewModel.didSucceed {
                HStack(spacing: DSTheme.Spacing.sm) {
                  Image(systemName: "checkmark.circle.fill")
                    .foregroundColor(DSTheme.Colors.success)
                  Text("Profile photo updated successfully")
                    .font(DSTheme.Typography.label)
                    .foregroundColor(DSTheme.Colors.success)
                }
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(DSTheme.Spacing.md)
                .background(DSTheme.Colors.cardBackground)
                .cornerRadius(DSTheme.Radius.md)
              }
            }
          }
        } else {
          DSInitializingView()
        }
      }
      .dsScreenPadding()
    }
    .navigationTitle("Profile Photo")
    .onAppear {
      if viewModel == nil {
        viewModel = UploadAvatarViewModel(authManager: authManager)
      }
    }
    .onChange(of: selectedItem) { _, newItem in
      guard let item = newItem, let viewModel = viewModel else { return }
      Task {
        await loadAndUpload(item: item, viewModel: viewModel)
      }
    }
    .onChange(of: viewModel?.didSucceed) { _, didSucceed in
      if didSucceed == true {
        dismiss()
      }
    }
  }

  private func loadAndUpload(item: PhotosPickerItem, viewModel: UploadAvatarViewModel) async {
    do {
      guard let data = try await item.loadTransferable(type: Data.self) else {
        viewModel.errorMessage = "Failed to load image data"
        return
      }

      let (uploadData, contentType, fileExtension) = prepareImageForUpload(data: data)
      let objectName = UUID().uuidString + "." + fileExtension

      await viewModel.uploadAvatar(
        imageData: uploadData,
        contentType: contentType,
        objectName: objectName
      )

      selectedItem = nil
    } catch {
      viewModel.errorMessage = "Failed to load image: \(error.localizedDescription)"
    }
  }

  private func prepareImageForUpload(data: Data) -> (Data, String, String) {
    if Self.isHEIC(data: data), let image = UIImage(data: data),
      let jpegData = image.jpegData(compressionQuality: 0.9)
    {
      return (jpegData, "image/jpeg", "jpg")
    }
    if Self.isPNG(data: data) {
      return (data, "image/png", "png")
    }
    if Self.isJPEG(data: data) {
      return (data, "image/jpeg", "jpg")
    }
    if Self.isGIF(data: data) {
      return (data, "image/gif", "gif")
    }
    if let image = UIImage(data: data), let jpegData = image.jpegData(compressionQuality: 0.9) {
      return (jpegData, "image/jpeg", "jpg")
    }
    return (data, "image/jpeg", "jpg")
  }

  private static func isJPEG(data: Data) -> Bool {
    data.count >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
  }

  private static func isPNG(data: Data) -> Bool {
    let signature: [UInt8] = [0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A]
    return data.count >= 8 && data.prefix(8).elementsEqual(signature)
  }

  private static func isGIF(data: Data) -> Bool {
    let signature = "GIF8"
    return data.count >= 4 && String(data: data.prefix(4), encoding: .ascii) == signature
  }

  private static func isHEIC(data: Data) -> Bool {
    guard data.count >= 12 else { return false }
    let ftyp = String(data: data.subdata(in: 4..<8), encoding: .ascii)
    guard ftyp == "ftyp" else { return false }
    let brand = String(data: data.subdata(in: 8..<12), encoding: .ascii)
    return brand == "heic" || brand == "heix" || brand == "mif1"
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "Test User"
  authManager.userID = "user123"

  return NavigationStack {
    UploadAvatarView()
      .environment(authManager)
  }
}
