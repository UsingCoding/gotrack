package config

type Parser interface {
	Parse(path string) (Config, error)
}
