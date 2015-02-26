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
	"compress/gzip"
	"fmt"
	"github.com/skarllot/flogviewer/common"
	"github.com/skarllot/gocli"
	"io"
	"os"
	"time"
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
	cmd.Find("load").Run = wfc.LoadFile
	cmd.Find("save").Run = wfc.SaveToFile

	cmdFilter := cmd.Find("filter")
	for _, v := range filterChilds {
		cmdFilter.AddChild(v)
	}
	cmdFilter.Find("user").Run = wfc.FilterUser
	cmdFilter.Find("reset").Run = wfc.ResetFilters
}

func (wfc *WebFilterCommand) LoadFile(cmd *gocli.Command, args []string) {
	args = common.ParseParameters(args)
	if len(args) != 2 {
		fmt.Println("Two parameters must be defined")
		fmt.Println("<plain|gzip> <path>")
		return
	}

	var reader io.Reader
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	switch args[0] {
	case "plain":
		reader = file
	case "gzip":
		gz, err := gzip.NewReader(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer gz.Close()
		reader = gz
	default:
		fmt.Println("Invalid first parameter, must be plain or gzip")
		return
	}

	dt1 := time.Now()
	fmt.Print("Parsing logs...")
	wfc.list = append(wfc.list, ParseFile(reader)...)
	if err != nil {
		fmt.Println(" error")
		fmt.Println(err)
		return
	}

	wfc.filter = wfc.list
	fmt.Printf(" done [%v  %d items]\n", time.Now().Sub(dt1), len(wfc.list))
}

func (wfc *WebFilterCommand) SaveToFile(cmd *gocli.Command, args []string) {
	args = common.ParseParameters(args)
	if len(args) != 1 {
		fmt.Println("One file name must be defined")
		return
	}

	file, err := os.Create(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString("Date;Device;Level;PolicyId;User;SourceIP;SourceIf;" +
		"DestIP;DestPort;DestIf;Service;Hostname;Profile;Status;Url;Message;" +
		"CategoryId;CategoryDesc\n")
	for _, v := range wfc.filter {
		file.WriteString(fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;"+
			"%v;%v;%v;%v;%v\n", v.Date, v.Device, v.LogLevel, v.PolicyId, v.User,
			v.SourceIP, v.SourceIf, v.DestIP, v.DestPort, v.DestIf, v.Service,
			v.Hostname, v.Profile, v.Status, v.Url, v.Message, v.CategoryId,
			v.CategoryDesc))
	}
	file.Sync()
}
