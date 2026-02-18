package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	paymentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/payments"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	productIDURLParamKey = "productID"
)

func (s *AdminFrontendServer) ProductCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Product", s.renderProductsError("Error: No API client available")), nil
	}

	var input *paymentssvc.ProductCreationRequestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Product", s.renderProductsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	createRes, err := c.CreateProduct(ctx, &paymentssvc.CreateProductRequest{
		Input: input,
	})
	if err != nil {
		return page("New Product", s.renderProductsError(fmt.Sprintf("Error creating product: %v", err))), nil
	}

	if createRes == nil || createRes.Created == nil {
		return page("New Product", s.renderProductsError("Error: No product returned from server")), nil
	}

	http.Redirect(res, req, fmt.Sprintf("/products/%s", createRes.Created.Id), http.StatusSeeOther)
	return g.El("div"), nil
}

func (s *AdminFrontendServer) ProductNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	if _, err := fetchClientFromContext(ctx); err != nil {
		return page("New Product", s.renderProductsError("Error: No API client available")), nil
	}

	emptyInput := &paymentssvc.ProductCreationRequestInput{
		Kind:     "recurring",
		Currency: "usd",
	}

	formPageResult, err := components.FormPage(&components.FormPageProps[*paymentssvc.ProductCreationRequestInput]{
		Title:        "Create New Product",
		BaseSubtitle: "Add a new product to the catalog",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*paymentssvc.ProductCreationRequestInput]{
			FormID: "create-product-form",
			Action: "/api/products",
			Method: "POST",
			EnabledFields: []string{
				"name",
				"description",
				"kind",
				"amount_cents",
				"currency",
				"billing_interval_months",
				"external_product_id",
			},
			FieldConfigs: map[string]*components.FieldConfig{
				"name": {
					Placeholder: "Enter product name...",
					Validation:  &components.FieldValidation{Required: true},
				},
				"description": {
					Placeholder: "Enter description...",
					Validation:  &components.FieldValidation{Required: true},
				},
				"kind": {
					Options: []*components.SelectOption{
						{Value: "recurring", Label: "Recurring", IsDefault: true},
						{Value: "one_time", Label: "One Time"},
					},
					Placeholder: "Select product kind...",
					Validation:  &components.FieldValidation{Required: true, CustomMessage: "Please select a product kind"},
				},
				"amount_cents": {
					Placeholder: "Amount in cents (e.g., 999 for $9.99)...",
					InputType:   "number",
				},
				"currency": {
					Placeholder: "Currency code (e.g., usd)...",
				},
				"billing_interval_months": {
					Placeholder: "Billing interval in months (recurring only, e.g., 1 for monthly)...",
					InputType:   "number",
				},
				"external_product_id": {
					Placeholder: "External product ID from payment provider...",
				},
			},
			FormRows: []*components.FormRow{
				{Fields: []string{"name"}, Columns: 1},
				{Fields: []string{"description"}, Columns: 1},
				{Fields: []string{"kind", "amount_cents", "currency"}, Columns: 3},
				{Fields: []string{"billing_interval_months", "external_product_id"}, Columns: 2},
			},
			SubmitButtonText: "Create Product",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/products",
			HTMXTarget:       "body",
			HTMXSwap:         "innerHTML",
			HTMXPushURL:      true,
			HTMXExtension:    "json-enc",
		},
		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Products", URL: "/products"},
			{Text: "New Product", URL: ""},
		},
	})
	if err != nil {
		return page("New Product", s.renderProductsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("New Product", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ProductPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Products", s.renderProductsError("Error: No API client available")), nil
	}

	productID := s.productIDRouteParamFetcher(req)
	if productID == "" {
		return page("Products", s.renderProductsError("Error: No product ID provided")), nil
	}

	productRes, err := c.GetProduct(ctx, &paymentssvc.GetProductRequest{ProductId: productID})
	if err != nil {
		return page("Products", s.renderProductsError(fmt.Sprintf("Error loading product: %v", err))), nil
	}
	product := productRes.Result

	formPageResult, err := components.FormPage(&components.FormPageProps[*paymentssvc.Product]{
		Title:        "Product Details",
		BaseSubtitle: "View and edit product information",
		Palette:      &design.StandardPalette,
		Data:         product,
		FormOptions: &components.FormOptions[*paymentssvc.Product]{
			FormID: "edit-product-form",
			Action: fmt.Sprintf("/api/products/%s", product.Id),
			Method: "POST",
			EnabledFields: []string{
				"name",
				"description",
				"kind",
				"amount_cents",
				"currency",
				"billing_interval_months",
				"external_product_id",
			},
			FieldConfigs: map[string]*components.FieldConfig{
				"name": {
					Validation: &components.FieldValidation{Required: true},
				},
				"description": {
					Validation: &components.FieldValidation{Required: true},
				},
				"kind": {
					Options: []*components.SelectOption{
						{Value: "recurring", Label: "Recurring"},
						{Value: "one_time", Label: "One Time"},
					},
					Validation: &components.FieldValidation{Required: true},
				},
				"amount_cents":            {InputType: "number"},
				"currency":                {},
				"billing_interval_months": {InputType: "number"},
				"external_product_id":     {},
			},
			FormRows: []*components.FormRow{
				{Fields: []string{"name"}, Columns: 1},
				{Fields: []string{"description"}, Columns: 1},
				{Fields: []string{"kind", "amount_cents", "currency"}, Columns: 3},
				{Fields: []string{"billing_interval_months", "external_product_id"}, Columns: 2},
			},
			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Products",
			CancelURL:        "/products",
			HTMXTarget:       "body",
			HTMXSwap:         "innerHTML",
			HTMXPushURL:      true,
			HTMXExtension:    "json-enc",
		},
		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Products", URL: "/products"},
			{Text: product.Name, URL: ""},
		},
		SubtitleGenerator: func(p *paymentssvc.Product) string {
			return fmt.Sprintf("Editing product: %s", p.Name)
		},
	})
	if err != nil {
		return page("Products", s.renderProductsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Products", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ProductUpdate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Products", s.renderProductsError("Error: No API client available")), nil
	}

	productID := s.productIDRouteParamFetcher(req)
	if productID == "" {
		return page("Products", s.renderProductsError("Error: No product ID provided")), nil
	}

	var input *paymentssvc.ProductUpdateRequestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("Products", s.renderProductsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	_, err = c.UpdateProduct(ctx, &paymentssvc.UpdateProductRequest{
		ProductId: productID,
		Input:     input,
	})
	if err != nil {
		return page("Products", s.renderProductsError(fmt.Sprintf("Error updating product: %v", err))), nil
	}

	http.Redirect(res, req, fmt.Sprintf("/products/%s", productID), http.StatusSeeOther)
	return g.El("div"), nil
}

func (s *AdminFrontendServer) ProductsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Products", s.renderProductsError("Error: No API client available")), nil
	}

	_, grpcFilter := buildQueryFilterFromRequest(req)
	productsRes, err := c.GetProducts(ctx, &paymentssvc.GetProductsRequest{Filter: grpcFilter})
	if err != nil {
		return page("Products", s.renderProductsError(fmt.Sprintf("Error loading products: %v", err))), nil
	}

	tablePageResult, err := components.TablePage(&components.TablePageProps[*paymentssvc.Product]{
		Title:             "Products",
		BaseSubtitle:      "Manage payment products",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search products...",
		HTMXSearchTarget:  "/api/products/search",
		Data:              productsRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Product", "/products/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*paymentssvc.Product]{
			TableID: "products-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"id",
				"name",
				"description",
				"kind",
				"amount_cents",
				"currency",
				"billing_interval_months",
				"external_product_id",
				"created_at",
				"last_updated_at",
				"archived_at",
			},
			FieldReplacements: map[string]string{
				"amount_cents":            "Amount (¢)",
				"billing_interval_months": "Billing Interval",
				"external_product_id":     "External ID",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"created_at":      renderTimestamp,
				"last_updated_at": renderTimestamp,
				"archived_at":     renderTimestamp,
			},
		},
		RowLinkGenerator: func(data *paymentssvc.Product) string {
			return fmt.Sprintf("/products/%s", data.Id)
		},
		EmptyStateTitle:       "No products found",
		EmptyStateDescription: "Get started by creating your first product.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Product", "/products/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage payment products"
			}
			return fmt.Sprintf("Manage %d products", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Products", s.renderProductsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Products", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ProductsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No API client available"),
			),
		), nil
	}

	searchQuery := req.URL.Query().Get("search")
	productsRes, err := c.GetProducts(ctx, &paymentssvc.GetProductsRequest{
		Filter: &grpcfiltering.QueryFilter{},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading products: %v", err)),
			),
		), nil
	}

	var filtered []*paymentssvc.Product
	if searchQuery == "" {
		filtered = productsRes.Results
	} else {
		q := strings.ToLower(searchQuery)
		for _, p := range productsRes.Results {
			if strings.Contains(strings.ToLower(p.Name), q) ||
				strings.Contains(strings.ToLower(p.Description), q) ||
				strings.Contains(strings.ToLower(p.Kind), q) {
				filtered = append(filtered, p)
			}
		}
	}

	if len(filtered) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No products found",
				fmt.Sprintf("No products match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{
					components.ActionButton("Create New Product", "/products/new", &design.StandardPalette, true),
				},
			),
		), nil
	}

	table, err := components.Table(filtered, &components.TableOptions[*paymentssvc.Product]{
		TableID: "products-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"id",
			"name",
			"description",
			"kind",
			"amount_cents",
			"currency",
			"billing_interval_months",
			"external_product_id",
			"created_at",
			"last_updated_at",
			"archived_at",
		},
		FieldReplacements: map[string]string{
			"amount_cents":            "Amount (¢)",
			"billing_interval_months": "Billing Interval",
			"external_product_id":     "External ID",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"created_at":      renderTimestamp,
			"last_updated_at": renderTimestamp,
			"archived_at":     renderTimestamp,
		},
		RowLinkGenerator: func(data *paymentssvc.Product) string {
			return fmt.Sprintf("/products/%s", data.Id)
		},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error creating table: %v", err)),
			),
		), nil
	}

	return g.El("div",
		g.Attr("class", "overflow-x-auto"),
		table,
	), nil
}

func (s *AdminFrontendServer) renderProductsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Products",
		Subtitle: "Manage payment products",
		Palette:  &design.StandardPalette,
	},
		components.Card(&design.StandardPalette,
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(errorMsg),
			),
		),
	)
}
