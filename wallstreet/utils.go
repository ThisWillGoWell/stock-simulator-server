package wallstreet

type Subscribe interface {
	update() chan Subscribe
} 