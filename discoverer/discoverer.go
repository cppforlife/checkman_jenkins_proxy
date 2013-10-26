package discoverer

type Peer interface {
	Uuid() string
}

type DiscoverCallback func(peer Peer)

type Discoverer interface {
	Discover(f DiscoverCallback) error
}
