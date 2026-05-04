package rollback

import "errors"

type RollBackManager struct {
	fns []func() error
}

func (m *RollBackManager) Add(fn func() error) {
	m.fns = append(m.fns, fn)
}

func (m *RollBackManager) Execute() error {
	var errs []error
	for i := len(m.fns) - 1; i >= 0; i-- {
		if err := m.fns[i](); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (m *RollBackManager) Clear() {
	m.fns = nil
}
func (m *RollBackManager) Join(other *RollBackManager) {
	m.fns = append(m.fns, other.fns...)
	// execution of rollback follows Last-In-First-Out order, so we append the other manager's functions to the end of the current manager's function slice.
}
