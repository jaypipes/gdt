package interfaces

type Fixture interface {
	Start()
	Cleanup()
}
