module github.com/prixfixeco/api_server

go 1.19

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
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/heimdalr/dag v1.2.1
	github.com/kyleconroy/sqlc v1.15.0
	github.com/lib/pq v1.10.6
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/makiuchi-d/gozxing v0.0.0-20210324052758-57132e828831
	github.com/moul/http2curl v1.0.0
	github.com/mssola/user_agent v0.5.2
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/o1egl/paseto v1.0.0
	github.com/pquerna/otp v1.3.0
	github.com/rs/xid v1.2.1
	github.com/rs/zerolog v1.21.0
	github.com/segmentio/ksuid v1.0.4
	github.com/sendgrid/rest v2.6.5+incompatible
	github.com/sendgrid/sendgrid-go v3.10.3+incompatible
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.1
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
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/tools v0.1.12
	gonum.org/v1/gonum v0.12.0
	google.golang.org/genproto v0.0.0-20220207185906-7721543eae58
	gopkg.in/launchdarkly/go-sdk-common.v2 v2.5.1
	gopkg.in/launchdarkly/go-server-sdk.v5 v5.10.0
	gopkg.in/mikespook/gorbac.v2 v2.1.0
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
	resenje.org/schulze v0.2.0
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/compute v1.2.0 // indirect
	cloud.google.com/go/iam v0.1.0 // indirect
	cloud.google.com/go/storage v1.15.0 // indirect
	cloud.google.com/go/trace v1.0.0 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/chacha20poly1305 v0.0.0-20170617001512-233f39982aeb // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/elastic-transport-go/v8 v8.0.0-20211202110751-50105067ef27 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/gregjones/httpcache v0.0.0-20171119193500-2bcd89a1743f // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/launchdarkly/ccache v1.1.0 // indirect
	github.com/launchdarkly/eventsource v1.6.2 // indirect
	github.com/launchdarkly/go-semver v1.0.2 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.26.0 // indirect
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.68.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.44.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/launchdarkly/go-jsonstream.v1 v1.0.1 // indirect
	gopkg.in/launchdarkly/go-sdk-events.v1 v1.1.1 // indirect
	gopkg.in/launchdarkly/go-server-sdk-evaluation.v1 v1.5.0 // indirect
)

require (
	cloud.google.com/go/kms v1.2.0 // indirect
	cloud.google.com/go/logging v1.4.2
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/cznic/ql v1.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1
	github.com/mikespook/gorbac v2.1.0+incompatible // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
