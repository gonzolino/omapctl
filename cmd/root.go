package cmd

import (
	"context"
	"fmt"
	"os"

	radoslib "github.com/ceph/go-ceph/rados"
	"github.com/gonzolino/omapctl/internal/ceph"
	"github.com/gonzolino/omapctl/internal/output"
	"github.com/spf13/cobra"
)

type contextKey int

const connKey contextKey = iota

type globalFlags struct {
	configFile string
	outputFmt  string
}

var rootFlags globalFlags

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "omapctl",
		Short: "Inspect and mutate Ceph RADOS omap stores",
		Long: `omapctl is a CLI tool for CRUD operations on Ceph RADOS omap stores.
Useful for debugging by introspecting and modifying omap key-value pairs on objects.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := output.ParseFormat(rootFlags.outputFmt); err != nil {
				return err
			}
			conn, err := ceph.NewConn(rootFlags.configFile)
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(cmd.Context(), connKey, conn))
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if conn := connFromContext(cmd.Context()); conn != nil {
				conn.Shutdown()
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&rootFlags.configFile, "config", "c", "", "path to ceph.conf (default: system default)")
	cmd.PersistentFlags().StringVarP(&rootFlags.outputFmt, "output", "o", "table", "output format: table|json")

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newSetCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newClearCmd())

	return cmd
}

func connFromContext(ctx context.Context) *radoslib.Conn {
	if ctx == nil {
		return nil
	}
	conn, _ := ctx.Value(connKey).(*radoslib.Conn)
	return conn
}

func outputFormat() output.Format {
	f, _ := output.ParseFormat(rootFlags.outputFmt)
	return f
}

func Execute() error {
	cmd := newRootCmd()
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	return cmd.Execute()
}

func usageError(msg string) error {
	return fmt.Errorf("%s", msg)
}
