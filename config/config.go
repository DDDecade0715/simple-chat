package config

import "time"

type Server struct {
	App App `mapstructure:"app" json:"app" yaml:"app"`

	Redis Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
	// gorm
	Database Database `mapstructure:"database" json:"database" yaml:"database"`

	Amqp Amqp `mapstructure:"amqp" json:"amqp" yaml:"amqp"`

	Jwt Jwt `mapstructure:"jwt" json:"jwt" yaml:"jwt"`

	Es Es `mapstructure:"elaticsearch" json:"elaticsearch" yaml:"elaticsearch"`
}

type App struct {
	Root                 string        //项目绝对路径
	Mode                 string        `mapstructure:"mode" json:"mode" yaml:"mode"` //
	Port                 string        `mapstructure:"port" json:"port" yaml:"port"` // 服务器地址:端口
	ReadTimeout          time.Duration `mapstructure:"readTimeOut" json:"readTimeOut" yaml:"readTimeOut"`
	WriteTimeout         time.Duration `mapstructure:"writeTimeout" json:"writeTimeout" yaml:"writeTimeout"`
	UploadImageMaxSize   int           `yaml:"uploadImageMaxSize"`
	UploadPath           string        `yaml:"uploadPath"`
	UploadUrl            string        `yaml:"uploadUrl"`
	UploadImageAllowExts []string      `yaml:"uploadImageAllowExts"`
}

type Amqp struct {
	Url    string `mapstructure:"url" json:"url" yaml:"url"`          // amqp地址
	Config Config `mapstructure:"config" json:"config" yaml:"config"` // 选择配置
}

type Config struct {
	Channels Channels `mapstructure:"channels" json:"channels" yaml:"channels"`
}

type Channels struct {
	Default Params `mapstructure:"default" json:"default" yaml:"default"`
	Test    Params `mapstructure:"test" json:"test" yaml:"test"`
}

type Params struct {
	Exchange string `mapstructure:"exchange" json:"exchange" yaml:"exchange"`
	Queues   string `mapstructure:"queues" json:"queues" yaml:"queues"`
	Key      string `mapstructure:"key" json:"key" yaml:"key"`
}

type Es struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
}

type Database struct {
	Driver      string        `mapstructure:"driver" json:"driver" yaml:"driver"`
	Protocol    string        `mapstructure:"protocol" json:"protocol" yaml:"protocol"`
	Host        string        `mapstructure:"host" json:"host" yaml:"host"`
	Port        int           `mapstructure:"port" json:"port" yaml:"port"`
	User        string        `mapstructure:"user" json:"user" yaml:"user"`
	Password    string        `mapstructure:"password" json:"password" yaml:"password"`
	Name        string        `mapstructure:"name" json:"name" yaml:"name"`
	Prefix      string        `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	RunMode     string        `mapstructure:"runMode" json:"maxOpenConns" yaml:"max-open-conns"`
	MaxOpens    int           `mapstructure:"maxOpens" json:"maxOpens" yaml:"maxOpens"`
	MaxIdles    int           `mapstructure:"maxIdles" json:"maxIdles" yaml:"maxIdles"`
	MaxLifetime time.Duration `mapstructure:"maxLifetime" json:"maxLifetime" yaml:"maxLifetime"`
}

type Jwt struct {
	Secret string        `mapstructure:"secret" json:"secret" yaml:"secret"`
	Issuer string        `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
	Expire time.Duration `mapstructure:"expire" json:"expire" yaml:"expire"`
}

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
}
