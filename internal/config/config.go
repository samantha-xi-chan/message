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

func GetWaiterPortPprof() (string, error) {
	port := viper.GetString("waiter.port_pprof")

	return port, nil
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
func GetSinkPortPprof() (string, error) {
	port := viper.GetString("sink.port_pprof")

	return port, nil
}

func GetGwPortHttp() (string, error) {
	port := viper.GetString("gateway.port_http")
	return port, nil
}

func GetNotifyPortHttp() (string, error) {
	port := viper.GetString("notify.port_http")
	return port, nil
}
