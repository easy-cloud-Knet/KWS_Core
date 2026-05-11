package creation

import (
	"fmt"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func BootExistingVM(d BootableDomain) error {
	if err := d.Create(); err != nil {
		return virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("error starting domain: %w", err))
	}
	return nil
}
