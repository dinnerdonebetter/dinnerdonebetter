module gitlab.com/prixfixe/prixfixe

go 1.17

require (
	github.com/Azure/azure-pipeline-go v0.2.3
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/GuiaBolso/darwin v0.0.0-20191218124601-fd6d2aa3d244
	github.com/Masterminds/squirrel v1.5.0
	github.com/RoaringBitmap/roaring v0.6.0 // indirect
	github.com/aead/chacha20poly1305 v0.0.0-20201124145622-1a5aba2a8b29 // indirect
	github.com/alexedwards/argon2id v0.0.0-20210326052512-e2135f7c9c77
	github.com/alexedwards/scs/postgresstore v0.0.0-20210407073823-f445396108a4
	github.com/alexedwards/scs/v2 v2.4.0
	github.com/aws/aws-sdk-go v1.38.35
	github.com/blevesearch/bleve/v2 v2.0.3
	github.com/boombuler/barcode v1.0.1
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/carolynvs/magex v0.5.0
	github.com/cznic/ql v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/cors v1.2.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/uuid v1.2.0
	github.com/google/wire v0.5.0
	github.com/gorilla/securecookie v1.1.1
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/vault/api v1.1.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/lib/pq v1.10.1
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/magefile/mage v1.11.0
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/makiuchi-d/gozxing v0.0.0-20210324052758-57132e828831
	github.com/mikespook/gorbac v2.1.0+incompatible // indirect
	github.com/moul/http2curl v1.0.0
	github.com/mssola/user_agent v0.5.2
	github.com/mxschmitt/playwright-go v0.1100.0
	github.com/nats-io/nats.go v1.11.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/o1egl/paseto v1.0.0
	github.com/pelletier/go-toml v1.9.0 // indirect
	github.com/pquerna/otp v1.3.0
	github.com/prometheus/common v0.23.0 // indirect
	github.com/rs/zerolog v1.21.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/stripe/stripe-go/v72 v72.45.0
	github.com/unrolled/secure v1.0.8
	github.com/wagslane/go-password-validator v0.3.0
	github.com/willf/bitset v1.1.11 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.19.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.19.0
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.19.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.19.0
	go.opentelemetry.io/otel/metric v0.19.0
	go.opentelemetry.io/otel/sdk v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	gocloud.dev v0.23.0
	gocloud.dev/pubsub/kafkapubsub v0.23.0
	gocloud.dev/pubsub/natspubsub v0.23.0
	gocloud.dev/pubsub/rabbitpubsub v0.23.0
	gocloud.dev/secrets/hashivault v0.23.0
	golang.org/x/net v0.0.0-20210505214959-0714010a04ed
	golang.org/x/oauth2 v0.0.0-20210427180440-81ed05c6b58c
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/mikespook/gorbac.v2 v2.1.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	cloud.google.com/go v0.81.0 // indirect
	cloud.google.com/go/pubsub v1.10.3 // indirect
	cloud.google.com/go/storage v1.15.0 // indirect
	github.com/Azure/azure-amqp-common-go/v3 v3.1.0 // indirect
	github.com/Azure/azure-sdk-for-go v54.0.0+incompatible // indirect
	github.com/Azure/azure-service-bus-go v0.10.11 // indirect
	github.com/Azure/go-amqp v0.13.7 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.13 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Shopify/sarama v1.28.0 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blevesearch/bleve_index_api v1.0.0 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/mmap-go v1.0.2 // indirect
	github.com/blevesearch/scorch_segment_api/v2 v2.0.1 // indirect
	github.com/blevesearch/segment v0.9.0 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blevesearch/upsidedown_store_api v1.0.1 // indirect
	github.com/blevesearch/vellum v1.0.3 // indirect
	github.com/blevesearch/zapx/v11 v11.2.0 // indirect
	github.com/blevesearch/zapx/v12 v12.2.0 // indirect
	github.com/blevesearch/zapx/v13 v13.2.0 // indirect
	github.com/blevesearch/zapx/v14 v14.2.0 // indirect
	github.com/blevesearch/zapx/v15 v15.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/vault/sdk v0.2.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/steveyen/gtreap v0.1.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib v0.19.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.19.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.19.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	golang.org/x/tools v0.1.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.46.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210506142907-4a47615972c2 // indirect
	google.golang.org/grpc v1.37.0 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
