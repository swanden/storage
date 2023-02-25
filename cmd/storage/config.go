package main

import (
	"github.com/swanden/storage/pkg/conf"
	"time"
)

const (
	defaultMemcachedPort                = 11211
	defaultGRPCListenerPort             = 8001
	defaultGRPCMAXRcv                   = 26214400
	defaultGRPCTimeOutConnection        = 3 * time.Second
	defaultHandlerWorkerPoolSize        = 100
	defaultMemcachedMaxIdleConns        = 10
	defaultMemcachedMaxOpenConns        = 10
	defaultMemcachedNewConnTimeout      = 3 * time.Second
	defaultMemcachedDefaultRetryTimeout = 3000 * time.Millisecond
)

type config struct {
	ServiceName                 string
	StorageType                 string
	MemcachedHost               string
	MemcachedPort               int
	MemcachedMaxIdleConns       int
	MemcachedMaxOpenConns       int
	MemcachedNewConnTimeout     time.Duration
	MemcachedConnRetryTimeout   time.Duration
	LogLevel                    string
	HandlerWorkerPoolSize       int
	GRPCServerListenerPort      int
	GRPCServerMaxRcvSize        int
	GRPCServerTimeOutConnection time.Duration
}

func newConfig() config {
	return config{
		ServiceName:                 conf.StrValueRequired("SERVICE_NAME"),
		StorageType:                 conf.StrValueRequired("STORAGE_TYPE"),
		MemcachedHost:               conf.StrValueRequired("MEMCACHED_HOST"),
		MemcachedPort:               conf.IntValue("MEMCACHED_Port", defaultMemcachedPort),
		MemcachedMaxIdleConns:       conf.IntValue("MEMCACHED_MAX_IDLE_CONNS", defaultMemcachedMaxIdleConns),
		MemcachedMaxOpenConns:       conf.IntValue("MEMCACHED_MAX_OPEN_CONNS", defaultMemcachedMaxOpenConns),
		MemcachedNewConnTimeout:     conf.TimeDurValue("MEMCACHED_NEW_CONN_TIMEOUT", defaultMemcachedNewConnTimeout),
		MemcachedConnRetryTimeout:   conf.TimeDurValue("MEMCACHED_CONN_RETRY_TIMEOUT", defaultMemcachedDefaultRetryTimeout),
		LogLevel:                    conf.StrValue("LOG_LEVEL", "info"),
		HandlerWorkerPoolSize:       conf.IntValue("HANDLER_WP_SIZE", defaultHandlerWorkerPoolSize),
		GRPCServerListenerPort:      conf.IntValue("GRPC_SERVER_LISTENER_PORT", defaultGRPCListenerPort),
		GRPCServerMaxRcvSize:        conf.IntValue("GRPC_SERVER_MAX_RECIVE_SIZE", defaultGRPCMAXRcv),
		GRPCServerTimeOutConnection: conf.TimeDurValue("GRPC_SERVER_TIME_OUT_CONNECTION", defaultGRPCTimeOutConnection),
	}
}
