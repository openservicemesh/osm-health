package envoy

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm/pkg/utils"
)

// Verify interface compliance
var _ common.Runnable = (*HasValidEnvoyCertificateCheck)(nil)

// HasValidEnvoyCertificateCheck implements common.Runnable
type HasValidEnvoyCertificateCheck struct {
	ConfigGetter
	pod             *corev1.Pod
	k8s             kubernetes.Interface
	certificateType SDSCertType
}

// SDSCertType is a type of a certificate requested by an Envoy proxy via SDS.
type SDSCertType string

const (
	// ServiceCertType is the prefix for the service certificate resource name.
	// Example: "service-cert:<service namespace>/<service service account>"
	ServiceCertType SDSCertType = "service-cert"

	// RootCertTypeForMTLSOutbound is the prefix for the mTLS root certificate
	// resource name for upstream connectivity.
	// Example: "root-cert-for-mtls-outbound:<service namespace>/<service name>"
	RootCertTypeForMTLSOutbound SDSCertType = "root-cert-for-mtls-outbound"

	// RootCertTypeForMTLSInbound is the prefix for the mTLS root certificate
	// resource name for downstream connectivity.
	// Example: "root-cert-for-mtls-inbound:<service namespace>/<service service account>"
	RootCertTypeForMTLSInbound SDSCertType = "root-cert-for-mtls-inbound"
)

func (ct SDSCertType) String() string {
	return string(ct)
}

// Run implements common.Runnable
func (c HasValidEnvoyCertificateCheck) Run() outcomes.Outcome {
	if c.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return outcomes.FailedOutcome{Error: ErrIncorrectlyInitializedConfigGetter}
	}
	envoyConfig, err := c.ConfigGetter.GetConfig()
	if err != nil {
		return outcomes.FailedOutcome{Error: err}
	}

	if envoyConfig == nil {
		return outcomes.FailedOutcome{Error: ErrEnvoyConfigEmpty}
	}

	// Checks if a secret of the specified certificate type exists in the
	// provided Envoy config with a name derived from the given Pod.
	// Secret name formats for each certificate type:
	// - root-cert-for-mtls-outbound:<service namespace>/<service name>
	// - root-cert-for-mtls-inbound:<service namespace>/<service service account>
	// - service-cert:<service namespace>/<service service account>
	//
	// The secret name for an outbound mTLS root certificate is formatted using
	// the destination Pod's service name. The Pod might back multiple services,
	// so construct the secret name for each possible service.
	possibleSecretNames := map[string]struct{}{}
	if c.certificateType == RootCertTypeForMTLSOutbound {
		svcs, err := kuberneteshelper.GetMatchingServices(c.k8s, c.pod.Labels, c.pod.Namespace)
		if err != nil {
			return outcomes.FailedOutcome{Error: errors.Wrapf(err, "failed to map Pod %s/%s to Kubernetes Services", c.pod.Namespace, c.pod.Name)}
		}
		for _, svc := range svcs {
			possibleSecretNames[fmt.Sprintf("%s:%s", c.certificateType.String(), utils.K8sSvcToMeshSvc(svc).String())] = struct{}{}
		}
	} else {
		possibleSecretNames[fmt.Sprintf("%s:%s/%s", c.certificateType.String(), c.pod.Namespace, c.pod.Spec.ServiceAccountName)] = struct{}{}
	}

	if len(possibleSecretNames) == 0 {
		return outcomes.FailedOutcome{Error: fmt.Errorf("no secrets listed in the Envoy config")}
	}

	// Check that at least one of the possible secret names is in the
	// Envoy config
	found := false
	var foundSecretNames []string
	for _, dynSecret := range envoyConfig.SecretsConfigDump.GetDynamicActiveSecrets() {
		foundSecretNames = append(foundSecretNames, dynSecret.Name)
		if _, exists := possibleSecretNames[dynSecret.Name]; exists {
			found = true
			break
		}
	}

	if !found {
		return outcomes.FailedOutcome{Error: fmt.Errorf("expected a secret named one of %v, but only found %v", possibleSecretNames, foundSecretNames)}
	}

	return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
}

// Suggestion implements common.Runnable
func (c HasValidEnvoyCertificateCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (c HasValidEnvoyCertificateCheck) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (c HasValidEnvoyCertificateCheck) Description() string {
	return fmt.Sprintf("Checking whether %s is configured with a %s envoy secret", c.ConfigGetter.GetObjectName(), c.certificateType)
}

// HasInboundRootCertificate creates a new common.Runnable, which checks whether the given Pod
// has an Envoy with a properly configured inbound root validation certificate.
func HasInboundRootCertificate(client kubernetes.Interface, dstConfigGetter ConfigGetter, dst *corev1.Pod) common.Runnable {
	return HasValidEnvoyCertificateCheck{
		ConfigGetter:    dstConfigGetter,
		k8s:             client,
		pod:             dst,
		certificateType: RootCertTypeForMTLSInbound,
	}
}

// HasOutboundRootCertificate creates a new common.Runnable, which checks whether the source Pod
// has an Envoy with a properly configured outbound root validation certificate for the given
// destination Pod.
func HasOutboundRootCertificate(client kubernetes.Interface, srcConfigGetter ConfigGetter, dst *corev1.Pod) common.Runnable {
	return HasValidEnvoyCertificateCheck{
		ConfigGetter:    srcConfigGetter,
		k8s:             client,
		pod:             dst,
		certificateType: RootCertTypeForMTLSOutbound,
	}
}

// HasServiceCertificate creates a new common.Runnable, which checks whether the given Pod
// has an Envoy with a properly configured service certificate.
func HasServiceCertificate(client kubernetes.Interface, configGetter ConfigGetter, pod *corev1.Pod) common.Runnable {
	return HasValidEnvoyCertificateCheck{
		ConfigGetter:    configGetter,
		k8s:             client,
		pod:             pod,
		certificateType: ServiceCertType,
	}
}
