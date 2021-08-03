package envoy

import (
	"bufio"
	"context"
	"fmt"
	"regexp"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// Verify interface compliance
var _ common.Runnable = (*NoBadEnvoyLogsCheck)(nil)

// NoBadEnvoyLogsCheck implements common.Runnable
type NoBadEnvoyLogsCheck struct {
	client kubernetes.Interface
	pod    *v1.Pod
}

// HasNoBadEnvoyLogsCheck checks whether the envoy container of the pod has bad/error/warning log messages
func HasNoBadEnvoyLogsCheck(client kubernetes.Interface, pod *v1.Pod) NoBadEnvoyLogsCheck {
	return NoBadEnvoyLogsCheck{
		client: client,
		pod:    pod,
	}
}

// Info implements common.Runnable
func (check NoBadEnvoyLogsCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s has bad/error logs in envoy container", check.pod.Name)
}

// Run implements common.Runnable
func (check NoBadEnvoyLogsCheck) Run() error {
	envoyLogsTailLines := int64(10)
	envoyContainerName := "envoy"
	podLogsOpt := v1.PodLogOptions{
		Container: envoyContainerName,
		Follow:    false,
		Previous:  false,
		TailLines: &envoyLogsTailLines,
	}

	request := check.client.CoreV1().Pods(check.pod.Namespace).GetLogs(check.pod.Name, &podLogsOpt)
	envoyPodLogsReader, err := request.Stream(context.TODO())
	if err != nil {
		// If there are issues obtaining current container logs, return previously terminated container logs.
		podLogsOpt.Previous = true
		request := check.client.CoreV1().Pods(check.pod.Namespace).GetLogs(check.pod.Name, &podLogsOpt)
		envoyPodLogsReader, err = request.Stream(context.TODO())
		if err != nil {
			return fmt.Errorf("could not obtain %s container logs of pod %s: %#v", envoyContainerName, check.pod.Name, err)
		}
	}
	defer envoyPodLogsReader.Close() //nolint: errcheck,gosec

	re := regexp.MustCompile("(?i)(error)|(warn)")
	scanner := bufio.NewScanner(envoyPodLogsReader)
	var badEnvoyLogLines string
	for scanner.Scan() {
		logLine := scanner.Text()
		if re.MatchString(logLine) {
			badEnvoyLogLines += logLine + "\n"
		}
	}

	if len(badEnvoyLogLines) != 0 {
		log.Error().Msgf("%s container of pod %s contains bad logs", envoyContainerName, check.pod.Name)
		log.Error().Msg(badEnvoyLogLines)
	}

	return nil
}

// Suggestion implements common.Runnable.
func (check NoBadEnvoyLogsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check NoBadEnvoyLogsCheck) FixIt() error {
	panic("implement me")
}
