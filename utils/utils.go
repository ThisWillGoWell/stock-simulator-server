package utils

type Subscribe interface {
	update() chan Subscribe
} 