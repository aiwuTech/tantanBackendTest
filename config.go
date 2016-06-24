// Copyright 2016 zm@huantucorp.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
/*
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
*/
package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/configor"
	"sync"
)

const (
	VERSION string = "0.0.1"
)

type GlobalConfig struct {
	Debug bool   `json:"debug" default:"true"`
	Addr  string `json:"addr" default:":3000"`
	DB    struct {
		DSN     string `json:"dsn"`
		MaxOpen int    `json:"maxOpen" default:"100"`
		MaxIdle int    `json:"maxIdle" default:"10"`
	} `json:"db"`
}

var (
	config *GlobalConfig
	lock   = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) *GlobalConfig {
	var c GlobalConfig
	configor.Load(&c, cfg)

	lock.Lock()
	config = &c
	lock.Unlock()

	logrus.Printf("ParseConfig ok, file: %s, content: %+v", cfg, c)
	return config
}

func Debug() bool {
	return Config().Debug
}
