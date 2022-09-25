package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddListCommand(client *s3.Client, root *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List encrypted file in store",
		RunE: func(cmd *cobra.Command, args []string) error {
			storeUrl := cmd.Flag("store").Value.String()
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

			w := tabwriter.NewWriter(os.Stdout, 4, 2, 1, '\t', 0)
			defer w.Flush()
			for _, fName := range files {
				fmt.Fprintf(w, "%s\t", fName)
			}
			fmt.Fprintln(w)

			return nil
		},
	}

	root.AddCommand(listCmd)
}
