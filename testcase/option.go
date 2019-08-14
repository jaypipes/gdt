package testcase

type WithOption struct {
	Name        string
	Description string
	Filepath    string
}

func WithName(name string) WithOption {
	return WithOption{Name: name}
}

func WithDescription(desc string) WithOption {
	return WithOption{Description: desc}
}

func WithFilepath(fp string) WithOption {
	return WithOption{Filepath: fp}
}

func MergeOptions(opts ...WithOption) WithOption {
	res := WithOption{}
	for opt := range opts {
		if opt.Name != "" {
			res.Name = opt.Name
		}
		if opt.Description != "" {
			res.Description = opt.Description
		}
		if opt.Filepath != "" {
			res.Filepath = opt.Filepath
		}
	}
	return res
}
