package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "op-env",
	Short: "Convert 1Password vaults into a dotenv file",
	Long: `op-env is a CLI tool that converts items from one or more 1Password vaults into a .env (dotenv) file, or injects items directly into a new process as environment variables.
It simplifies integrating 1Password with application deployments by exporting vault contents as environment variables.`,
	RunE: handleRoot,
}

func handleRoot(cmd *cobra.Command, args []string) error {
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
	out, err := cmd.Flags().GetString("out")
	if err != nil {
		return err
	}
	if out == "" {
		return fmt.Errorf("out not specified")
	}
	env, err := GetEnv(vaults, token)
	if err != nil {
		return err
	}
	var writer io.Writer
	if out == "-" {
		writer = os.Stdout
	} else {
		fileWriter, err := os.Create(out)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := fileWriter.Close(); cerr != nil && err == nil {
				err = cerr
			}
		}()
		writer = fileWriter
	}
	for key, value := range env {
		if _, err := fmt.Fprintf(writer, "%s=\"%s\"\n", key, strings.ReplaceAll(strings.ReplaceAll(value, "\n", "\\n"), "\"", "\\\"")); err != nil {
			return err
		}
	}
	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Normalize(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "file":
		name = "out"
	}
	return pflag.NormalizedName(name)
}

func init() {
	rootCmd.Flags().SetNormalizeFunc(Normalize)
	rootCmd.Flags().StringP("token", "t", "", "1Password service account token. Can also be set via the OP_SERVICE_ACCOUNT_TOKEN environment variable.")
	rootCmd.Flags().StringSliceP("vault", "v", nil, "Name or ID of a 1Password vault to export. Can be specified multiple times. Defaults to all accessible vaults.")
	rootCmd.Flags().StringP("out", "o", ".env", "Output to a file. Use \"-\" to write to stdout.")
}
