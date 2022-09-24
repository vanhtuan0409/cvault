package main

import (
	"errors"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/spf13/cobra"
	"github.com/vanhtuan0409/cvault"
)

func AddDecryptCommand(client *kms.Client, root *cobra.Command) {
	decryptCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt a file from storage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyId := cmd.Flag("key-id").Value.String()
			if keyId == "" {
				return errors.New("invalid key id")
			}

			input, err := ioutil.ReadFile("encrypted")
			if err != nil {
				return err
			}

			decrypted, err := cvault.Decrypt(cmd.Context(), client, keyId, input)
			if err != nil {
				return err
			}

			return ioutil.WriteFile("decrypted", decrypted, 0644)
		},
	}
	decryptCmd.PersistentFlags().StringP("output", "o", "", "Output file")

	root.AddCommand(decryptCmd)
}
