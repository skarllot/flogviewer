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
	"errors"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/skarllot/flogviewer/common"
	"github.com/skarllot/gocli"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

func (wfc *WebFilterCommand) Load(cmd *gocli.Command, args []string) {
	args = common.ParseParameters(args)
	if len(args) != 2 {
		fmt.Println("Two parameters must be defined")
		fmt.Println("<dir|file> <path>")
		return
	}

	dt1 := time.Now()
	done := false
	var err error
	defer func() {
		if !done {
			fmt.Println("Parsing log files error")
			fmt.Println(err)
		}
	}()

	switch args[0] {
	case "dir":
		files, err := ioutil.ReadDir(args[1])
		if err != nil {
			return
		}
		for _, f := range files {
			fname := path.Join(args[1], f.Name())
			if strings.Index(fname, ".log.gz") == -1 ||
				strings.Index(fname, "wlog") == -1 {
				continue
			}
			err = wfc.LoadFile(fname, wfc.Dbm)
			if err != nil {
				return
			}
		}
	case "file":
		if err := wfc.LoadFile(args[1], wfc.Dbm); err != nil {
			return
		}
	default:
		err = errors.New("You must choose between dir and file")
		return
	}

	wfc.filter = wfc.list
	done = true
	fmt.Printf("Parsing log files done [%v  %d items]\n", time.Now().Sub(dt1), len(wfc.list))
}

func (wfc *WebFilterCommand) LoadFile(fname string, dbm *gorp.DbMap) error {
	var reader io.Reader
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()
	reader = gz

	err = ParseFile(reader, dbm)
	if err != nil {
		return err
	}
	return nil
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
		file.WriteString(fmt.Sprintf("\"%v\";\"%v\";\"%v\";\"%v\";\"%v\";"+
			"\"%v\";\"%v\";\"%v\";\"%v\";\"%v\";\"%v\";\"%v\";\"%v\";"+
			"\"%v\";\"%v\";\"%v\";\"%v\";\"%v\"\n",
			v.Date, v.Device, v.LogLevel, v.PolicyId, v.User, v.SourceIP,
			v.SourceIf, v.DestIP, v.DestPort, v.DestIf, v.Service,
			v.Hostname, v.Profile, v.Status, v.Url, v.Message, v.CategoryId,
			v.CategoryDesc))
	}
	file.Sync()
}
