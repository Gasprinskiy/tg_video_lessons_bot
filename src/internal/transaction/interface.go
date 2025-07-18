package transaction

type Session interface {
	Start() error
	Rollback() error
	Commit() error
	Tx() interface{}
	TxIsActive() bool
	CreateNewSession() Session
	// context.Context NOT USED YET
	// WithContext(ctx context.Context) Session
	// TakeContext() context.Context
}

type SessionManager interface {
	CreateSession() Session
}
