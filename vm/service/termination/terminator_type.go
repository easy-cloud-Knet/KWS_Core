package termination

type DomainDeleteType uint

const (
	HardDelete DomainDeleteType = iota
	SoftDelete
)

type DomainDeletion interface {
	DeleteDomain() error
}

type DomainTermination interface {
	TerminateDomain() error
}

type DomainTerminator struct {
	domain Domain
}
type DomainDeleter struct {
	uuid         string
	domain       Domain
	DeletionType DomainDeleteType
}
