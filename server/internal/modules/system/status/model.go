package status

type Info struct {
	Status   string
	Checks   Checks
	Warnings []string
}

type Checks struct {
	Database Dependency
}

type Dependency struct {
	Configured bool
	OK         bool
	Message    string
}
