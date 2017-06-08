package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/goregraph/lib-r/types"
	"github.com/synw/terr"
)

func GetConf(dev bool, verbosity int) (*types.Conf, *terr.Trace) {
	var conf *types.Conf
	// set some defaults for conf
	if dev {
		viper.SetConfigName("dev_config")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(".")
	viper.SetDefault("grg_type", "rethinkdb")
	viper.SetDefault("grg_host", "localhost:8081")
	viper.SetDefault("grg_addr", "localhost:28015")
	viper.SetDefault("grg_user", "")
	viper.SetDefault("grg_password", "")
	viper.SetDefault("grg_cors", []interface{}{})
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigParseError:
			trace := terr.New("conf.GetConf", err)
			return conf, trace
		default:
			err := errors.New("Unable to locate config file")
			trace := terr.New("conf.GetConf", err)
			return conf, trace
		}
	}
	dbtype := viper.GetString("grg_type")
	host := viper.GetString("grg_host")
	addr := viper.GetString("grg_addr")
	user := viper.GetString("grg_user")
	pwd := viper.GetString("grg_password")
	// headers
	cors := viper.GetStringSlice("grg_cors")
	endconf := &types.Conf{dbtype, host, addr, user, pwd, dev, verbosity, cors}
	return endconf, nil
}
