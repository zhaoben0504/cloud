package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/go-xorm/xorm"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/text/language"
	xormcore "xorm.io/core"

	// not need
	_ "github.com/go-sql-driver/mysql"
)

var server = &Server{}

// Server server define
type Server struct {
	Config      *Config
	Engine      *xorm.EngineGroup
	Mode        string
	Node        *snowflake.Node
	redisClient *redis.Client
	minioClient *minio.Client
	bundle      *tool.Bundle
}

func NewServer(configPath, mode string) error {
	config, err := LoadLocalConfig(configPath, mode)
	if err != nil {
		tool.Logger.Error(err.Error())
		return err
	}
	server.Config = config

	engine, err := initEngine(config.DB)
	if err != nil {
		tool.Logger.Error(err.Error())
		return err
	}
	server.Engine = engine

	node, err := snowflake.NewNode(getServerID())
	if err != nil {
		tool.Logger.Error(err.Error())
		return err
	}
	server.Node = node

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Pwd, // no password set
		DB:       0,                // use default DB
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if nil != err {
		tool.Logger.Error(err.Error())
		return err
	}
	server.redisClient = redisClient

	minioClient, err := minio.New(config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.AccessKeyID, config.Minio.SecretAccessKey, ""),
		Secure: config.Minio.UseSSL,
	})
	if err != nil {
		tool.Logger.Error(err.Error())
		return err
	}
	server.minioClient = minioClient
	server.bundle = tool.NewBundle(language.Chinese)

	return nil
}

func LoadMessageFile(messageFiles []string) {
	for _, f := range messageFiles {
		server.bundle.MustLoadMessageFile(f)
	}
}

func getServerID() int64 {
	node := os.Getenv(ServiceID)

	if node == "" {
		tool.Logger.Fatalf("env %s is absent", ServiceID)
	}

	num, err := strconv.Atoi(node)
	if err != nil {
		tool.Logger.Fatalf("env %s must be int type , wrong value is %s", ServiceID, node)
	}

	return int64(num)
}

func initEngine(config *DBConfig) (*xorm.EngineGroup, error) {
	if config == nil || len(config.DataSources) == 0 {
		tool.Logger.Error("the db config of data sources is empty, Server.Engine is not init")
		return nil, errors.New("the db config of data sources is empty, Server.Engine is not init")
	}

	engineGroup, err := xorm.NewEngineGroup("mysql", config.DataSources)
	if nil != err {
		tool.Logger.Error(err.Error())
		return nil, err
	}
	err = engineGroup.Ping()
	if nil != err {
		tool.Logger.Error(err.Error())
		return nil, err
	}
	engineGroup.SetMapper(xormcore.GonicMapper{})
	engineGroup.SetMaxIdleConns(config.MaxIdleCon)
	engineGroup.SetMaxOpenConns(config.MaxCon)
	engineGroup.ShowSQL(true)
	engineGroup.ShowExecTime(true)

	tool.Logger.Debugf("connected to databases: %s", formatDataSources(config.DataSources))

	return engineGroup, nil
}

// GetID id generage
func GetID() int64 {
	return int64(server.Node.Generate())
}

// GetRedisClient redis client
func GetRedisClient() *redis.Client {
	return server.redisClient
}

// formatDataSources 格式化data source, 去掉敏感的用户名密码
func formatDataSources(sources []string) string {
	var formatSources []string

	for _, item := range sources {
		formatSources = append(formatSources, strings.Split(item, "@")[1])
	}

	return fmt.Sprintf("%v", formatSources)
}

func GetEngine() *xorm.EngineGroup {
	return server.Engine
}

func GetPort() int {
	return server.Config.Port
}

func GetMsgByCode(lang string, code int) string {
	return server.bundle.GetMsgByCode(lang, code)
}
