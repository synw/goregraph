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
	viper.SetDefault("type", "")
	viper.SetDefault("host", "localhost:8080")
	viper.SetDefault("addr", "localhost:28015")
	viper.SetDefault("user", "")
	viper.SetDefault("password", "")
	viper.SetDefault("cors", []string{})
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigParseError:
			trace := terr.New("conf.getConf", err)
			return conf, trace
		default:
			err := errors.New("Unable to locate config file")
			trace := terr.New("conf.getConf", err)
			return conf, trace
		}
	}
	dbtype := viper.Get("type").(string)
	host := viper.Get("host").(string)
	addr := viper.Get("addr").(string)
	user := viper.Get("user").(string)
	pwd := viper.Get("password").(string)
	if dbtype == "" {
		err := errors.New("Please set the database type into your config file: ex: \"type\":\"rethinkdb\"")
		tr := terr.New("conf.GetConf", err)
		terr.Fatal("loading configuration", tr)
	}
	// headers
	crs := viper.Get("cors").([]interface{})
	var cors []string
	for _, c := range crs {
		cors = append(cors, c.(string))
	}
	endconf := &types.Conf{dbtype, host, addr, user, pwd, dev, verbosity, cors}
	return endconf, nil
}
