module github.com/Vernacular-ai/vcore

go 1.12

require (
	github.com/Vernacular-ai/gorm v1.10.1
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/getsentry/sentry-go v0.12.0
	github.com/google/go-cmp v0.5.5
	github.com/hashicorp/go-getter v1.5.6
	github.com/hashicorp/vault/api v1.5.0 // indirect
	github.com/hashicorp/vault/api/auth/approle v0.1.1 // indirect
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mediocregopher/radix.v2 v0.0.0-20181115013041-b67df6e626f9
	github.com/mediocregopher/radix/v3 v3.4.2
	github.com/newrelic/go-agent v2.10.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/streadway/amqp v0.0.0-20190815230801-eade30b20f1d
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	google.golang.org/genproto v0.0.0-20210729151513-df9385d47c1b // indirect
	google.golang.org/grpc v1.41.0
	gopkg.in/yaml.v2 v2.2.5
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.2

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
