/*
Copyright 2014 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"

	"github.com/golang/glog"
	marathon "github.com/gambol99/go-marathon"
)

var marathon_url string

func init() {
	flag.StringVar(&marathon_url, "url", "http://127.0.0.1:8080", "the url for the marathon endpoint")
}

func Assert(err error) {
	if err != nil {
		glog.Fatalf("Failed, error: %s", err)
	}
}

func main() {
	flag.Parse()
	config := marathon.NewDefaultConfig()
	config.URL = marathon_url
	if client, err := marathon.NewClient(config); err != nil {
		glog.Fatalf("Failed to create a client for marathon, error: %s", err)
	} else {
		applications, err := client.Applications()
		Assert(err)
		glog.Infof("Found %d application running", len(applications.Apps))
		for _, application := range applications.Apps {
			glog.Infof("Application: %s", application)
			details, err := client.Application(application.ID)
			Assert(err)
			if details.Tasks != nil && len(details.Tasks) > 0 {
				for _, task := range details.Tasks {
					glog.Infof("task: %s", task)
				}
				health, err := client.ApplicationOK(details.ID)
				Assert(err)
				glog.Infof("Application: %s, healthy: %t", details.ID, health )
			}
		}
	}
}