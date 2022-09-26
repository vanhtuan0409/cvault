package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddRemoveCommand(client *s3.Client, root *cobra.Command) {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove encrypted file from store",
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

			for _, fName := range args {
				if err := s.Remove(ctx, fName); err != nil {
					fmt.Printf("[Warning] Unable to remove file %s. ERR: %v\n", fName, err)
				}
			}

			return nil
		},
	}

	root.AddCommand(removeCmd)
}
