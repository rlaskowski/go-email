package email

type Config struct {
	Key         string     `yaml:"key"`
	Description string     `yaml:"description"`
	SMTP        ServerInfo `yaml:"smtp"`
	POP3        ServerInfo `yaml:"pop3"`
	Email       string     `yaml:"email"`
	Username    string     `yaml:"username"`
	Password    string     `yaml:"password"`
}

type ServerInfo struct {
	Hostname   string `yaml:"hostname"`
	Port       int    `yaml:"port"`
	Encryption bool   `yaml:"encryption"`
}
