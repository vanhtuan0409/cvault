package main

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/storage"
)

func expandStorageWildcard(ctx context.Context, s storage.Storage, args []string) []string {
	if len(args) == 1 && args[0] == "*" {
		items, err := s.List(ctx)
		if err != nil {
			return []string{}
		}
		return cvault.SliceMap(items, func(i *storage.VaultItem) string {
			return i.Key
		})
	}
	return args
}

func completeStoreFile() func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		storeUrl := viper.GetString("store")
		if storeUrl == "" {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
		s, err := storage.GetStorage(storeUrl)
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		files, err := s.List(cmd.Context())
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}

		matches := cvault.SliceFilter(files, func(item *storage.VaultItem) bool {
			return strings.HasPrefix(item.Key, toComplete)
		})
		items := cvault.SliceMap(matches, func(item *storage.VaultItem) string {
			return item.Key
		})
		return items, cobra.ShellCompDirectiveNoFileComp
	}
}
