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

var wlogChilds = gocli.Commands{
	&gocli.Command{
		Name: "filter",
		Help: "Apply filters to current loaded logs",
	},
	&gocli.Command{
		Name: "load",
		Help: "Load specified log file into memory",
	},
	&gocli.Command{
		Name: "save",
		Help: "Saves current filtered data",
	},
	&gocli.Command{
		Name: "stats",
		Help: "Run statistical analysis on filtered data",
	},
}

var filterChilds = gocli.Commands{
	&gocli.Command{
		Name: "category",
		Help: "Filter logs to specified category",
	},
	&gocli.Command{
		Name: "dstip",
		Help: "Filter logs to specified destination IP",
	},
	&gocli.Command{
		Name: "hostname",
		Help: "Filter logs to specified hostname",
	},
	&gocli.Command{
		Name: "month",
		Help: "Filter logs to specified month",
	},
	&gocli.Command{
		Name: "reset",
		Help: "Reset all applied filters",
	},
	&gocli.Command{
		Name: "srcip",
		Help: "Filter logs to specified source IP",
	},
	&gocli.Command{
		Name: "status",
		Help: "Filter logs to specified status",
	},
	&gocli.Command{
		Name: "user",
		Help: "Filter logs to specified user",
	},
}

var statsChilds = gocli.Commands{
	&gocli.Command{
		Name: "hits",
		Help: "Group hits by specified target",
	},
	&gocli.Command{
		Name: "trafficin",
		Help: "Group input traffic by specified target",
	},
}
