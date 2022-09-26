package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddListCommand(client *s3.Client, root *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List encrypted file in store",
		RunE: func(cmd *cobra.Command, args []string) error {
			storeUrl := viper.GetString("store")
			if storeUrl == "" {
				return errors.New("invalid store url")
			}
			ctx := cmd.Context()

			s, err := storage.GetStorage(storeUrl, client)
			if err != nil {
				return err
			}

			files, err := s.List(ctx)
			if err != nil {
				return err
			}
			for _, fName := range files {
				fmt.Println(fName)
			}
			return nil
		},
	}

	root.AddCommand(listCmd)
}
