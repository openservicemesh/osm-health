package podhelper

import (
	"bufio"
	"context"
	"fmt"
	"regexp"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// HasNoBadLogs checks whether the logs of the pod container contain bad (fatal/error/warning/fail) logs
func HasNoBadLogs(client kubernetes.Interface, pod *v1.Pod, containerName string) error {
	logsTailLines := int64(10)
	podLogsOpt := v1.PodLogOptions{
		Container: containerName,
		Follow:    false,
		Previous:  false,
		TailLines: &logsTailLines,
	}

	request := client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogsOpt)
	podLogsReader, err := request.Stream(context.TODO())
	if err != nil {
		// If there are issues obtaining current container logs, return previously terminated container logs.
		podLogsOpt.Previous = true
		request := client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogsOpt)
		podLogsReader, err = request.Stream(context.TODO())
		if err != nil {
			return fmt.Errorf("could not obtain %s container logs of pod %s: %#v", containerName, pod.Name, err)
		}
	}
	defer podLogsReader.Close() //nolint: errcheck,gosec

	re := regexp.MustCompile("(?i)(fatal)|(error)|(warn)|(fail)")
	scanner := bufio.NewScanner(podLogsReader)
	var badLogLines string
	for scanner.Scan() {
		logLine := scanner.Text()
		if re.MatchString(logLine) {
			badLogLines += logLine + "\n"
		}
	}

	if len(badLogLines) != 0 {
		log.Error().Msgf("%s container of pod %s contains bad logs", containerName, pod.Name)
		log.Error().Msg(badLogLines)
	}

	return nil
}
