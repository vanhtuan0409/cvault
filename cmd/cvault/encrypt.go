package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddEncryptCommand(root *cobra.Command) {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a file and push it into storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyId := viper.GetString("keyId")
			storeUrl := viper.GetString("store")

			ctx := cmd.Context()
			s, err := storage.GetStorage(storeUrl)
			if err != nil {
				return err
			}

			fmt.Printf("Using key: %s\n", keyId)
			fmt.Printf("Using storage: %s\n", storeUrl)
			fmt.Println("--------------------------------")

			for _, inputPath := range args {
				fmt.Printf("Source file: %s\n", inputPath)

				err := func() error {
					input, err := ioutil.ReadFile(inputPath)
					if err != nil {
						return err
					}

					encrypted, err := cvault.Encrypt(ctx, keyId, input)
					if err != nil {
						return err
					}

					fileKey := cvault.ToEncryptedName(filepath.Base(inputPath))
					return s.Put(ctx, fileKey, encrypted)
				}()
				if err != nil {
					fmt.Printf("[Warning] Unable to encrypt file %s. ERR: %v\n", inputPath, err)
				}
			}

			return nil
		},
	}

	root.AddCommand(encryptCmd)
}
