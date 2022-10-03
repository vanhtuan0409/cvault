package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddDecryptCommand(root *cobra.Command) {
	decryptCmd := &cobra.Command{
		Use:               "decrypt",
		Short:             "Decrypt a file from storage",
		ValidArgsFunction: completeStoreFile(),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyId := viper.GetString("keyId")
			storeUrl := viper.GetString("store")

			outputDir := cmd.Flag("output-dir").Value.String()
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				if err := os.MkdirAll(outputDir, 0644); err != nil {
					return err
				}
			}

			ctx := cmd.Context()
			s, err := storage.GetStorage(storeUrl)
			if err != nil {
				return err
			}

			fmt.Printf("Using key: %s\n", keyId)
			fmt.Printf("Using storage: %s\n", storeUrl)
			fmt.Printf("Output dir: %s\n", outputDir)
			fmt.Println("--------------------------------")

			for _, fileKey := range args {
				fileName := cvault.ToDecryptedName(fileKey)
				decryptedFile := filepath.Join(outputDir, fileName)

				fmt.Printf("Source file: %s\n", fileKey)
				fmt.Printf("Decrypted file: %s\n", decryptedFile)

				err := func() error {
					encrypted, err := s.Get(ctx, fileKey)
					if err != nil {
						return err
					}

					decrypted, err := cvault.Decrypt(ctx, keyId, encrypted)
					if err != nil {
						return err
					}

					return ioutil.WriteFile(decryptedFile, decrypted, 0644)
				}()
				if err != nil {
					fmt.Printf("[Warning] Unable to decrypt file %s. ERR: %v\n", fileKey, err)
				}
			}

			return nil
		},
	}
	decryptCmd.PersistentFlags().StringP("output-dir", "o", "/tmp", "Output directory")

	root.AddCommand(decryptCmd)
}
