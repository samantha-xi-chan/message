package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

const PATH = "config/msg.yaml"

func init() {
	viper.SetConfigFile(PATH)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	log.Println("config init()")
}

func GetWaiterPortRpc() (string, error) {
	port := viper.GetString("waiter.port_rpc")
	return port, nil
}

func GetDependQueue() (string, error) {
	v := viper.GetString("depend.queue")

	return v, nil
}

func GetDependMongo() (string, error) {
	v := viper.GetString("depend.mongo_dsn")

	return v, nil
}

func GetGwPortHttp() (string, error) {
	port := viper.GetString("gateway.port_http")
	return port, nil
}

func GetNotifyPortHttp() (string, error) {
	port := viper.GetString("notify.port_http")
	return port, nil
}

func GetDebugMode() (bool, error) {
	v := viper.GetBool("debug.debug_mode")
	return v, nil
}
func GetDebugPprofWaiter() (string, error) {
	v := viper.GetString("debug.pprof_waiter")
	return v, nil
}
func GetDebugPprofSink() (string, error) {
	v := viper.GetString("debug.pprof_sink")
	return v, nil
}
func GetDebugPprofGateway() (string, error) {
	v := viper.GetString("debug.pprof_gateway")
	return v, nil
}
func GetDebugPprofNotify() (string, error) {
	v := viper.GetString("debug.pprof_notify")
	return v, nil
}

func GetDebugSlowThreshold() (int, error) {
	v := viper.GetInt("debug.slow_threshold")
	return v, nil
}

func GetDebugLogRt() (bool, error) {
	v := viper.GetBool("debug.log_rt")
	return v, nil
}
