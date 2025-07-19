package config

type etcdConfig struct {
	DialTimeout          int `mapstructure:"dial_timeout"`
	DialKeepAliveTime    int `mapstructure:"dial_keep_alive_time"`
	DialKeepAliveTimeout int `mapstructure:"dial_keep_alive_timeout"`
	Endpoint             endpoint
}
