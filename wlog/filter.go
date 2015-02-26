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
	"strconv"
	"strings"
	"time"
)

func (wfc *WebFilterCommand) Filter(f func(WebFilter) bool) {
	dt1 := time.Now()
	fmt.Print("Processing filter...")

	result := make(WebFilterList, 0)
	for _, v := range wfc.filter {
		if f(v) {
			result = append(result, v)
		}
	}
	wfc.filter = result
	fmt.Printf(" done [%v  %d items]\n", time.Now().Sub(dt1), len(wfc.filter))
}

func (wfc *WebFilterCommand) FilterDstIp(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("One destination IP must be specified")
		return
	}

	wfc.Filter(func(wf WebFilter) bool {
		return (strings.Index(wf.DestIP.String(), args[0]) == 0)
	})
}

func (wfc *WebFilterCommand) FilterHostname(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("One hostname must be specified")
		return
	}

	hostname := strings.ToLower(args[0])
	wfc.Filter(func(wf WebFilter) bool {
		return (strings.Index(strings.ToLower(wf.Hostname), hostname) != -1)
	})
}

func (wfc *WebFilterCommand) FilterMonth(cmd *gocli.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("One month and year value must be specified")
		return
	}

	month, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("The month value must be an integer")
		return
	}
	if month < 1 || month > 12 {
		fmt.Println("Invalid month value")
		return
	}
	year, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("The year value must be an integer")
		return
	}
	if year < 2000 {
		fmt.Println("Invalid year value")
		return
	}

	wfc.Filter(func(wf WebFilter) bool {
		return (int(wf.Date.Month()) == month && wf.Date.Year() == year)
	})
}

func (wfc *WebFilterCommand) FilterSrcIp(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("One source IP must be specified")
		return
	}

	wfc.Filter(func(wf WebFilter) bool {
		return (strings.Index(wf.SourceIP.String(), args[0]) == 0)
	})
}

func (wfc *WebFilterCommand) FilterUser(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("One user name must be specified")
		return
	}

	args[0] = strings.ToLower(args[0])
	wfc.Filter(func(wf WebFilter) bool {
		return (strings.ToLower(wf.User) == args[0])
	})
}

func (wfc *WebFilterCommand) ResetFilters(cmd *gocli.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("This command takes no parameter")
		return
	}

	wfc.filter = wfc.list
}
