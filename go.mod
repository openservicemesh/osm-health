module github.com/openservicemesh/osm-health

go 1.16

require (
	github.com/axw/gocov v1.0.0 // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/envoyproxy/go-control-plane v0.9.9
	github.com/fatih/color v1.12.0
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/openservicemesh/osm v0.8.2-0.20210913092604-a59e75de0400
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/servicemeshinterface/smi-sdk-go v0.5.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	google.golang.org/protobuf v1.26.0
	helm.sh/helm/v3 v3.6.2
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/cli-runtime v0.21.2
	k8s.io/client-go v0.21.2
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
