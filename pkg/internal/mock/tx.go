package mock

type Tx struct{}

func (m *Tx) Commit() error { return nil }

func (m *Tx) Rollback() error { return nil }
