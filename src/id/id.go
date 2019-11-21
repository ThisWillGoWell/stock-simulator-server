package id

/**
Identifiable are used throughout the design, its just something
that can be Identified by type and uuid
*/
type Identifiable interface {
	GetId() string
	GetType() string
}
