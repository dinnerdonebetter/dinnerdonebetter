package config

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config represents our database configuration.
	Config struct {
		_ struct{} `json:"-"`

		OAuth2TokenEncryptionKey string            `json:"oauth2TokenEncryptionKey" toml:"oauth2_token_encryption_key,omitempty"`
		ConnectionDetails        ConnectionDetails `json:"connectionDetails"        toml:"connection_details,omitempty"`
		Debug                    bool              `json:"debug"                    toml:"debug,omitempty"`
		LogQueries               bool              `json:"logQueries"               toml:"log_queries,omitempty"`
		RunMigrations            bool              `json:"runMigrations"            toml:"run_migrations,omitempty"`
		MaxPingAttempts          uint64            `json:"maxPingAttempts"          toml:"max_ping_attempts,omitempty"`
		PingWaitPeriod           time.Duration     `json:"pingWaitPeriod"           toml:"ping_wait_period,omitempty"`
	}

	ConnectionDetails struct {
		_ struct{} `json:"-"`
		
		Username   string `json:"username"`
		Password   string `json:"password"`
		Database   string `json:"database"`
		Host       string `json:"hostname"`
		Port       uint16 `json:"port"`
		DisableSSL bool   `json:"disableSSL"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates an DatabaseSettings struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.ConnectionDetails, validation.Required),
		validation.Field(&cfg.OAuth2TokenEncryptionKey, validation.Required),
	)
}

// LoadConnectionDetailsFromURL wraps an inner function.
func (cfg *Config) LoadConnectionDetailsFromURL(u string) error {
	return cfg.ConnectionDetails.LoadFromURL(u)
}

// ValidateWithContext validates an DatabaseSettings struct.
func (x *ConnectionDetails) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Host, validation.Required),
		validation.Field(&x.Database, validation.Required),
		validation.Field(&x.Username, validation.Required),
		validation.Field(&x.Password, validation.Required),
		validation.Field(&x.Port, validation.Required),
	)
}

var _ fmt.Stringer = (*ConnectionDetails)(nil)

func (x *ConnectionDetails) String() string {
	return fmt.Sprintf(
		"user=%s password=%s database=%s host=%s port=%d",
		x.Username,
		x.Password,
		x.Database,
		x.Host,
		x.Port,
	)
}

func (x *ConnectionDetails) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		x.Username,
		x.Password,
		net.JoinHostPort(x.Host, strconv.FormatUint(uint64(x.Port), 10)),
		x.Database,
	)
}

// LoadFromURL accepts a Postgres connection string and parses it into the ConnectionDetails struct.
func (x *ConnectionDetails) LoadFromURL(u string) error {
	z, err := url.Parse(u)
	if err != nil {
		return err
	}

	port, err := strconv.ParseUint(z.Port(), 10, 64)
	if err != nil {
		return err
	}

	x.Username = z.User.Username()
	x.Password, _ = z.User.Password()
	x.Host = z.Hostname()
	x.Port = uint16(port)
	x.Database = strings.TrimPrefix(z.Path, "/")
	x.DisableSSL = z.Query().Get("sslmode") == "disable"

	return nil
}
