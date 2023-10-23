package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
)

const PATH = "config/msg.yaml"

func init() {
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		log.Println("Failed to get POD_NAME environment variable")

		viper.SetConfigFile(PATH)
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	} else {
		log.Printf("Pod Name: %s\n", podName)
	}

	log.Println("config init()")
}

func GetWaiterPortRpc() (string, error) {
	value := os.Getenv("WAITER_PORT_RPC")
	if value == "" {
		v := viper.GetString("waiter.port_rpc")
		return v, nil
	} else {
		return value, nil
	}
}

func GetDependQueue() (string, error) {
	value := os.Getenv("DEPEND_QUEUE")
	if value == "" {
		v := viper.GetString("depend.queue")
		return v, nil
	} else {
		return value, nil
	}
}

//func GetDependMongo() (string, error) {
//	v := viper.GetString("depend.mongo_dsn")
//
//	return v, nil
//}

func GetDependRedisDsn() (string, error) {

	value := os.Getenv("DEPEND_REDIS_DSN")
	if value == "" {
		v := viper.GetString("depend.redis_dsn")
		return v, nil
	} else {
		return value, nil
	}
}

//
//func GetStoreMaxCount() (int64, error) {
//
//	value := os.Getenv("STORE_MAX_COUNT")
//	if value == "" {
//		v := viper.GetInt64("store.max_count")
//		return v, nil
//	} else {
//		return value, nil
//	}
//}

func GetGwPortHttp() (string, error) {

	value := os.Getenv("GATEWAY_PORT_HTTP")
	if value == "" {
		port := viper.GetString("gateway.port_http")
		return port, nil
	} else {
		return value, nil
	}
}

func GetNotifyPortHttp() (string, error) {

	value := os.Getenv("NOTIFY_PORT_HTTP")
	if value == "" {
		port := viper.GetString("notify.port_http")
		return port, nil
	} else {
		return value, nil
	}
}

func GetDebugMode() (bool, error) {
	value := os.Getenv("DEBUG_DEBUG_MODE")
	if value == "" {
		v := viper.GetBool("debug.debug_mode")
		return v, nil
	} else {
		b, err := strconv.ParseBool(value)
		if err != nil {
			fmt.Println("转换失败:", err)
			return false, err
		} else {
			fmt.Println("转换后的 bool 值:", b)

			return b, nil
		}
	}
}
func GetDebugPprofWaiter() (string, error) {
	value := os.Getenv("DEBUG_PPROF_WAITER")
	if value == "" {
		v := viper.GetString("debug.pprof_waiter")
		return v, nil
	} else {
		return value, nil
	}
}
func GetDebugPprofSink() (string, error) {

	value := os.Getenv("DEBUG_PPROF_SINK")
	if value == "" {
		v := viper.GetString("debug.pprof_sink")
		return v, nil
	} else {
		return value, nil
	}
}
func GetDebugPprofGateway() (string, error) {

	value := os.Getenv("DEBUG_PPROF_GATEWAY")
	if value == "" {
		v := viper.GetString("debug.pprof_gateway")
		return v, nil
	} else {
		return value, nil
	}
}
func GetDebugPprofNotify() (string, error) {

	value := os.Getenv("DEBUG_PPROF_NOTIFY")
	if value == "" {
		v := viper.GetString("debug.pprof_notify")
		return v, nil
	} else {
		return value, nil
	}
}
