package status

import (
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"github.com/easy-cloud-Knet/KWS_Core/pkg/virtXML/fetcher"
	"go.uber.org/zap"
)

type XMLStatus struct{}

func (ds *XMLStatus) RetrieveStatus(dom Domain, sources map[vmtypes.SourceType]int, _ *zap.Logger) (map[vmtypes.SourceType]int, error) {
	return fetcher.NewXMLFetcher().Fetch(dom, sources)
}
