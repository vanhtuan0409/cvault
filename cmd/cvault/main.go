package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "cvault",
		Short:        "A cold vault encryption tool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	rootCmd.PersistentFlags().StringP("key-id", "k", "", "KMS key id")
	rootCmd.PersistentFlags().StringP("store", "s", "local://.", "Location of storage")

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	kmsClient := kms.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	AddEncryptCommand(kmsClient, s3Client, rootCmd)
	AddDecryptCommand(kmsClient, s3Client, rootCmd)
	AddListCommand(s3Client, rootCmd)
	AddRemoveCommand(s3Client, rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
