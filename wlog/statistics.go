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
	"github.com/skarllot/flogviewer/common"
	"github.com/skarllot/gocli"
	"os"
	"sort"
)

type TrafficInStatistics struct {
	Hostname     string
	TrafficOut   uint64
	TrafficIn    uint64
	CategoryId   int
	CategoryDesc string
}

type TrafficInStatisticsList []TrafficInStatistics

func (wfc *WebFilterCommand) StatisticsTrafficIn(cmd *gocli.Command, args []string) {
	args = common.ParseParameters(args)
	if len(args) != 1 {
		fmt.Println("The destination file must be defined")
		return
	}

	list := make(map[string]TrafficInStatistics, 0)
	for _, v := range wfc.filter {
		wf, ok := list[v.Hostname]
		if !ok {
			list[v.Hostname] = TrafficInStatistics{
				Hostname:     v.Hostname,
				TrafficOut:   v.TrafficOut,
				TrafficIn:    v.TrafficIn,
				CategoryId:   v.CategoryId,
				CategoryDesc: v.CategoryDesc,
			}
		} else {
			wf.TrafficIn += v.TrafficIn
			wf.TrafficOut += v.TrafficOut
			list[v.Hostname] = wf
		}
	}

	result := make(TrafficInStatisticsList, 0, len(list))
	for _, v := range list {
		result = append(result, v)
	}
	sort.Sort(sort.Reverse(result))

	file, err := os.Create(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString("Hostname;TrafficIn;TrafficOut;CategoryId;CategoryDesc\n")
	for _, v := range result {
		file.WriteString(fmt.Sprintf(
			"\"%v\";\"%v\";\"%v\";\"%v\";\"%v\"\n",
			v.Hostname, v.TrafficIn, v.TrafficOut, v.CategoryId, v.CategoryDesc))
	}
	file.Sync()
}

func (a TrafficInStatisticsList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TrafficInStatisticsList) Len() int           { return len(a) }
func (a TrafficInStatisticsList) Less(i, j int) bool { return a[i].TrafficIn < a[j].TrafficIn }
