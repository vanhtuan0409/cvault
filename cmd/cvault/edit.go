package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddEditCommand(root *cobra.Command) {
	editCmd := &cobra.Command{
		Use:               "edit",
		Short:             "Edit an encrypted file by temporary decrypt it",
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

			encrypted, err := s.Get(ctx, fileKey)
			if err != nil {
				return err
			}
			decrypted, err := cvault.Decrypt(ctx, keyId, encrypted)
			if err != nil {
				return err
			}

			tempFile, err := os.CreateTemp("", "")
			if err != nil {
				return err
			}
			defer tempFile.Close()
			defer os.Remove(tempFile.Name())

			if _, err := io.Copy(tempFile, bytes.NewBuffer(decrypted)); err != nil {
				return err
			}
			tempFile.Seek(0, io.SeekStart)

			editor := cvault.LookupEditor()
			if editor == "" {
				return errors.New("unknown editor var")
			}
			absEditor, err := exec.LookPath(editor)
			if err != nil {
				return err
			}

			editCmd := exec.CommandContext(ctx, absEditor, tempFile.Name())
			editCmd.Stdin = os.Stdin
			editCmd.Stdout = os.Stdout
			editCmd.Stderr = os.Stderr
			if err := editCmd.Run(); err != nil {
				return err
			}

			modified, err := ioutil.ReadAll(tempFile)
			if err != nil {
				return err
			}
			modifiedEncrypted, err := cvault.Encrypt(ctx, keyId, modified)
			if err != nil {
				return err
			}

			return s.Put(ctx, fileKey, modifiedEncrypted)
		},
	}

	root.AddCommand(editCmd)
}
