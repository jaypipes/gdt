// Code generated by mockery v1.0.0. DO NOT EDIT.

// DO NOT EDIT MANUALLY. If you make changes to anything in interfaces.go, run make generate-mocks.

package mocks

import gdt "github.com/jaypipes/gdt"
import mock "github.com/stretchr/testify/mock"

// Parser is an autogenerated mock type for the Parser type
type Parser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: ca, contents
func (_m *Parser) Parse(ca gdt.ContextAppendable, contents []byte) error {
	ret := _m.Called(ca, contents)

	var r0 error
	if rf, ok := ret.Get(0).(func(gdt.ContextAppendable, []byte) error); ok {
		r0 = rf(ca, contents)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
