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

package wlog

import (
	"fmt"
	"github.com/skarllot/gocli"
	"strings"
	"time"
)

func (wfc *WebFilterCommand) FilterUser(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("One user name must be specified")
		return
	}

	dt1 := time.Now()
	fmt.Print("Processing filter...")

	result := make(WebFilterList, 0)
	args[0] = strings.ToLower(args[0])
	for _, v := range wfc.filter {
		if strings.ToLower(v.User) == args[0] {
			result = append(result, v)
		}
	}
	wfc.filter = result
	fmt.Printf(" done [%v  %d items]\n", time.Now().Sub(dt1), len(wfc.filter))
}

func (wfc *WebFilterCommand) ResetFilters(cmd *gocli.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("This command takes no parameter")
		return
	}

	wfc.filter = wfc.list
}
