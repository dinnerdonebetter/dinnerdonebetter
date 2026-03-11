package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	uploadedmediagrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"

	g "maragu.dev/gomponents"
)

const (
	maxMediaUploadSize       = 5 * 1024 * 1024 // 5 MB
	uploadChunkSize          = 32 * 1024       // 32 KB
	defaultUploadFilename    = "upload"
	defaultUploadContentType = "application/octet-stream"
)

// streamFileToPreparationMedia uploads a file to the preparation media gRPC stream.
func (s *AdminFrontendServer) streamFileToPreparationMedia(ctx context.Context, validPreparationID, forIngredientID, filename, contentType string, fileData []byte) (string, error) {
	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return "", err
	}

	stream, err := c.UploadPreparationMedia(ctx)
	if err != nil {
		return "", fmt.Errorf("opening upload stream: %w", err)
	}

	// First message: metadata
	firstReq := &mealplanningsvc.UploadPreparationMediaRequest{
		ValidPreparationId: validPreparationID,
		Upload: &uploadedmediagrpc.UploadRequest{
			Payload: &uploadedmediagrpc.UploadRequest_Metadata{
				Metadata: &uploadedmediagrpc.UploadMetadata{
					ObjectName:  filename,
					ContentType: contentType,
				},
			},
		},
	}
	if forIngredientID != "" {
		firstReq.ForIngredientId = &forIngredientID
	}
	if err = stream.Send(firstReq); err != nil {
		return "", fmt.Errorf("sending metadata: %w", err)
	}

	// Stream chunks
	for offset := 0; offset < len(fileData); offset += uploadChunkSize {
		end := min(offset+uploadChunkSize, len(fileData))
		chunk := fileData[offset:end]
		if err = stream.Send(&mealplanningsvc.UploadPreparationMediaRequest{
			Upload: &uploadedmediagrpc.UploadRequest{
				Payload: &uploadedmediagrpc.UploadRequest_Chunk{Chunk: chunk},
			},
		}); err != nil {
			return "", fmt.Errorf("sending chunk: %w", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("closing stream: %w", err)
	}
	if resp == nil || resp.UploadedMediaId == nil {
		return "", fmt.Errorf("no uploaded media id in response")
	}
	return *resp.UploadedMediaId, nil
}

// streamFileToIngredientMedia uploads a file to the ingredient media gRPC stream.
func (s *AdminFrontendServer) streamFileToIngredientMedia(ctx context.Context, validIngredientID, filename, contentType string, fileData []byte) (string, error) {
	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return "", err
	}

	stream, err := c.UploadIngredientMedia(ctx)
	if err != nil {
		return "", fmt.Errorf("opening upload stream: %w", err)
	}

	// First message: metadata
	if err = stream.Send(&mealplanningsvc.UploadIngredientMediaRequest{
		ValidIngredientId: validIngredientID,
		Upload: &uploadedmediagrpc.UploadRequest{
			Payload: &uploadedmediagrpc.UploadRequest_Metadata{
				Metadata: &uploadedmediagrpc.UploadMetadata{
					ObjectName:  filename,
					ContentType: contentType,
				},
			},
		},
	}); err != nil {
		return "", fmt.Errorf("sending metadata: %w", err)
	}

	// Stream chunks
	for offset := 0; offset < len(fileData); offset += uploadChunkSize {
		end := min(offset+uploadChunkSize, len(fileData))
		chunk := fileData[offset:end]
		if err = stream.Send(&mealplanningsvc.UploadIngredientMediaRequest{
			Upload: &uploadedmediagrpc.UploadRequest{
				Payload: &uploadedmediagrpc.UploadRequest_Chunk{Chunk: chunk},
			},
		}); err != nil {
			return "", fmt.Errorf("sending chunk: %w", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("closing stream: %w", err)
	}
	if resp == nil || resp.UploadedMediaId == nil {
		return "", fmt.Errorf("no uploaded media id in response")
	}
	return *resp.UploadedMediaId, nil
}

// streamFileToRecipeStepImage uploads a file to the recipe step image gRPC stream.
func (s *AdminFrontendServer) streamFileToRecipeStepImage(ctx context.Context, recipeID, recipeStepID, filename, contentType string, fileData []byte) (string, error) {
	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return "", err
	}

	stream, err := c.UploadRecipeStepImage(ctx)
	if err != nil {
		return "", fmt.Errorf("opening upload stream: %w", err)
	}

	// First message: metadata
	if err = stream.Send(&mealplanningsvc.UploadRecipeStepImageRequest{
		RecipeId:     recipeID,
		RecipeStepId: recipeStepID,
		Upload: &uploadedmediagrpc.UploadRequest{
			Payload: &uploadedmediagrpc.UploadRequest_Metadata{
				Metadata: &uploadedmediagrpc.UploadMetadata{
					ObjectName:  filename,
					ContentType: contentType,
				},
			},
		},
	}); err != nil {
		return "", fmt.Errorf("sending metadata: %w", err)
	}

	// Stream chunks
	for offset := 0; offset < len(fileData); offset += uploadChunkSize {
		end := min(offset+uploadChunkSize, len(fileData))
		chunk := fileData[offset:end]
		if err = stream.Send(&mealplanningsvc.UploadRecipeStepImageRequest{
			Upload: &uploadedmediagrpc.UploadRequest{
				Payload: &uploadedmediagrpc.UploadRequest_Chunk{Chunk: chunk},
			},
		}); err != nil {
			return "", fmt.Errorf("sending chunk: %w", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("closing stream: %w", err)
	}
	if resp == nil || resp.UploadedMediaId == nil {
		return "", fmt.Errorf("no uploaded media id in response")
	}
	return *resp.UploadedMediaId, nil
}

// UploadPreparationMedia handles POST /api/valid_preparations/{id}/media (multipart).
func (s *AdminFrontendServer) UploadPreparationMedia(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	validPreparationID := s.validPreparationIDRouteParamFetcher(req)
	if validPreparationID == "" {
		http.Error(res, "valid_preparation_id is required", http.StatusBadRequest)
		return g.El("div"), nil
	}

	if err := req.ParseMultipartForm(maxMediaUploadSize); err != nil {
		http.Error(res, fmt.Sprintf("failed to parse multipart form: %v", err), http.StatusBadRequest)
		return g.El("div"), nil
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(res, fmt.Sprintf("file is required: %v", err), http.StatusBadRequest)
		return g.El("div"), nil
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			s.logger.Error("failed to close uploaded file", closeErr)
		}
	}()

	fileData, err := io.ReadAll(io.LimitReader(file, maxMediaUploadSize))
	if err != nil {
		http.Error(res, fmt.Sprintf("failed to read file: %v", err), http.StatusInternalServerError)
		return g.El("div"), nil
	}
	if len(fileData) == 0 {
		http.Error(res, "file is empty", http.StatusBadRequest)
		return g.El("div"), nil
	}

	filename := header.Filename
	if filename == "" {
		filename = defaultUploadFilename
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = defaultUploadContentType
	}

	forIngredientID := req.FormValue("for_ingredient_id")

	_, err = s.streamFileToPreparationMedia(ctx, validPreparationID, forIngredientID, filename, contentType, fileData)
	if err != nil {
		s.logger.Error("failed to upload preparation media", err)
		http.Error(res, fmt.Sprintf("upload failed: %v", err), http.StatusInternalServerError)
		return g.El("div"), nil
	}

	http.Redirect(res, req, fmt.Sprintf("/valid_preparations/%s", validPreparationID), http.StatusSeeOther)
	return g.El("div"), nil
}

// UploadIngredientMedia handles POST /api/valid_ingredients/{id}/media (multipart).
func (s *AdminFrontendServer) UploadIngredientMedia(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	validIngredientID := s.validIngredientIDRouteParamFetcher(req)
	if validIngredientID == "" {
		http.Error(res, "valid_ingredient_id is required", http.StatusBadRequest)
		return g.El("div"), nil
	}

	if err := req.ParseMultipartForm(maxMediaUploadSize); err != nil {
		http.Error(res, fmt.Sprintf("failed to parse multipart form: %v", err), http.StatusBadRequest)
		return g.El("div"), nil
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(res, fmt.Sprintf("file is required: %v", err), http.StatusBadRequest)
		return g.El("div"), nil
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			s.logger.Error("failed to close uploaded file", closeErr)
		}
	}()

	fileData, err := io.ReadAll(io.LimitReader(file, maxMediaUploadSize))
	if err != nil {
		http.Error(res, fmt.Sprintf("failed to read file: %v", err), http.StatusInternalServerError)
		return g.El("div"), nil
	}
	if len(fileData) == 0 {
		http.Error(res, "file is empty", http.StatusBadRequest)
		return g.El("div"), nil
	}

	filename := header.Filename
	if filename == "" {
		filename = defaultUploadFilename
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = defaultUploadContentType
	}

	_, err = s.streamFileToIngredientMedia(ctx, validIngredientID, filename, contentType, fileData)
	if err != nil {
		s.logger.Error("failed to upload ingredient media", err)
		http.Error(res, fmt.Sprintf("upload failed: %v", err), http.StatusInternalServerError)
		return g.El("div"), nil
	}

	http.Redirect(res, req, fmt.Sprintf("/valid_ingredients/%s", validIngredientID), http.StatusSeeOther)
	return g.El("div"), nil
}

// UploadRecipeStepImage handles POST /api/recipes/{recipeID}/steps/{recipeStepID}/images (multipart).
func (s *AdminFrontendServer) UploadRecipeStepImage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	recipeID := s.recipeIDRouteParamFetcher(req)
	recipeStepID := s.recipeStepIDRouteParamFetcher(req)
	if recipeID == "" || recipeStepID == "" {
		http.Error(res, "recipe_id and recipe_step_id are required", http.StatusBadRequest)
		return g.El("div"), nil
	}

	if err := req.ParseMultipartForm(maxMediaUploadSize); err != nil {
		http.Error(res, fmt.Sprintf("failed to parse multipart form: %v", err), http.StatusBadRequest)
		return g.El("div"), nil
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(res, fmt.Sprintf("file is required: %v", err), http.StatusBadRequest)
		return g.El("div"), nil
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			s.logger.Error("failed to close uploaded file", closeErr)
		}
	}()

	fileData, err := io.ReadAll(io.LimitReader(file, maxMediaUploadSize))
	if err != nil {
		http.Error(res, fmt.Sprintf("failed to read file: %v", err), http.StatusInternalServerError)
		return g.El("div"), nil
	}
	if len(fileData) == 0 {
		http.Error(res, "file is empty", http.StatusBadRequest)
		return g.El("div"), nil
	}

	filename := header.Filename
	if filename == "" {
		filename = defaultUploadFilename
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = defaultUploadContentType
	}

	_, err = s.streamFileToRecipeStepImage(ctx, recipeID, recipeStepID, filename, contentType, fileData)
	if err != nil {
		s.logger.Error("failed to upload recipe step image", err)
		http.Error(res, fmt.Sprintf("upload failed: %v", err), http.StatusInternalServerError)
		return g.El("div"), nil
	}

	http.Redirect(res, req, fmt.Sprintf("/recipes/%s", recipeID), http.StatusSeeOther)
	return g.El("div"), nil
}
