module github.com/jcfug8/daylear/server

go 1.24

require (
	github.com/PuerkitoBio/goquery v1.10.3
	github.com/arran4/golang-ical v0.3.2
	github.com/aws/aws-sdk-go-v2 v1.32.8
	github.com/aws/aws-sdk-go-v2/config v1.28.10
	github.com/aws/aws-sdk-go-v2/credentials v1.17.51
	github.com/aws/aws-sdk-go-v2/service/s3 v1.72.2
	github.com/goccy/go-yaml v1.15.12
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1
	github.com/joho/godotenv v1.5.1
	github.com/justinas/alice v1.2.0
	github.com/lib/pq v1.10.9
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/rs/cors v1.11.1
	github.com/rs/zerolog v1.33.0
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.9.0
	github.com/teambition/rrule-go v1.8.2
	go.einride.tech/aip v0.68.1
	go.uber.org/fx v1.23.0
	golang.org/x/oauth2 v0.24.0
	google.golang.org/genai v1.15.0
	google.golang.org/genproto v0.0.0-20240903143218-8af14fe29dc1
	google.golang.org/genproto/googleapis/api v0.0.0-20241219192143-6b3ec007d9bb
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.0
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)

require (
	cloud.google.com/go v0.116.0 // indirect
	cloud.google.com/go/auth v0.9.3 // indirect
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.7 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.6 // indirect
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/rs/xid v1.5.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241219192143-6b3ec007d9bb // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
