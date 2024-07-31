package dependencies

type DependencyHandler func(Dependency) ([]byte, error)

type DependencyResolver struct {
	Type string
	Handler DependencyHandler
}

var resolvers = []DependencyResolver{}

func RegisterResolver(resolver DependencyResolver) {
	resolvers = append(resolvers, resolver)
}