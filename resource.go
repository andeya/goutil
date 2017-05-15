package srcpool

// Resource is a resource that can be stored in the Pool.
type Resource interface {
	// SetAvatar stores the contact with pool
	// Do not call it yourself, it is only called by (*Pool).get, and will only be called once
	SetAvatar(*Avatar)
	// GetAvatar gets the contact with pool
	// Do not call it yourself, it is only called by (*Pool).Put
	GetAvatar() *Avatar
	// Close closes the original source
	// No need to call it yourself, it is only called by (*Avatar).close
	Close() error
}
