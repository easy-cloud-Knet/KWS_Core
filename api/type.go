package api

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type InstHandler struct {
	LibvirtInst   *libvirt.Connect
	DomainControl *domCon.DomListControl
	Logger        *zap.Logger
}
