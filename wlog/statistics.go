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
	"github.com/dustin/go-humanize"
	"github.com/skarllot/flogviewer/common"
	"github.com/skarllot/gocli"
	"os"
	"sort"
	"strings"
)

type TrafficInStats struct {
	Target     string
	TrafficOut uint64
	TrafficIn  uint64
}

type TrafficInStatsList []TrafficInStats

func (wfc *WebFilterCommand) statisticsTrafficIn(
	columnName string,
	fname string,
	f func(*WebFilter) string) {
	list := make(map[string]TrafficInStats, 0)
	for _, v := range wfc.filter {
		target := strings.ToLower(f(&v))
		wf, ok := list[target]
		if !ok {
			list[target] = TrafficInStats{
				Target:     target,
				TrafficOut: v.TrafficOut,
				TrafficIn:  v.TrafficIn,
			}
		} else {
			wf.TrafficIn += v.TrafficIn
			wf.TrafficOut += v.TrafficOut
			list[target] = wf
		}
	}

	result := make(TrafficInStatsList, 0, len(list))
	for _, v := range list {
		result = append(result, v)
	}
	sort.Sort(sort.Reverse(result))

	file, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString(columnName +
		";TrafficIn;TrafficOut;TrafficInBytes;TrafficOutBytes\n")
	for _, v := range result {
		file.WriteString(fmt.Sprintf(
			"\"%v\";\"%v\";\"%v\";\"%v\";\"%v\"\n",
			v.Target,
			humanize.Bytes(v.TrafficIn), humanize.Bytes(v.TrafficOut),
			v.TrafficIn, v.TrafficOut))
	}
	file.Sync()
}

func (wfc *WebFilterCommand) StatsTrafficIn(cmd *gocli.Command, args []string) {
	args = common.ParseParameters(args)
	if len(args) != 2 {
		fmt.Println("Two parameters must be defined")
		fmt.Println("<category|hostname|user> <path>")
		return
	}

	option := strings.ToLower(args[0])
	switch option {
	case "category":
		wfc.statisticsTrafficIn("Category", args[1],
			func(wf *WebFilter) string { return wf.CategoryDesc })
	case "hostname":
		wfc.statisticsTrafficIn("Hostname", args[1],
			func(wf *WebFilter) string { return wf.Hostname })
	case "user":
		wfc.statisticsTrafficIn("User", args[1],
			func(wf *WebFilter) string { return wf.User })
	default:
		fmt.Println("Invalid option. Must be category, hostname or user.")
		return
	}
}

func (a TrafficInStatsList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TrafficInStatsList) Len() int           { return len(a) }
func (a TrafficInStatsList) Less(i, j int) bool { return a[i].TrafficIn < a[j].TrafficIn }
