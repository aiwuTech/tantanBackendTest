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
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/aiwuTech/tantanBackendTest/models"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"github.com/aiwuTech/tantanBackendTest/handlers"
)

var (
	version = kingpin.Flag("version", "show version").Short('v').Default("false").Bool()
	cfgFile = kingpin.Flag("cfg", "config file location").Short('c').Default("config.json").String()
	db      *gorm.DB
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	kingpin.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	cfg := ParseConfig(*cfgFile)
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var err error
	db, err = models.GetDB(cfg.DB.DSN, cfg.DB.MaxOpen, cfg.DB.MaxIdle, cfg.Debug)
	if err != nil {
		logrus.Fatalf("get db err: %v, exiting", err)
	}

	// init user adapter
	if err = handlers.InitUserAdapter(db); err != nil {
		logrus.Fatalf("init user err: %v, exiting", err)
	}

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if cfg.Debug {
		martini.Env = martini.Dev
	} else {
		martini.Env = martini.Prod
	}

	handlers.InitRouter(m)

	go m.RunOnAddr(cfg.Addr)
	handleSignal()
}

func handleSignal() {
	pid := os.Getpid()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	logrus.WithField("pid", pid).Info("has registered signal notify.")

	for {
		s := <-sigs
		logrus.Infof("has received signal: %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logrus.Info("is graceful shutting down...")

			db.Close()

			logrus.WithField("pid", pid).Info("has exited")
			os.Exit(0)
		}
	}
}
