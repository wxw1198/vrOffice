package main

import (
	"flag"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/nats-io/nats-server/v2/conf"
	"github.com/wxw1198/vrOffice/db"
	"github.com/wxw1198/vrOffice/log"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
	"github.com/wxw1198/vrOffice/userbaseoperation/server"
	"os"
	"reflect"
	"strings"

	"time"
)

type options struct {
	mysqlAddr  string //ip:port
	mysqlUser  string
	mysqlPwd   string
	etcdv3Addr string //ip:port

	cacheAddr  string //ip:port
	cacheUser  string
	cachePwd   string
	cacheDBSeq int
}

var defautOpts = options{
	mysqlAddr:  "127.0.0.1:3306",
	mysqlUser:  "wxw",
	mysqlPwd:   "123",
	etcdv3Addr: "127.0.0.1:2379",

	cacheAddr:  "127.0.0.1:6379",
	cacheUser:  "",
	cachePwd:   "123",
	cacheDBSeq: 0,
}

var usageStr = `
Usage: register server [options]

 register Server Options:
	--etcdaddr                       etcd addr (default:127.0.0.1:2379)
	--configfile                     config file path and name

register server mysql  Options:
    --mysqladdr  <string>            mysql addr (default: 127.0.0.1:3306)
    --mysqluser  <string>            mysql user name (default: test)
    --mysqlpwd   <string>            mysql password
   
 register Server redis Options:
    --cacheaddr  <string>            redis addr(default:127.0.0.1:6379)
    --cacheuser  <string>            redis user name
    --cachepwd   <string>            redis password
    --cacheDBSeq <int>               Database to be selected after connecting to the server

Common Options:
    --help                           Show this message
    --version                        Show version       
`

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func main() {
	//todo config
	opts := parseConfig(&defautOpts)

	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			opts.etcdv3Addr,
		}
	})

	service := micro.NewService(
		micro.Name("go.micro.srv.register"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(reg),
		//micro.Address("127.0.0.1:5555"),
	)

	os.Args = os.Args[:1]
	// optionally setup command line usage
	service.Init()

	sqlOpts := db.MysqlOptions{
		User: opts.mysqlUser,
		Pwd:  opts.mysqlPwd,
		Addr: opts.mysqlAddr,
	}

	cacheOpts := db.CacheOptions{
		Addr:  opts.cacheAddr,
		Pwd:   opts.cachePwd,
		User:  opts.cacheUser,
		DBSeq: opts.cacheDBSeq,
	}

	proto.RegisterUserBaseOpsHandler(service.Server(), server.NewUserBaseOpsServer(&sqlOpts, &cacheOpts))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(o *options) *options {
	fs := flag.NewFlagSet("register", flag.ExitOnError)
	fs.Usage = usage

	opts := configureOpts(fs,
		os.Args[1:],
		func() {
			fmt.Printf("register server version %s, ", "v1.10")
		},
		fs.Usage)
	return opts

}

func configureOpts(fs *flag.FlagSet, args []string, printVersion, printHelp func()) *options {
	opts := defautOpts
	var (
		configFile string
	)

	fs.StringVar(&opts.mysqlAddr, "mysqladdr", "127.0.0.1:3306", "stan.ID")
	fs.StringVar(&opts.etcdv3Addr, "etcdaddr", "127.0.0.1:2379", "stan.ID")
	fs.StringVar(&opts.mysqlUser, "mysqluser", "test", "db user name")
	fs.StringVar(&opts.mysqlPwd, "mysqlpwd", "test", "db pwd")

	fs.StringVar(&opts.cacheAddr, "cacheaddr", "127.0.0.1:6379", "stan.ID")
	fs.StringVar(&opts.cacheUser, "cacheuser", "user", "db user name")
	fs.StringVar(&opts.cachePwd, "cachepwd", "test", "db pwd")
	fs.IntVar(&opts.cacheDBSeq, "cacheDBSeq", 0, "Database to be selected after connecting to the server.")

	fs.StringVar(&configFile, "configFile", "", "config file")

	fs.Parse(args)

	if configFile != "" {
		if err := ProcessConfigFile(configFile, &opts); err != nil {
			return nil
		}
		//cmd first
		fs.Parse(args)
	}

	return &opts
}

func checkType(name string, kind reflect.Kind, v interface{}) error {
	actualKind := reflect.TypeOf(v).Kind()
	if actualKind != kind {
		return fmt.Errorf("parameter %q value is expected to be %v, got %v",
			name, kind.String(), actualKind.String())
	}
	return nil
}

func ProcessConfigFile(configFile string, opts *options) error {
	m, err := conf.ParseFile(configFile)
	if err != nil {
		return err
	}

	for k, v := range m {
		name := strings.ToLower(k)
		switch name {
		case "mysqladdr":
			if err := checkType(k, reflect.String, v); err != nil {
				return err
			}
			opts.mysqlAddr = v.(string)
		case "etcdaddr":
			if err := checkType(k, reflect.String, v); err != nil {
				return err
			}
			opts.etcdv3Addr = v.(string)
		case "mysqlusername":
			if err := checkType(k, reflect.String, v); err != nil {
				return err
			}
			opts.mysqlUser = v.(string)

		case "mysqlpwd":
			if err := checkType(k, reflect.String, v); err != nil {
				return err
			}
			opts.mysqlPwd = v.(string)
		case "cacheaddr":
			if err := checkType(k, reflect.String, v); err != nil {
				return err
			}
			opts.cacheAddr = v.(string)
		case "cacheuser":
			if err := checkType(k, reflect.String, v); err != nil {
				return err
			}
			opts.cacheUser = v.(string)
		case "cachepwd":
			if err := checkType(k, reflect.String, v); err != nil {
				return err
			}
			opts.cachePwd = v.(string)
		case "cacheDBSeq":
			if err := checkType(k, reflect.Int, v); err != nil {
				return err
			}
			opts.cacheDBSeq = v.(int)
		}
	}
	return nil
}
