package termination

type Domain interface {
	IsActive() (bool, error)
	Destroy() error
	Undefine() error
}
