/*
* Copyright 2015 Fabrício Godoy
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/skarllot/flogviewer/bll"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := bll.Configuration{}
	if err := cfg.Load("flogviewer.gcfg"); err != nil {
		fmt.Println("Could not load configuration file:", err)
		return
	}

	dbm, err := cfg.CreateDbMap()
	if err != nil {
		fmt.Println("Could not initialize database:", err)
		return
	}

	for _, v := range RootChilds(dbm) {
		rootCmd.AddChild(v)
	}

	fmt.Println("Type help for help\n")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
