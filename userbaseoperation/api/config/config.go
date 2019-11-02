package config

const (
	etcdAddr = ""
)

type config struct {
	Etcdv3          string //用于服务发现的地址
	MicroName       string //监听的url
	ListenLocalAddr string //本地绑定的地址
	AdvertiseAddr   string //本地服务器，所在的公网地址
}

var DefaultConfig = config{
	Etcdv3:          "127.0.0.1:2379",
	MicroName:       "go.micro.api.UserBaseOps",
	ListenLocalAddr: "0.0.0.0:8081",
	AdvertiseAddr:   "39.98.39.224:8081",
}

func ParseConfig(c *config) bool {
	return true
}
