package interfaces

type Named interface {
	// Name returns a string name
	Name() string
	// Describe returns a longer string description
	Describe() string
}
