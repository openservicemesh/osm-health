module github.com/openservicemesh/osm-health

go 1.16

require (
	github.com/openservicemesh/osm v0.9.1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/servicemeshinterface/smi-sdk-go v0.5.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
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
