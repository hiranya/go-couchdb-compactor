// Copyright 2016 Hiranya Samarasekera
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/hiranya/gocouchlib"
	"github.com/hiranya/goroutinepool"
)

var fServerURL = flag.String("s", "http://localhost:5984", "CouchDB server url. Defaults to http://localhost:5984")
var fConcurrency = flag.Int("c", 5, "Concurrency level required for compaction")
var fUsername = flag.String("u", "", "Username to access the CouchDB server")
var fPassword = flag.String("p", "", "Password to access the CouchDB server")

type compactorCommand struct {
	DbNameStream chan string
	Server       *gocouchlib.Server
}

func (ccmd *compactorCommand) Exec() {
	startTime := time.Now()
	dbName := <-ccmd.DbNameStream

	db := &gocouchlib.Database{dbName, ccmd.Server}
	accepted, _ := db.Compact()
	if accepted {
		log.Println("Compaction running for", dbName)

		for {
			dbInfo, _ := db.Info()
			compactRunning := dbInfo.(map[string]interface{})["compact_running"]
			if !bool(compactRunning.(bool)) {
				timeDiff := time.Now().Sub(startTime)
				log.Printf("Compaction completed for %s (ET: %fs)", dbName, timeDiff.Seconds())
				break
			}
			// sleep before making the next check
			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func main() {
	flag.Parse()

	server := &gocouchlib.Server{
		*fServerURL, url.UserPassword(*fUsername, *fPassword),
	}

	allDbs, _ := json.Marshal(server.AllDbs())
	alldbs_value_json, _ := jason.NewValueFromBytes(allDbs)

	alldbs_array, _ := alldbs_value_json.Array()
	dbsToCompact := make([]string, 1)

	// remove meta DBs
	for _, db := range alldbs_array {
		dbName, _ := db.String()

		if !strings.HasPrefix(dbName, "_") && !strings.HasSuffix(dbName, "_") {
			dbsToCompact = append(dbsToCompact, dbName)
		}
	}
	log.Println(dbsToCompact)

	compactorCmd := &compactorCommand{}
	compactorCmd.DbNameStream = make(chan string, len(dbsToCompact))
	compactorCmd.Server = server

	// stream DB names to Command(s)
	for _, db := range dbsToCompact {
		compactorCmd.DbNameStream <- db
	}

	pool := &goroutinepool.Pool{}
	pool.Run(len(dbsToCompact), *fConcurrency, compactorCmd)
}
