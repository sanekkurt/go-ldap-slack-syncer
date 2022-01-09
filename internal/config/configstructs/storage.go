package configstructs

type Storage struct {
	MySQL MySQL `yaml:"mysql"`
}

type MySQL struct {
	Address  string `yaml:"address"`
	Port     uint   `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
