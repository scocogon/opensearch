package httpclient

type Config struct {
	accessKeyID     string
	accessKeySecret string
	host            string
}

func NewConfig(accessKeyID, accessKeySecret, host string) *Config {
	return &Config{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		host:            host,
	}
}
