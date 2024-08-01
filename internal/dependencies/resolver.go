package dependencies

type DependencyResolver struct {
	Type     string
	Resolve  func(dir string, dep Dependency) ([]byte, error)
	Finisher func(dir string) error
}

var resolvers = map[string]DependencyResolver{}

func RegisterResolver(resolver DependencyResolver) {
	resolvers[resolver.Type] = resolver
}

func GetResolver(depType string) *DependencyResolver {
	if resolver, exists := resolvers[depType]; exists {
		return &resolver
	}
	return nil
}
