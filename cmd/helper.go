package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/1password/onepassword-sdk-go"
)

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

func GetEnv(vaults []string, token string) (map[string]string, error) {
	client, err := onepassword.NewClient(
		context.TODO(),
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("OP-ENV", "v1.1.0"),
	)
	if err != nil {
		return nil, err
	}
	vaultOverviews, err := client.Vaults().List(context.TODO())
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("client.items.list '%s': %w", vault, err)
		}
		for _, itemSummary := range items {
			item, err := client.Items().Get(context.TODO(), vault, itemSummary.ID)
			if err != nil {
				return nil, fmt.Errorf("client.items.get: %w", err)
			}
			key, err := FormatKey(item.Title)
			if err != nil {
				return nil, fmt.Errorf("format key: %w", err)
			}
			if len(item.Notes) > 0 {
				env[key] = item.Notes
			}
			for _, field := range item.Fields {
				fieldKey, err := FormatKey(field.Title)
				if err != nil {
					return nil, fmt.Errorf("format key: %w", err)
				}
				fieldKey = key + "_" + fieldKey
				env[fieldKey] = field.Value
			}
		}
	}
	return env, nil
}
