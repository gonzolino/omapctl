package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/gonzolino/omapctl/internal/ceph"
	"github.com/spf13/cobra"
)

func newSetCmd() *cobra.Command {
	var fromStdin bool

	cmd := &cobra.Command{
		Use:   "set <pool> <oid> <key> [value]",
		Short: "Set an omap key to a value",
		Long: `Set an omap key to a value. The value can be provided as a positional
argument or read from stdin with --stdin. Exactly one must be used.`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, oid, key := args[0], args[1], args[2]

			hasPositional := len(args) == 4
			if fromStdin && hasPositional {
				return usageError("cannot use both a positional value argument and --stdin")
			}
			if !fromStdin && !hasPositional {
				return usageError("provide a value argument or use --stdin")
			}

			var value []byte
			if fromStdin {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("read value from stdin: %w", err)
				}
				value = data
			} else {
				value = []byte(args[3])
			}

			conn := connFromContext(cmd.Context())
			ioctx, err := ceph.OpenIOContext(conn, pool)
			if err != nil {
				return err
			}
			defer ioctx.Destroy()

			if err := ioctx.SetOmap(oid, map[string][]byte{key: value}); err != nil {
				return fmt.Errorf("set omap key %q on %q in pool %q: %w", key, oid, pool, err)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&fromStdin, "stdin", false, "read value from stdin")

	return cmd
}
