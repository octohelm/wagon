module github.com/octohelm/wagon

go 1.20

replace (
	// follow from dagger
	cloud.google.com/go => cloud.google.com/go v0.100.2
	cuelang.org/go => cuelang.org/go v0.5.0-beta.1
	github.com/moby/buildkit => github.com/moby/buildkit v0.11.0-rc3.0.20230529232739-baffc1bda21b
	github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.1.0-rc3
	github.com/spdx/tools-golang => github.com/spdx/tools-golang v0.5.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.40.0
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/exporters/jaeger => go.opentelemetry.io/otel/exporters/jaeger v1.14.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.14.0
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v0.37.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.14.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.14.0
	go.opentelemetry.io/proto/otlp => go.opentelemetry.io/proto/otlp v0.19.0
)

require (
	cuelang.org/go v0.5.0-beta.1
	dagger.io/dagger v0.7.2
	github.com/dagger/dagger v0.6.2
	github.com/davecgh/go-spew v1.1.1
	github.com/go-courier/logr v0.2.0
	github.com/google/uuid v1.3.0
	github.com/innoai-tech/infra v0.0.0-20230615040720-354f8c1272d4
	github.com/mattn/go-colorable v0.1.13
	github.com/octohelm/courier v0.0.0-20230612084835-5bae73959f5b
	github.com/octohelm/cuemod v0.6.3
	github.com/octohelm/gengo v0.0.0-20230602052920-d1e8eaa72959
	github.com/octohelm/storage v0.0.0-20230511032824-7f3085d1c289
	github.com/octohelm/x v0.0.0-20220811034253-019992077a5d
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.9.5
	github.com/vito/progrock v0.7.0
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
	golang.org/x/mod v0.12.0
	golang.org/x/net v0.11.0
	golang.org/x/sync v0.3.0
)

require (
	github.com/99designs/gqlgen v0.17.33 // indirect
	github.com/Khan/genqlient v0.6.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230528122434-6f98819771a1 // indirect
	github.com/acomagu/bufpipe v1.0.4 // indirect
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/charmbracelet/bubbles v0.16.1 // indirect
	github.com/charmbracelet/bubbletea v0.24.2 // indirect
	github.com/charmbracelet/harmonica v0.2.0 // indirect
	github.com/charmbracelet/lipgloss v0.7.1 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/containerd/containerd v1.7.2 // indirect
	github.com/containerd/continuity v0.4.1 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.14.3 // indirect
	github.com/containerd/typeurl/v2 v2.1.1 // indirect
	github.com/dagger/graphql v0.0.0-20230601100125-137fc3a90735 // indirect
	github.com/dagger/graphql-go-tools v0.0.0-20230418214324-32c52f390881 // indirect
	github.com/docker/cli v24.0.2+incompatible // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.2+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/emicklei/proto v1.11.2 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fogleman/ease v0.0.0-20170301025033-8da417bf1776 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.4.1 // indirect
	github.com/go-git/go-git/v5 v5.7.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-containerregistry v0.15.2 // indirect
	github.com/google/go-github/v50 v50.2.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.16.6 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/buildkit v0.11.0-rc3.0.20230608232644-8a28fe6bc051 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/signal v0.7.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mpvl/unique v0.0.0-20150818121801-cbe035fff7de // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.1 // indirect
	github.com/onsi/gomega v1.27.6 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc3 // indirect
	github.com/opencontainers/runc v1.1.7 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/protocolbuffers/txtpbfmt v0.0.0-20230412060525-fa9f017c0ded // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skeema/knownhosts v1.1.1 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tonistiigi/fsutil v0.0.0-20230407161946-9e7a6df48576 // indirect
	github.com/tonistiigi/units v0.0.0-20180711220420-6950e57a87ea // indirect
	github.com/tonistiigi/vt100 v0.0.0-20210615222946-8066bb97264f // indirect
	github.com/vbatts/tar-split v0.11.3 // indirect
	github.com/vektah/gqlparser/v2 v2.5.3 // indirect
	github.com/vito/vt100 v0.1.2 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	github.com/zmb3/spotify/v2 v2.3.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.40.0 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.14.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/oauth2 v0.9.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/term v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.56.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.27.3 // indirect
	k8s.io/apimachinery v0.27.3 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/utils v0.0.0-20230505201702-9f6742963106 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
