package interfaces

import "testing"

type Testable interface {
	T() *testing.T
	RunWithFixtures(FixtureRegistry)
}
