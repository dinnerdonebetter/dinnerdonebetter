module github.com/prixfixeco/api_server

go 1.16

require (
	cloud.google.com/go/pubsub v1.10.3
	cloud.google.com/go/secretmanager v1.1.0
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.3.0
	github.com/GuiaBolso/darwin v0.0.0-20191218124601-fd6d2aa3d244
	github.com/Masterminds/squirrel v1.5.0
	github.com/alexedwards/argon2id v0.0.0-20210326052512-e2135f7c9c77
	github.com/alexedwards/scs/postgresstore v0.0.0-20210407073823-f445396108a4
	github.com/alexedwards/scs/v2 v2.4.0
	github.com/aws/aws-sdk-go v1.40.43
	github.com/boombuler/barcode v1.0.1
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20211207161625-b8fa12c97f1d
	github.com/go-chi/chi/v5 v5.0.4
	github.com/go-chi/cors v1.2.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-redis/redis/v8 v8.11.3
	github.com/goccy/go-graphviz v0.0.9
	github.com/google/uuid v1.2.0
	github.com/google/wire v0.5.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/lib/pq v1.10.2
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/makiuchi-d/gozxing v0.0.0-20210324052758-57132e828831
	github.com/moul/http2curl v1.0.0
	github.com/mssola/user_agent v0.5.2
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/o1egl/paseto v1.0.0
	github.com/pquerna/otp v1.3.0
	github.com/rs/zerolog v1.21.0
	github.com/segmentio/ksuid v1.0.4
	github.com/sendgrid/rest v2.6.5+incompatible
	github.com/sendgrid/sendgrid-go v3.10.3+incompatible
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/stripe/stripe-go/v72 v72.72.0
	github.com/tkrajina/typescriptify-golang-structs v0.1.7
	github.com/unrolled/secure v1.0.8
	github.com/wagslane/go-password-validator v0.3.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.27.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/exporters/prometheus v0.26.0
	go.opentelemetry.io/otel/metric v0.26.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/sdk/export/metric v0.26.0
	go.opentelemetry.io/otel/sdk/metric v0.26.0
	go.opentelemetry.io/otel/trace v1.3.0
	gocloud.dev v0.23.0
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	gonum.org/v1/gonum v0.12.0
	google.golang.org/genproto v0.0.0-20220207185906-7721543eae58
	gopkg.in/launchdarkly/go-sdk-common.v2 v2.5.1
	gopkg.in/launchdarkly/go-server-sdk.v5 v5.10.0
	gopkg.in/mikespook/gorbac.v2 v2.1.0
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
	resenje.org/schulze v0.2.0
)

require (
	cloud.google.com/go/kms v1.2.0 // indirect
	cloud.google.com/go/logging v1.4.2
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/cznic/ql v1.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1
	github.com/kr/text v0.2.0 // indirect
	github.com/mikespook/gorbac v2.1.0+incompatible // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20220317061510-51cd9980dadf // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
