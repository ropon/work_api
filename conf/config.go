package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/ropon/logger"
	"io/ioutil"
	"log"
	"os"
)

const SERVERNAME = "work_api"

var (
	configFileName string
	Cfg            Config
	MysqlDb        *gorm.DB
	RedisCi        *redis.Client
)

type LogCfg struct {
	Level     string
	FilePath  string
	FileName  string
	MaxSize   int64
	SplitFlag bool
	TimeDr    float64
}

type MysqlCfg struct {
	MysqlConn            string
	MysqlConnectPoolSize int
}

type RedisCfg struct {
	RedisConn   string
	RedisPasswd string
	RedisDb     int
}

// Config 配置文件结构体
type Config struct {
	LogCfg        LogCfg
	MysqlCfg      map[string]MysqlCfg
	RedisCfg      map[string]RedisCfg
	External      map[string]string
	ExternalInt64 map[string]int64
	Listen        string
}

func initConfigFile(filename string, cfg *Config) error {
	//fmt.Println("filename", filename)
	_, err := os.Stat(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return err
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return err
	}
	if err := json.Unmarshal(bytes, cfg); err != nil {
		err = fmt.Errorf("unmarshal error :%s", string(bytes))
		log.Println(err.Error())
		return err
	}
	//fmt.Println("config :", *cfg)
	return nil
}

func initConf() (err error) {
	flag.StringVar(&configFileName, "c", configFileName, "config file")
	flag.Parse()
	if configFileName != "" {
		err = initConfigFile(configFileName, &Cfg)
	}
	if err != nil {
		return
	}
	logCfg := logger.LogCfg(Cfg.LogCfg)
	err = logger.InitLog(&logCfg)
	return
}

func initGormDbPool(cfg *MysqlCfg, setLog bool) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", cfg.MysqlConn)
	if err != nil {
		fmt.Println("init db err : ", cfg, err)
		return nil, err
	}
	db.DB().SetMaxOpenConns(cfg.MysqlConnectPoolSize)
	db.DB().SetMaxIdleConns(cfg.MysqlConnectPoolSize / 2)
	if setLog {
		db.LogMode(true)
		db.SetLogger(logger.Log)
	}
	db.SingularTable(true)
	return db, nil
}

func initMysql() (err error) {
	mysqlConn := Cfg.MysqlCfg[SERVERNAME]
	MysqlDb, err = initGormDbPool(&mysqlConn, true)
	return
}

func initRedisPool(cfg *RedisCfg) (*redis.Client, error) {
	redisCi := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConn,
		Password: cfg.RedisPasswd,
		DB:       cfg.RedisDb,
	})
	_, err := redisCi.Ping().Result()
	return redisCi, err
}

func initRedis() (err error) {
	redisConn := Cfg.RedisCfg[SERVERNAME]
	RedisCi, err = initRedisPool(&redisConn)
	return
}

// Init 初始化配置文件、Mysql、Redis、etcd等
func Init() (err error) {
	err = initConf()
	if err != nil {
		return
	}
	return nil
	//return initMysql()

	//return initRedis()
}
