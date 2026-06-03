package ceph

import (
	"fmt"

	"github.com/ceph/go-ceph/rados"
)

// NewConn creates, configures, and connects a *rados.Conn.
// If configFile is empty, the system default (/etc/ceph/ceph.conf) is used.
func NewConn(configFile string) (*rados.Conn, error) {
	conn, err := rados.NewConn()
	if err != nil {
		return nil, fmt.Errorf("create rados connection: %w", err)
	}
	if configFile == "" {
		if err := conn.ReadDefaultConfigFile(); err != nil {
			return nil, fmt.Errorf("read default ceph config: %w", err)
		}
	} else {
		if err := conn.ReadConfigFile(configFile); err != nil {
			return nil, fmt.Errorf("read ceph config %q: %w", configFile, err)
		}
	}
	if err := conn.Connect(); err != nil {
		return nil, fmt.Errorf("connect to ceph cluster: %w", err)
	}
	return conn, nil
}

// OpenIOContext opens an IOContext for the given pool.
// The caller must call ioctx.Destroy() when done.
func OpenIOContext(conn *rados.Conn, pool string) (*rados.IOContext, error) {
	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return nil, fmt.Errorf("open io context for pool %q: %w", pool, err)
	}
	return ioctx, nil
}
