package email

type Config struct {
	Key         string `yaml:"key"`
	Description string `yaml:"description"`
	Hostname    string `yaml:"hostname"`
	Port        int    `yaml:"port"`
	Email       string `yaml:"email"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
}
