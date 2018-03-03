package cwl

type InputValue interface{}

type InputValues map[string]InputValue

type File struct {
	Location       string
	Path           string
	Basename       string
	Dirname        string
	Nameroot       string
	Nameext        string
	Checksum       string
	Size           int64
	Format         string
	Contents       string
	SecondaryFiles []Expression
}

type Directory struct {
	Location string
	Path     string
	Basename string
	Listing  []string
}
