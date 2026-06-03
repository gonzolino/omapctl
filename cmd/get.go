package cmd

import (
	"fmt"
	"os"

	"github.com/gonzolino/omapctl/internal/ceph"
	"github.com/gonzolino/omapctl/internal/output"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <pool> <oid> <key>",
		Short: "Get a single omap value by key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, oid, key := args[0], args[1], args[2]
			conn := connFromContext(cmd.Context())
			ioctx, err := ceph.OpenIOContext(conn, pool)
			if err != nil {
				return err
			}
			defer ioctx.Destroy()

			// Use the key as prefix with maxReturn=2; then verify exact match.
			pairs, err := ioctx.GetOmapValues(oid, "", key, 2)
			if err != nil {
				return fmt.Errorf("get omap value %q on %q in pool %q: %w", key, oid, pool, err)
			}
			val, ok := pairs[key]
			if !ok {
				return fmt.Errorf("key %q not found on object %q in pool %q", key, oid, pool)
			}

			return output.PrintEntry(os.Stdout, output.Entry{Key: key, Value: val}, outputFormat())
		},
	}
}
