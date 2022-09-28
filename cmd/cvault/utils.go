package main

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/storage"
)

func completeStoreFile(s3Client *s3.Client) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		storeUrl := viper.GetString("store")
		if storeUrl == "" {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
		s, err := storage.GetStorage(storeUrl, s3Client)
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
