//
//  UploadMediaViewModel.swift
//  ios
//

import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2
import SwiftUI

private let uploadChunkSize = 64 * 1024  // 64 KB
private let uploadBucketName = "avatars"

@Observable
@MainActor
class UploadMediaViewModel {
  var isUploading = false
  var errorMessage: String?
  var lastUploadedPath: String?

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  func uploadPhoto(imageData: Data, contentType: String, objectName: String) async {
    isUploading = true
    errorMessage = nil
    lastUploadedPath = nil

    do {
      let (clientManager, metadata) = try await getClientManagerAndMetadata()

      var uploadOptions = GRPCCore.CallOptions.defaults
      uploadOptions.timeout = .seconds(60)

      let response = try await clientManager.client.uploadedMedia.upload(
        metadata: metadata,
        options: uploadOptions,
        requestProducer: { writer in
          // 1. Send metadata
          var meta = UploadedMedia_UploadMetadata()
          meta.bucket = uploadBucketName
          meta.objectName = objectName
          meta.contentType = contentType

          var metadataReq = UploadedMedia_UploadRequest()
          metadataReq.payload = .metadata(meta)
          try await writer.write(metadataReq)

          // 2. Send chunks
          var offset = 0
          while offset < imageData.count {
            let end = min(offset + uploadChunkSize, imageData.count)
            let chunk = imageData.subdata(in: offset..<end)
            offset = end

            var chunkReq = UploadedMedia_UploadRequest()
            chunkReq.payload = .chunk(chunk)
            try await writer.write(chunkReq)
          }
        }
      )

      lastUploadedPath = response.objectURL
    } catch let error as GRPCCore.RPCError {
      let statusMessage = formatRPCError(error)
      errorMessage = "Upload failed: \(statusMessage)"
      print("❌ Upload RPC error: \(error.code), \(error.message)")
    } catch {
      errorMessage = "Upload failed: \(error.localizedDescription)"
    }

    isUploading = false
  }

  private func getClientManagerAndMetadata() async throws -> (
    ClientManager<HTTP2ClientTransport.TransportServices>, GRPCCore.Metadata
  ) {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "UploadMediaViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "UploadMediaViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)
    return (clientManager, metadata)
  }

  private func formatRPCError(_ error: GRPCCore.RPCError) -> String {
    switch error.code {
    case .cancelled:
      return "Request was cancelled. This can happen if the connection was interrupted or the request took too long."
    case .deadlineExceeded:
      return "Request timed out. Try a smaller image."
    case .unauthenticated:
      return "Session expired. Please sign in again."
    case .unavailable:
      return "Server unavailable. Is the backend running at \(APIConfiguration.grpcHost):\(APIConfiguration.grpcPort)?"
    case .permissionDenied:
      return "Permission denied."
    case .invalidArgument:
      return "Invalid request: \(error.message)"
    default:
      return "\(error.code): \(error.message)"
    }
  }
}
