package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerPort      int
	DBPort          int
	DBHost          string
	DBUser          string
	DBPassword      string
	DBName          string
	Secret          string
	RefreshSecret   string
	CDN_Cloud_Name  string
	CDN_API_Key     string
	CDN_API_Secret  string
	CDN_Folder_Name string
}

func loadConfig() *Config {
	var res = new(Config)
	if val, found := os.LookupEnv("SERVER"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : invalid port value, ", err.Error())
			return nil
		}
		res.ServerPort = port
	}

	if val, found := os.LookupEnv("DBPORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : invalid port value, ", err.Error())
			return nil
		}
		res.DBPort = port
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		res.DBHost = val
	}

	if val, found := os.LookupEnv("DBUSER"); found {
		res.DBUser = val
	}

	if val, found := os.LookupEnv("DBPASSWORD"); found {
		res.DBPassword = val
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		res.DBName = val
	}

	if val, found := os.LookupEnv("SECRET"); found {
		res.Secret = val
	}

	if val, found := os.LookupEnv("REFSECRET"); found {
		res.RefreshSecret = val
	}
	if val, found := os.LookupEnv("CDN_Cloud_Name"); found {
		res.CDN_Cloud_Name = val
	}
	if val, found := os.LookupEnv("CDN_API_Key"); found {
		res.CDN_API_Key = val
	}
	if val, found := os.LookupEnv("CDN_API_Secret"); found {
		res.CDN_API_Secret = val
	}
	if val, found := os.LookupEnv("CDN_Folder_Name"); found {
		res.CDN_Folder_Name = val
	}

	return res
}

func InitConfig() *Config {
	var res = new(Config)

	res = loadConfig()
	if res == nil {
		logrus.Fatal("Config : Cannot start program, failed to load configuration")
		return nil
	}

	return res
}
