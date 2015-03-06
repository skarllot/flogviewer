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
	"github.com/go-gorp/gorp"
	"github.com/skarllot/gocli"
)

type WebFilterCommand struct {
	Dbm    *gorp.DbMap
	list   WebFilterList
	filter WebFilterList
}

func NewWebFilterCommand(dbm *gorp.DbMap) *WebFilterCommand {
	return &WebFilterCommand{
		Dbm:    dbm,
		list:   make(WebFilterList, 0),
		filter: make(WebFilterList, 0),
	}
}

func (self *WebFilterCommand) LoadWlog(cmd *gocli.Command) {
	for _, v := range wlogChilds {
		switch v.Name {
		case "load":
			v.Run = self.Load
		case "save":
			v.Run = self.SaveToFile
		}
		cmd.AddChild(v)
	}

	cmdFilter := cmd.Find("filter")
	for _, v := range filterChilds {
		switch v.Name {
		case "category":
			v.Run = self.FilterCategory
		case "dstip":
			v.Run = self.FilterDstIp
		case "hostname":
			v.Run = self.FilterHostname
		case "month":
			v.Run = self.FilterMonth
		case "reset":
			v.Run = self.ResetFilters
		case "srcip":
			v.Run = self.FilterSrcIp
		case "status":
			v.Run = self.FilterStatus
		case "user":
			v.Run = self.FilterUser
		}
		cmdFilter.AddChild(v)
	}

	cmdStatistics := cmd.Find("stats")
	for _, v := range statsChilds {
		switch v.Name {
		case "hits":
			v.Run = self.StatsHits
		case "trafficin":
			v.Run = self.StatsTrafficIn
		}
		cmdStatistics.AddChild(v)
	}
}
