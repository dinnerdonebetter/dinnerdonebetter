//
//  UploadMediaView.swift
//  ios
//

import PhotosUI
import SwiftUI

private struct PreparedImage {
  let data: Data
  let contentType: String
  let fileExtension: String
}

struct UploadMediaView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: UploadMediaViewModel?
  @State private var selectedItem: PhotosPickerItem?

  var body: some View {
    ScrollView {
      VStack(spacing: DSTheme.Spacing.xl) {
        if let viewModel = viewModel {
          DSSection("Upload Photo", subtitle: "Select a photo to upload to the backend") {
            VStack(spacing: DSTheme.Spacing.lg) {
              PhotosPicker(
                selection: $selectedItem,
                matching: .images,
                photoLibrary: .shared()
              ) {
                HStack(spacing: DSTheme.Spacing.md) {
                  Image(systemName: "photo.on.rectangle.angled")
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

              if let path = viewModel.lastUploadedPath {
                VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
                  Text("Uploaded successfully")
                    .font(DSTheme.Typography.label)
                    .foregroundColor(DSTheme.Colors.success)
                  Text(path)
                    .font(DSTheme.Typography.caption)
                    .foregroundColor(DSTheme.Colors.textSecondary)
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
    .navigationTitle("Upload Media")
    .onAppear {
      if viewModel == nil {
        viewModel = UploadMediaViewModel(authManager: authManager)
      }
    }
    .onChange(of: selectedItem) { _, newItem in
      guard let item = newItem, let viewModel = viewModel else { return }
      Task {
        await loadAndUpload(item: item, viewModel: viewModel)
      }
    }
  }

  private func loadAndUpload(item: PhotosPickerItem, viewModel: UploadMediaViewModel) async {
    do {
      guard let data = try await item.loadTransferable(type: Data.self) else {
        viewModel.errorMessage = "Failed to load image data"
        return
      }

      let prepared = prepareImageForUpload(data: data)
      let objectName = UUID().uuidString + "." + prepared.fileExtension

      await viewModel.uploadPhoto(
        imageData: prepared.data,
        contentType: prepared.contentType,
        objectName: objectName
      )

      selectedItem = nil
    } catch {
      viewModel.errorMessage = "Failed to load image: \(error.localizedDescription)"
    }
  }

  /// Prepares image data for upload. Converts HEIC to JPEG since backend doesn't support HEIC.
  private func prepareImageForUpload(data: Data) -> PreparedImage {
    if Self.isHEIC(data: data), let image = UIImage(data: data),
      let jpegData = image.jpegData(compressionQuality: 0.9)
    {
      return PreparedImage(data: jpegData, contentType: "image/jpeg", fileExtension: "jpg")
    }
    if Self.isPNG(data: data) {
      return PreparedImage(data: data, contentType: "image/png", fileExtension: "png")
    }
    if Self.isJPEG(data: data) {
      return PreparedImage(data: data, contentType: "image/jpeg", fileExtension: "jpg")
    }
    if Self.isGIF(data: data) {
      return PreparedImage(data: data, contentType: "image/gif", fileExtension: "gif")
    }
    if let image = UIImage(data: data), let jpegData = image.jpegData(compressionQuality: 0.9) {
      return PreparedImage(data: jpegData, contentType: "image/jpeg", fileExtension: "jpg")
    }
    return PreparedImage(data: data, contentType: "image/jpeg", fileExtension: "jpg")
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
    UploadMediaView()
      .environment(authManager)
  }
}
