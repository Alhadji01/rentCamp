package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerPort    int
	DBPort        int
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	Secret        string
	RefreshSecret string
	OpenAIKey     string
}

func loadConfig() *Config {
	var res = new(Config)

	var err = godotenv.Load(".ENV")
	if err != nil {
		logrus.Error("Config : Cannot load config file, ", err.Error())
		return nil
	}

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
	if val, found := os.LookupEnv("OPENAI_API_KEY"); found {
		res.OpenAIKey = val
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
