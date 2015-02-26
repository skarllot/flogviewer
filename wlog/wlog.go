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
		cmd.AddChild(v)
	}
	cmd.Find("load").Run = wfc.Load
	cmd.Find("save").Run = wfc.SaveToFile

	cmdFilter := cmd.Find("filter")
	for _, v := range filterChilds {
		cmdFilter.AddChild(v)
	}
	cmdFilter.Find("month").Run = wfc.FilterMonth
	cmdFilter.Find("reset").Run = wfc.ResetFilters
	cmdFilter.Find("user").Run = wfc.FilterUser
}
