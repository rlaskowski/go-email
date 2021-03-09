package email

type Config struct {
	Key         string `yaml:"key"`
	Description string `yaml:"description"`
	Hostname    string `yaml:"hostname"`
	Port        string `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
}
