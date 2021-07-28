package main

import "github.com/spf13/cobra"

func newIngressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "Checks ingress to Kubernetes resources",
		Long:  `Checks ingressKubernetes resources`,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newIngressToPodCmd())
	return cmd
}
