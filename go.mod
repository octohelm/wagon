module github.com/octohelm/wagon

go 1.20

replace github.com/docker/docker => github.com/docker/docker v20.10.3-0.20220414164044-61404de7df1a+incompatible

replace (
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.12.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.12.0
)

require (
	cuelang.org/go v0.5.0-beta.1
	dagger.io/dagger v0.4.5
	github.com/dagger/dagger v0.3.13
	github.com/davecgh/go-spew v1.1.1
	github.com/go-courier/logr v0.1.1
	github.com/google/uuid v1.3.0
	github.com/innoai-tech/infra v0.0.0-20230209064224-7873d371b0ba
	github.com/mattn/go-colorable v0.1.13
	github.com/moby/buildkit v0.11.1
	github.com/octohelm/courier v0.0.0-20230110034437-4f26125f8237
	github.com/octohelm/cuemod v0.6.2
	github.com/octohelm/gengo v0.0.0-20230203081555-b66e95cb6110
	github.com/octohelm/storage v0.0.0-20230103034019-3817a122e336
	github.com/octohelm/x v0.0.0-20220811034253-019992077a5d
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.9.3
	golang.org/x/exp v0.0.0-20230129154200-a960b3787bd2
	golang.org/x/mod v0.8.0
	golang.org/x/net v0.7.0
	golang.org/x/sync v0.1.0
)

require (
	github.com/Khan/genqlient v0.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/containerd v1.7.0-beta.4 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/containerd/typeurl/v2 v2.1.0 // indirect
	github.com/dagger/graphql v0.0.0-20221102000338-24d5e47d3b72 // indirect
	github.com/dagger/graphql-go-tools v0.0.0-20221102001222-e68b44170936 // indirect
	github.com/docker/cli v23.0.0+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v23.0.0+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/emicklei/proto v1.11.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-containerregistry v0.13.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/klauspost/compress v1.15.15 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/signal v0.7.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mpvl/unique v0.0.0-20150818121801-cbe035fff7de // indirect
	github.com/onsi/gomega v1.24.1 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2.0.20221005185240-3a7f492d3f1b // indirect
	github.com/protocolbuffers/txtpbfmt v0.0.0-20220926135727-61ed6f8e4d6e // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/cobra v1.6.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tonistiigi/fsutil v0.0.0-20230105215944-fb433841cbfa // indirect
	github.com/tonistiigi/units v0.0.0-20180711220420-6950e57a87ea // indirect
	github.com/tonistiigi/vt100 v0.0.0-20210615222946-8066bb97264f // indirect
	github.com/vektah/gqlparser/v2 v2.5.1 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.37.0 // indirect
	go.opentelemetry.io/otel v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.12.0 // indirect
	go.opentelemetry.io/otel/metric v0.34.0 // indirect
	go.opentelemetry.io/otel/sdk v1.13.0 // indirect
	go.opentelemetry.io/otel/trace v1.13.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/tools v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20230202175211-008b39050e57 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.25.4 // indirect
	k8s.io/apimachinery v0.25.4 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
