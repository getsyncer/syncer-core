package config

type Syncs struct {
	Logic  Name    `yaml:"logic"`
	ID     string  `yaml:"id"`
	Config Dynamic `yaml:"config"`
}

type Name string

func (n Name) String() string {
	return string(n)
}
