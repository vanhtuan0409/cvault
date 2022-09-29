package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddPeekCommand(kmsClient *kms.Client, root *cobra.Command) {
	peekCmd := &cobra.Command{
		Use:               "peek",
		Short:             "Peek an encrypted file to stdout",
		ValidArgsFunction: completeStoreFile(),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyId := viper.GetString("keyId")
			if keyId == "" {
				return errors.New("invalid key id")
			}
			storeUrl := viper.GetString("store")
			if storeUrl == "" {
				return errors.New("invalid store url")
			}

			ctx := cmd.Context()
			s, err := storage.GetStorage(storeUrl)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Using KMS key: %s\n", keyId)
			fmt.Fprintf(os.Stderr, "Using storage: %s\n", storeUrl)
			fmt.Fprintln(os.Stderr, "--------------------------------")

			for _, fileKey := range args {
				fmt.Fprintf(os.Stderr, "# Source file: %s\n", fileKey)

				err := func() error {
					encrypted, err := s.Get(ctx, fileKey)
					if err != nil {
						return err
					}

					decrypted, err := cvault.Decrypt(ctx, kmsClient, keyId, encrypted)
					if err != nil {
						return err
					}

					os.Stdout.Write(decrypted)
					return nil
				}()
				if err != nil {
					fmt.Fprintf(os.Stderr, "[Warning] Unable to decrypt file %s. ERR: %v\n", fileKey, err)
				}
			}

			return nil
		},
	}

	root.AddCommand(peekCmd)
}
