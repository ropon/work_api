package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/ropon/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const SERVERNAME = "work_api"

var (
	etcdHost       string
	configFileName string
	Cfg            Config
	MysqlDb        *gorm.DB
	RedisCi        *redis.Client
	KafkaProducer  sarama.SyncProducer
	KafkaConsumer  sarama.Consumer
	Etcd           *clientv3.Client
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

type EtcdRes struct {
	Node struct {
		Key           string `json:"key"`
		Value         string `json:"value"`
		ModifiedIndex int    `json:"modifiedIndex"`
		CreatedIndex  int    `json:"createdIndex"`
	} `json:"node"`
}

func initConf() (err error) {
	flag.StringVar(&configFileName, "c", configFileName, "config file")
	flag.StringVar(&etcdHost, "etcd", etcdHost, "etcd addr")
	flag.Parse()
	if configFileName != "" {
		err = initConfigFile(configFileName, &Cfg)
	} else {
		err = loadCfgFromEtcd([]string{etcdHost}, SERVERNAME, &Cfg)
	}
	if err != nil {
		return
	}
	logCfg := logger.LogCfg(Cfg.LogCfg)
	err = logger.InitLog(&logCfg)
	return
}

func initConfigFile(filename string, cfg *Config) error {
	_, err := os.Stat(filename)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, cfg); err != nil {
		err = fmt.Errorf("unmarshal error :%s", string(bytes))
		return err
	}
	return nil
}

func loadCfgFromEtcd(addrs []string, service string, cfg interface{}) error {
	environment := strings.ToLower(os.Getenv("GOENV"))
	if environment == "" {
		environment = "online"
	}
	data, err := cfgFromEtcd(addrs, service, environment)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), cfg)
}

func cfgFromEtcd(addrs []string, service, env string) (string, error) {
	if len(addrs) == 0 || addrs[0] == "" {
		return "", fmt.Errorf("etcd地址不能为空")
	}
	addr := fmt.Sprintf("%s/v2/keys%s", addrs[0], etcdKey(service, env))
	resp, err := http.Get(addr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	etcdRes := EtcdRes{}
	err = json.Unmarshal(data, &etcdRes)
	if err != nil {
		log.Printf("unmarshal [%s] failed:%v", data, err)
		return "", err
	}
	return etcdRes.Node.Value, nil
}

func etcdKey(service, env string) string {
	return fmt.Sprintf("/config/%s/%s", service, env)
}

func initMysql() (err error) {
	mysqlConn := Cfg.MysqlCfg[SERVERNAME]
	MysqlDb, err = initGormDbPool(&mysqlConn, true)
	return
}

func initGormDbPool(cfg *MysqlCfg, setLog bool) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", cfg.MysqlConn)
	if err != nil {
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

func initRedis() (err error) {
	redisConn := Cfg.RedisCfg[SERVERNAME]
	RedisCi, err = initRedisPool(&redisConn)
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

func initKafka() (err error) {
	//生产者
	kafkaAddrs := strings.Split(Cfg.External["KafkaAddr"], ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	// 连接kafka
	KafkaProducer, err = sarama.NewSyncProducer(kafkaAddrs, config)
	if err != nil {
		return
	}
	//消费者
	KafkaConsumer, err = sarama.NewConsumer(kafkaAddrs, nil)
	return
}

func initEtcd() (err error) {
	Etcd, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(Cfg.External["EtcdConn"], ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	return nil
}

// Init 初始化配置文件、Mysql、Redis、etcd等
func Init() (err error) {
	err = initConf()
	if err != nil {
		return
	}
	//return nil
	return initMysql()
	//return initKafka()
	//return initEtcd()

	//return initRedis()
}
