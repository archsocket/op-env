package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/1password/onepassword-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "op-env",
	Short: "Convert 1Password vaults into a dotenv file",
	Long: `op-env is a CLI tool that converts items from one or more 1Password vaults into a .env (dotenv) file.
It simplifies integrating 1Password with application deployments by exporting vault contents as environment variables.`,
	RunE: RunE,
}

func FormatKey(s string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		return "", fmt.Errorf("regexp.Compile: %w", err)
	}
	key := strings.ToUpper(reg.ReplaceAllString(strings.ReplaceAll(s, " ", "_"), ""))
	if len(key) == 0 {
		return "", fmt.Errorf("resulting key length is zero")
	}
	return key, nil
}

func RunE(cmd *cobra.Command, args []string) error {
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
	client, err := onepassword.NewClient(
		context.TODO(),
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("OP-ENV", "v1.1.0"),
	)
	if err != nil {
		return err
	}
	vaultOverviews, err := client.Vaults().List(context.TODO())
	if err != nil {
		return err
	}

	if len(vaults) == 0 {
		vaults = make([]string, len(vaultOverviews))
		for i, vaultOverview := range vaultOverviews {
			vaults[i] = vaultOverview.ID
		}
	} else {
		for i, vault := range vaults {
			for _, vaultOverview := range vaultOverviews {
				if vault == vaultOverview.Title {
					vaults[i] = vaultOverview.ID
				}
			}
		}
	}
	env := make(map[string]string)
	for _, vault := range vaults {
		items, err := client.Items().List(context.TODO(), vault)
		if err != nil {
			return fmt.Errorf("client.items.list '%s': %w", vault, err)
		}
		for _, itemSummary := range items {
			item, err := client.Items().Get(context.TODO(), vault, itemSummary.ID)
			if err != nil {
				return fmt.Errorf("client.items.get: %w", err)
			}
			key, err := FormatKey(item.Title)
			if err != nil {
				return fmt.Errorf("format key: %w", err)
			}
			if len(item.Notes) > 0 {
				env[key] = item.Notes
			}
			for _, field := range item.Fields {
				fieldKey, err := FormatKey(field.Title)
				if err != nil {
					return fmt.Errorf("format key: %w", err)
				}
				fieldKey = key + "_" + fieldKey
				env[fieldKey] = field.Value
			}
		}
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
