package testcase

// WithOption contains attributes describing a Testcase
type WithOption struct {
	Name        string
	Description string
}

// WithName returns a WithOption that modifies a constructor for the Testcase
// with a Name attribute
func WithName(name string) WithOption {
	return WithOption{Name: name}
}

// WithDescription returns a WithOption that modifies a constructor for the
// Testcase with a Description attribute
func WithDescription(desc string) WithOption {
	return WithOption{Description: desc}
}

// mergeOptions takes zero or more WithOption structs and merges the values
// contained in those options into a single WithOption containing all non-zero
// values
func mergeOptions(opts ...WithOption) WithOption {
	res := WithOption{}
	for opt := range opts {
		if opt.Name != "" {
			res.Name = opt.Name
		}
		if opt.Description != "" {
			res.Description = opt.Description
		}
	}
	return res
}
