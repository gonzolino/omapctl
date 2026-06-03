package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/gonzolino/omapctl/internal/ceph"
	"github.com/gonzolino/omapctl/internal/output"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var prefix string
	var startAfter string
	var iteratorSize int64

	cmd := &cobra.Command{
		Use:   "list <pool> <oid>",
		Short: "List all omap key-value pairs for an object",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, oid := args[0], args[1]
			conn := connFromContext(cmd.Context())
			ioctx, err := ceph.OpenIOContext(conn, pool)
			if err != nil {
				return err
			}
			defer ioctx.Destroy()

			pairs, err := ioctx.GetAllOmapValues(oid, startAfter, prefix, iteratorSize)
			if err != nil {
				return fmt.Errorf("list omap values for %q in pool %q: %w", oid, pool, err)
			}

			keys := make([]string, 0, len(pairs))
			for k := range pairs {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			entries := make([]output.Entry, len(keys))
			for i, k := range keys {
				entries[i] = output.Entry{Key: k, Value: pairs[k]}
			}

			return output.PrintEntries(os.Stdout, entries, outputFormat())
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "filter keys by prefix")
	cmd.Flags().StringVar(&startAfter, "start-after", "", "list keys after this one")
	cmd.Flags().Int64Var(&iteratorSize, "iterator-size", 100, "internal iterator chunk size")

	return cmd
}
