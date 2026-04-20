package libvirtconn

import (
	"os"

	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func Connect(logger *zap.Logger) (*libvirt.Connect, error) {
	conn, err := libvirt.NewConnect(config.LibvirtURI)
	if err != nil {
		return nil, err
	}
	logger.Info("Libvirt connection successfully done.", zap.Int("pid", os.Getegid()))
	defer logger.Sync()
	return conn, nil
}

func IsAlive(conn *libvirt.Connect) bool {
	if conn == nil {
		return false
	}
	alive, err := conn.IsAlive()
	return err == nil && alive
}
