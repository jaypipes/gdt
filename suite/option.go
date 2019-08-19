package suite

// WithOption contains attributes describing a Suite
type WithOption struct {
	Name        string
	Description string
}

// WithName returns a WithOption that modifies a constructor for the Suite with
// a Name attribute
func WithName(name string) WithOption {
	return WithOption{Name: name}
}

// WithDescription returns a WithOption that modifies a constructor for the
// Suite with a Description attribute
func WithDescription(desc string) WithOption {
	return WithOption{Description: desc}
}

// mergeOptions takes zero or more WithOption structs and merges the values
// contained in those options into a single WithOption containing all non-zero
// values
func mergeOptions(opts ...WithOption) WithOption {
	res := WithOption{}
	for _, opt := range opts {
		if opt.Name != "" {
			res.Name = opt.Name
		}
		if opt.Description != "" {
			res.Description = opt.Description
		}
	}
	return res
}
