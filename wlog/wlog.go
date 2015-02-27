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
	"github.com/skarllot/gocli"
)

type WebFilterCommand struct {
	list   WebFilterList
	filter WebFilterList
}

func LoadWlog(cmd *gocli.Command) {
	wfc := &WebFilterCommand{
		list:   make(WebFilterList, 0),
		filter: make(WebFilterList, 0),
	}

	for _, v := range wlogChilds {
		switch v.Name {
		case "load":
			v.Run = wfc.Load
		case "save":
			v.Run = wfc.SaveToFile
		}
		cmd.AddChild(v)
	}

	cmdFilter := cmd.Find("filter")
	for _, v := range filterChilds {
		switch v.Name {
		case "category":
			v.Run = wfc.FilterCategory
		case "dstip":
			v.Run = wfc.FilterDstIp
		case "hostname":
			v.Run = wfc.FilterHostname
		case "month":
			v.Run = wfc.FilterMonth
		case "reset":
			v.Run = wfc.ResetFilters
		case "srcip":
			v.Run = wfc.FilterSrcIp
		case "status":
			v.Run = wfc.FilterStatus
		case "user":
			v.Run = wfc.FilterUser
		}
		cmdFilter.AddChild(v)
	}

	cmdStatistics := cmd.Find("stats")
	for _, v := range statsChilds {
		switch v.Name {
		case "hits":
			v.Run = wfc.StatsHits
		case "trafficin":
			v.Run = wfc.StatsTrafficIn
		}
		cmdStatistics.AddChild(v)
	}
}
