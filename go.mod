module github.com/dinnerdonebetter/backend

go 1.20

require (
	cloud.google.com/go/pubsub v1.31.0
	cloud.google.com/go/secretmanager v1.11.0
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.15.0
	github.com/GuiaBolso/darwin v0.0.0-20191218124601-fd6d2aa3d244
	github.com/Masterminds/squirrel v1.5.4
	github.com/alexedwards/argon2id v0.0.0-20230305115115-4b3c3280a736
	github.com/alexedwards/scs/postgresstore v0.0.0-20230327161757-10d4299e3b24
	github.com/alexedwards/scs/v2 v2.5.1
	github.com/algolia/algoliasearch-client-go/v3 v3.30.0
	github.com/aws/aws-sdk-go v1.44.280
	github.com/boombuler/barcode v1.0.1
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/cockroachdb/cockroachdb-parser v0.0.0-20230515042840-c9c144eab71a
	github.com/daixiang0/gci v0.2.9
	github.com/doug-martin/goqu/v9 v9.18.0
	github.com/elastic/go-elasticsearch/v8 v8.8.1
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/cors v1.2.1
	github.com/go-oauth2/oauth2/v4 v4.5.2
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/goccy/go-graphviz v0.1.1
	github.com/google/wire v0.5.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-retryablehttp v0.7.4
	github.com/heimdalr/dag v1.2.1
	github.com/invopop/jsonschema v0.7.0
	github.com/jinzhu/copier v0.3.5
	github.com/keith-turner/ecoji/v2 v2.0.1
	github.com/kyleconroy/sqlc v1.15.0
	github.com/launchdarkly/go-sdk-common/v3 v3.0.1
	github.com/launchdarkly/go-server-sdk/v6 v6.1.0
	github.com/lib/pq v1.10.9
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/mailgun/mailgun-go/v4 v4.9.0
	github.com/mailjet/mailjet-apiv3-go/v4 v4.0.1
	github.com/matcornic/hermes/v2 v2.1.0
	github.com/mjibson/sqlfmt v0.5.0
	github.com/moul/http2curl v1.0.0
	github.com/mssola/useragent v1.0.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pquerna/otp v1.4.0
	github.com/rs/xid v1.5.0
	github.com/rs/zerolog v1.29.1
	github.com/rudderlabs/analytics-go/v4 v4.1.0
	github.com/sendgrid/rest v2.6.9+incompatible
	github.com/sendgrid/sendgrid-go v3.12.0+incompatible
	github.com/shopspring/decimal v1.3.1
	github.com/stretchr/testify v1.8.4
	github.com/stripe/stripe-go/v72 v72.122.0
	github.com/wagslane/go-password-validator v0.3.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.42.0
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/exporters/jaeger v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
	go.uber.org/automaxprocs v1.5.2
	go.uber.org/zap v1.24.0
	gocloud.dev v0.29.0
	golang.org/x/net v0.10.0
	golang.org/x/oauth2 v0.8.0
	golang.org/x/tools v0.9.3
	gonum.org/v1/gonum v0.13.0
	gopkg.in/mikespook/gorbac.v2 v2.3.0
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
	resenje.org/schulze v0.4.2
)

require (
	cloud.google.com/go v0.110.2 // indirect
	cloud.google.com/go/compute v1.20.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.0 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	cloud.google.com/go/trace v1.10.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.39.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go-v2 v1.18.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.25 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.67 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.33.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.19.0 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/biogo/store v0.0.0-20201120204734-aad293a2328f // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cockroachdb/apd/v3 v3.2.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.4 // indirect
	github.com/dave/dst v0.27.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.3.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/getsentry/sentry-go v0.21.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/geo v0.0.0-20230421003525-6adc56603217 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/google/subcommands v1.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.4 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/launchdarkly/ccache v1.1.0 // indirect
	github.com/launchdarkly/eventsource v1.7.1 // indirect
	github.com/launchdarkly/go-jsonstream/v3 v3.0.0 // indirect
	github.com/launchdarkly/go-sdk-events/v2 v2.0.2 // indirect
	github.com/launchdarkly/go-semver v1.0.2 // indirect
	github.com/launchdarkly/go-server-sdk-evaluation/v2 v2.0.2 // indirect
	github.com/mailjet/mailjet-apiv3-go/v3 v3.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/petermattis/goid v0.0.0-20230518223814-80aa455d8761 // indirect
	github.com/pierrre/geohash v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twpayne/go-geom v1.5.2 // indirect
	github.com/twpayne/go-kml/v2 v2.0.1 // indirect
	github.com/vanng822/css v1.0.1 // indirect
	github.com/vanng822/go-premailer v1.20.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/image v0.7.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/sync v0.2.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

require (
	github.com/cznic/ql v1.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.10.0
	github.com/mikespook/gorbac v2.1.0+incompatible // indirect
	github.com/segmentio/backo-go v1.0.1 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	golang.org/x/crypto v0.9.0
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
