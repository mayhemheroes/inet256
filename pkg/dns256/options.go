package dns256

type ResolverOpt func(c *Resolver)

type resolveConfig struct {
	maxHops int
	labels  map[string]string
}

type ResolveOpt func(c *resolveConfig)

func WithMaxHops(n int) ResolveOpt {
	return func(c *resolveConfig) {
		c.maxHops = n
	}
}

func WithLabels(filter map[string]string) ResolveOpt {
	return func(c *resolveConfig) {
		c.labels = filter
	}
}
