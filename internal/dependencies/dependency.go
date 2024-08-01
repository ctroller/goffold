package dependencies

type Dependency struct {
	Pkg     string `yaml:"pkg"`
	Version string `yaml:"version"`
	Args    any    `yaml:"args"`
}
