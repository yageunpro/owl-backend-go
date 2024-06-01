package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const PathEnvKey = "APP__CONFIG__PATH"

type oAuthConfig struct {
	Data   []byte   `json:"data"`
	Scopes []string `json:"scopes"`
}

type jwtConfig struct {
	AccessKey  string `json:"accessKey"`
	RefreshKey string `json:"refreshKey"`
}

type jsonData struct {
	OAuth oAuthConfig `json:"oauth"`
	JWT   jwtConfig   `json:"jwt"`
	Dsn   string      `json:"dsn"`
}

var OAuth *oAuthConfig
var JWT *jwtConfig
var DBDsn string

func init() {
	path, ok := os.LookupEnv(PathEnvKey)
	if !ok {
		excPath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(filepath.Dir(excPath), "config.json")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	data := new(jsonData)
	err = json.Unmarshal(raw, data)
	if err != nil {
		panic(err)
	}

	OAuth = &data.OAuth
	JWT = &data.JWT
	DBDsn = data.Dsn
}
