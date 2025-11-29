/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [flags] [command]",
	Short: "Pass 1Password vault into a new process as environment variables",
	RunE:  handleRun,
}

func handleRun(cmd *cobra.Command, args []string) error {
	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return err
	}
	if token == "" {
		token = os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("token not specified")
	}
	vaults, err := cmd.Flags().GetStringSlice("vault")
	if err != nil {
		return err
	}
	env, err := GetEnv(vaults, token)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("command not specified")
	}
	process := exec.Command(args[0], args[1:]...)
	process.Env = os.Environ()
	for key, value := range env {
		process.Env = append(process.Env, fmt.Sprintf("%s=%s", key, value))
	}
	process.Stderr = os.Stderr
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	if err := process.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("token", "t", "", "1Password service account token. Can also be set via the OP_SERVICE_ACCOUNT_TOKEN environment variable.")
	runCmd.Flags().StringSliceP("vault", "v", nil, "Name or ID of a 1Password vault to export. Can be specified multiple times. Defaults to all accessible vaults.")
}
