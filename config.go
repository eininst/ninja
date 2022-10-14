package ninja

import (
	"encoding/json"
	"github.com/eininst/flog"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var once sync.Once

type Config struct {
	Data      map[string]any
	Ret       gjson.Result
	AppConfig *AppConfig
}

type AppConfig struct {
	Port         string `json:"port"`
	Prefork      bool   `json:"prefork"`
	ReadTimeout  int64  `json:"readTimeout"`
	WriteTimeout int64  `json:"writeTimeout"`
}

func NewConfig(path string, profile ...string) *Config {
	var data map[string]any
	var ret gjson.Result

	prf := ""
	if len(profile) > 0 {
		prf = profile[0]
	} else {
		prfEnv := os.Getenv("profile")
		if prfEnv != "" {
			prf = prfEnv
		}
	}
	if prf != "" {
		flog.Infof("profile is: %s", prf)
	}

	file, err := os.Open(path)
	defer func() { _ = file.Close() }()
	if err != nil {
		log.Fatal(err)
	}
	dec := yaml.NewDecoder(file)
	err = dec.Decode(&data)

	for {
		var t map[string]interface{}
		err = dec.Decode(&t)
		if err != nil {
			break
		}
		if p, ok := t["profile"]; ok {
			if p == prf {
				for k, v := range t {
					data[k] = v
				}
				break
			}
		}
	}
	v, er := json.Marshal(&data)
	if er != nil {
		log.Println(er)
	}
	ret = gjson.Parse(string(v))

	return &Config{
		Data:      data,
		Ret:       ret,
		AppConfig: nil,
	}
}

func (c *Config) Get(path ...string) gjson.Result {
	if len(path) == 0 {
		return c.Ret
	}

	var r gjson.Result
	for _, p := range path {
		if r.Value() == nil {
			r = c.Ret.Get(p)
		} else {
			r = r.Get(p)
		}
	}
	return r
}

func (c *Config) GetAppConfig() *AppConfig {
	once.Do(func() {
		var acfg *AppConfig
		mstr := c.Get("ninja").String()
		_ = json.Unmarshal([]byte(mstr), &acfg)

		defaultAppConfig := &AppConfig{
			Port:         ":3000",
			Prefork:      false,
			ReadTimeout:  15,
			WriteTimeout: 15,
		}
		if acfg == nil {
			acfg = defaultAppConfig
		} else {
			if acfg.WriteTimeout == 0 {
				acfg.WriteTimeout = defaultAppConfig.WriteTimeout
			}
			if acfg.ReadTimeout == 0 {
				acfg.ReadTimeout = defaultAppConfig.ReadTimeout
			}
			if acfg.Port == "" {
				acfg.Port = defaultAppConfig.Port
			}
		}
		c.AppConfig = acfg
	})

	return c.AppConfig
}
func (c *Config) IsDev() bool {
	return c.Get("profile").String() == "dev"
}

func (c *Config) IsTest() bool {
	return c.Get("profile").String() == "test"
}

func (c *Config) IsUat() bool {
	return c.Get("profile").String() == "uat"
}

func (c *Config) IsProd() bool {
	return c.Get("profile").String() == "prod"
}
