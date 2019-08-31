package interfaces

type Fixture interface {
	Start()
	Stop()
	State(string) string
}

// FixtureRegistry describes something that can register and return fixtures
type FixtureRegistry interface {
	Register(string, Fixture)
	Get(string) Fixture
}
