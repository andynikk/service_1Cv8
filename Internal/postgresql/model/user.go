package model

type User struct {
	Type         string `yaml:"type"`
	Name         string `yaml:"login"`
	Password     string `yaml:"password"`
	HashPassword string `yaml:"hash_password"`
	Event        string `yaml:"event"`
	New          bool   `yaml:"new"`
}
