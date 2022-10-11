package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddPeekCommand(root *cobra.Command) {
	peekCmd := &cobra.Command{
		Use:               "peek",
		Short:             "Peek an encrypted file to stdout",
		ValidArgsFunction: completeStoreFile(),
		Args:              cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyId := viper.GetString("keyId")
			storeUrl := viper.GetString("store")
			fileKey := args[0]

			ctx := cmd.Context()
			s, err := storage.GetStorage(storeUrl)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Using key: %s\n", keyId)
			fmt.Fprintf(os.Stderr, "Using storage: %s\n", storeUrl)
			fmt.Fprintf(os.Stderr, "Source file: %s\n", fileKey)
			fmt.Fprintln(os.Stderr, "--------------------------------")

			encrypted, err := s.Get(ctx, fileKey)
			if err != nil {
				return err
			}

			decrypted, err := cvault.Decrypt(ctx, keyId, encrypted)
			if err != nil {
				return err
			}

			os.Stdout.Write(decrypted)
			return nil
		},
	}

	root.AddCommand(peekCmd)
}
