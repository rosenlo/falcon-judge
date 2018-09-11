// Copyright 2018 RosenLo

// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/**
 * This code was originally worte by Xiaomi, Inc. modified by RosenLo.
**/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/falcon-plus/modules/judge/cron"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
	"github.com/open-falcon/falcon-plus/modules/judge/http"
	"github.com/open-falcon/falcon-plus/modules/judge/rpc"
	"github.com/open-falcon/falcon-plus/modules/judge/store"
)

var (
	GitCommit string
	BuildTime string
	Version   string
)

func init() {
	if GitCommit == "" {
		GitCommit = "UnKnown"
	}

	if BuildTime == "" {
		BuildTime = "UnKnown"
	}
}

func displayVersion() {
	fmt.Println("Git Commit: ", GitCommit)
	fmt.Println("Build Time: ", BuildTime)
}

func start_signal(pid int) {
	sigs := make(chan os.Signal, 1)
	log.Println(pid, "register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			store.SaveEventToFile()
			log.Println("graceful shut down")
			log.Println(pid, "exit")
			os.Exit(0)
		}
	}
}

func main() {
	displayVersion()
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitRedisConnPool()
	g.InitHbsClient()

	store.InitHistoryBigMap()
	store.ReadEventFromFile()

	go http.Start()
	go rpc.Start()

	go cron.SyncStrategies()
	go cron.CleanStale()

	start_signal(os.Getpid())
}
