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
		Name:  "filter",
		Short: "Apply filters to current loaded logs",
	},
	&gocli.Command{
		Name:  "load",
		Short: "Load specified log file into memory",
	},
	&gocli.Command{
		Name:  "save",
		Short: "Saves current filtered data",
	},
}

var filterChilds = gocli.Commands{
	&gocli.Command{
		Name:  "user",
		Short: "Filter logs to specified user",
	},
	&gocli.Command{
		Name:  "reset",
		Short: "Reset all applied filters",
	},
}
