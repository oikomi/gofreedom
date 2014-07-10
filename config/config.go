//
// Copyright 2014 Hong Miao. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"encoding/json"
	"flag"
	"fmt"
)

type Config struct {
	Mode   string
	Listen string
	Server string

	Logfile    string
	Loglevel   string
	AdminIface string

	Cipher    string
	Key       string
	Blackfile string

	Username string
	Password string
	Auth     map[string]string
}

func LoadConfig(configpath string) (cfg Config, err error) {
	var configfile string
	flag.StringVar(&configfile, "config", configpath, "config file")
	flag.Parse()

	file, err := os.Open(configfile)
	if err != nil {
		return
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&cfg)
	if err != nil {
		return
	}
	return
}

func DumpConfig(cfg *Config) {
	fmt.Printf("Mode: %s\nListen: %s\nServer: %s\nLogfile: %s\n", 
	cfg.Mode, cfg.Listen, cfg.Server, cfg.Logfile)
}