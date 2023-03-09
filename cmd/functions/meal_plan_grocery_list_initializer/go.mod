module github.com/prixfixeco/backend/cmd/functions/meal_plan_grocery_list_initializer

go 1.19

replace github.com/prixfixeco/backend => ../../../

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.6.1
	github.com/cloudevents/sdk-go/v2 v2.12.0
	github.com/prixfixeco/backend v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.11.1
	go.uber.org/automaxprocs v1.5.1
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/compute v1.6.1 // indirect
	cloud.google.com/go/functions v1.0.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/pubsub v1.10.3 // indirect
	cloud.google.com/go/secretmanager v1.1.0 // indirect
	cloud.google.com/go/storage v1.22.1 // indirect
	cloud.google.com/go/trace v1.0.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.3.0 // indirect
	github.com/GuiaBolso/darwin v0.0.0-20191218124601-fd6d2aa3d244 // indirect
	github.com/Masterminds/squirrel v1.5.0 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/chacha20poly1305 v0.0.0-20170617001512-233f39982aeb // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/alexedwards/argon2id v0.0.0-20210326052512-e2135f7c9c77 // indirect
	github.com/alexedwards/scs/postgresstore v0.0.0-20210407073823-f445396108a4 // indirect
	github.com/alexedwards/scs/v2 v2.4.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.40.43 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/cznic/fileutil v0.0.0-20181122101858-4d67cfea8c87 // indirect
	github.com/cznic/golex v0.0.0-20181122101858-9c343928389c // indirect
	github.com/cznic/internal v0.0.0-20181122101858-3279554c546e // indirect
	github.com/cznic/lldb v1.1.0 // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/cznic/sortutil v0.0.0-20181122101858-f5f958428db8 // indirect
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537 // indirect
	github.com/cznic/zappy v0.0.0-20181122101859-ca47d358d4b1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/go-redis/redis/v8 v8.11.3 // indirect
	github.com/goccy/go-graphviz v0.0.9 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/heimdalr/dag v1.2.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/keith-turner/ecoji/v2 v2.0.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/luna-duclos/instrumentedsql v1.1.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/mssola/user_agent v0.5.2 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/o1egl/paseto v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/otp v1.3.0 // indirect
	github.com/prometheus/client_golang v1.11.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20220927061507-ef77025ab5aa // indirect
	github.com/rs/xid v1.2.1 // indirect
	github.com/rs/zerolog v1.21.0 // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
	github.com/sendgrid/rest v2.6.5+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.10.3+incompatible // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/wagslane/go-password-validator v0.3.0 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.27.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.26.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.3.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	gocloud.dev v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
	google.golang.org/api v0.78.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220518221133-4f43b3371335 // indirect
	google.golang.org/grpc v1.46.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/mikespook/gorbac.v2 v2.1.0 // indirect
	gopkg.in/segmentio/analytics-go.v3 v3.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	resenje.org/schulze v0.2.0 // indirect
)
