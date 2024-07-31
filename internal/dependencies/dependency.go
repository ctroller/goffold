package dependencies

type Dependency struct {
	Pkg     string `yaml:"pkg"`
	Version string `yaml:"version" default:"latest"`
	Args    any    `yaml:"args"`
}
