package main

import (
	"errors"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/spf13/cobra"
	"github.com/vanhtuan0409/cvault"
)

func AddEncryptCommand(client *kms.Client, root *cobra.Command) {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a file and push it into storage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyId := cmd.Flag("key-id").Value.String()
			if keyId == "" {
				return errors.New("invalid key id")
			}
			inputPath := args[0]

			input, err := ioutil.ReadFile(inputPath)
			if err != nil {
				return err
			}

			encrypted, err := cvault.Encrypt(cmd.Context(), client, keyId, input)
			if err != nil {
				return err
			}

			return ioutil.WriteFile("encrypted", encrypted, 0644)
		},
	}

	root.AddCommand(encryptCmd)
}
