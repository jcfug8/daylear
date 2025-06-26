package domain

// Domain defines how to interact with the domain or business logic layer.
type Domain interface {
	circleDomain
	recipeDomain
	userDomain
	tokenDomain
}
