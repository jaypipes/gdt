package interfaces

type Fixture interface {
	Start()
	Stop()
	HasState(string) bool
	State(string) string
}

// FixtureRegistry describes something that can register and return fixtures
type FixtureRegistry interface {
	Register(string, Fixture)
	Get(string) Fixture
	List() []Fixture
}
