package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"

	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/memblob"
	"gocloud.dev/blob/s3blob"
	"gocloud.dev/gcp"
)

const (
	// MemoryProvider indicates we'd like to use the memory adapter for blob.
	MemoryProvider = "memory"
)

var (
	// ErrNilConfig denotes that the provided configuration is nil.
	ErrNilConfig = errors.New("nil config provided")
)

type (
	// Uploader implements our UploadManager struct.
	Uploader struct {
		bucket          *blob.Bucket
		logger          logging.Logger
		tracer          tracing.Tracer
		filenameFetcher func(req *http.Request) string
	}

	// Config configures our UploadManager.
	Config struct {
		_ struct{} `json:"-"`

		FilesystemConfig  *FilesystemConfig `env:"init"                envPrefix:"FILESYSTEM_"            json:"filesystem,omitempty"`
		S3Config          *S3Config         `env:"init"                envPrefix:"S3_"                    json:"s3,omitempty"`
		GCP               *GCPConfig        `env:"init"                envPrefix:"GCP_"                   json:"gcpConfig,omitempty"`
		BucketPrefix      string            `env:"BUCKET_PREFIX"       json:"bucketPrefix,omitempty"`
		BucketName        string            `env:"BUCKET_NAME"         json:"bucketName,omitempty"`
		UploadFilenameKey string            `env:"UPLOAD_FILENAME_KEY" json:"uploadFilenameKey,omitempty"`
		Provider          string            `env:"PROVIDER"            json:"provider,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the Config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.BucketName, validation.Required),
		validation.Field(&c.Provider, validation.In(S3Provider, FilesystemProvider, MemoryProvider, GCPCloudStorageProvider)),
		validation.Field(&c.S3Config, validation.When(c.Provider == S3Provider, validation.Required).Else(validation.Nil)),
		validation.Field(&c.GCP, validation.When(c.Provider == GCPCloudStorageProvider, validation.Required).Else(validation.Nil)),
		validation.Field(&c.FilesystemConfig, validation.When(c.Provider == FilesystemProvider, validation.Required).Else(validation.Nil)),
	)
}

// NewUploadManager provides a new uploads.UploadManager.
func NewUploadManager(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config, routeParamManager routing.RouteParamManager) (*Uploader, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	serviceName := fmt.Sprintf("%s_uploader", cfg.BucketName)
	u := &Uploader{
		logger:          logging.EnsureLogger(logger).WithName(serviceName),
		tracer:          tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		filenameFetcher: routeParamManager.BuildRouteParamStringIDFetcher(cfg.UploadFilenameKey),
	}

	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, fmt.Errorf("upload manager provided invalid config: %w", err)
	}

	if err := u.selectBucket(ctx, cfg); err != nil {
		return nil, fmt.Errorf("initializing bucket: %w", err)
	}

	return u, nil
}

func (u *Uploader) selectBucket(ctx context.Context, cfg *Config) (err error) {
	switch strings.TrimSpace(strings.ToLower(cfg.Provider)) {
	case S3Provider:
		if cfg.S3Config == nil {
			return ErrNilConfig
		}

		if u.bucket, err = s3blob.OpenBucketV2(ctx, s3v2.New(s3v2.Options{}), cfg.S3Config.BucketName, &s3blob.Options{
			UseLegacyList: false,
		}); err != nil {
			return fmt.Errorf("initializing s3 bucket: %w", err)
		}
	case GCPCloudStorageProvider:
		creds, credsErr := gcp.DefaultCredentials(ctx)
		if credsErr != nil {
			return fmt.Errorf("initializing GCP objectstorage: %w", credsErr)
		}

		client, clientErr := gcp.NewHTTPClient(gcp.DefaultTransport(), creds.TokenSource)
		if clientErr != nil {
			return fmt.Errorf("initializing GCP objectstorage: %w", clientErr)
		}

		u.bucket, err = gcsblob.OpenBucket(ctx, client, cfg.GCP.BucketName, nil)
		if err != nil {
			return fmt.Errorf("initializing GCP objectstorage: %w", err)
		}

		if available, availabilityErr := u.bucket.IsAccessible(ctx); availabilityErr != nil {
			return fmt.Errorf("verifying bucket accessibility: %w", availabilityErr)
		} else if !available {
			return fmt.Errorf("bucket %q is unavailable", cfg.BucketName)
		}

	case MemoryProvider:
		u.bucket = memblob.OpenBucket(&memblob.Options{})
	default:
		if cfg.FilesystemConfig == nil {
			return ErrNilConfig
		}

		if u.bucket, err = fileblob.OpenBucket(cfg.FilesystemConfig.RootDirectory, &fileblob.Options{
			URLSigner: nil,
			CreateDir: true,
		}); err != nil {
			return fmt.Errorf("initializing filesystem bucket: %w", err)
		}
	}

	if cfg.BucketPrefix != "" {
		u.bucket = blob.PrefixedBucket(u.bucket, cfg.BucketPrefix)
	}

	return err
}
