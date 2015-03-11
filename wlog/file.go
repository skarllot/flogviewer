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
	"github.com/skarllot/flogviewer/dal"
	"github.com/skarllot/flogviewer/models"
	"github.com/skarllot/gocli"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"
)

const (
	FILENAME_WLOG_PATTERN1     = `[^_]+_wlog_(\d{8}-\d{4})-(\d{8}-\d{4}).log.gz$`
	FILENAME_WLOG_PATTERN2     = `[^_]+_wlog_(\d{8}-\d{4})-Present.log.gz$`
	FILENAME_WLOG_DATE_PATTERN = `20060102-1504`
)

func (wfc *WebFilterCommand) Load(cmd *gocli.Command, args []string) {
	args = common.ParseParameters(args)
	if len(args) != 2 {
		fmt.Println("Two parameters must be defined")
		fmt.Println("<dir|file> <path>")
		return
	}

	dt1 := time.Now()
	fError := func(e error) {
		fmt.Println("Parsing log files error")
		fmt.Println(e)
	}

	switch args[0] {
	case "dir":
		files, err := ioutil.ReadDir(args[1])
		if err != nil {
			fError(err)
			return
		}
		for _, f := range files {
			fname := path.Join(args[1], f.Name())
			err := wfc.LoadFile(fname, wfc.Dbm)
			if err != nil {
				fError(err)
				return
			}
		}
	case "file":
		if err := wfc.LoadFile(args[1], wfc.Dbm); err != nil {
			fError(err)
			return
		}
	default:
		fError(errors.New("You must choose between dir and file"))
		return
	}

	wfc.filter = wfc.list
	fmt.Printf("Parsing log files done [%v  %d items]\n", time.Now().Sub(dt1), len(wfc.list))
}

func (wfc *WebFilterCommand) LoadFile(fname string, dbm *gorp.DbMap) error {
	r1, _ := regexp.Compile(FILENAME_WLOG_PATTERN1)
	r2, _ := regexp.Compile(FILENAME_WLOG_PATTERN2)

	var dt1, dt2 time.Time
	var isPresent bool
	if m := r1.FindStringSubmatch(fname); m != nil {
		dt1, _ = time.Parse(FILENAME_WLOG_DATE_PATTERN, m[1])
		dt2, _ = time.Parse(FILENAME_WLOG_DATE_PATTERN, m[2])
		isPresent = false
	} else if m = r2.FindStringSubmatch(fname); m != nil {
		dt1, _ = time.Parse(FILENAME_WLOG_DATE_PATTERN, m[1])
		isPresent = true
	} else {
		fmt.Println("Skipped file:", fname)
		return nil
	}

	txn, err := dbm.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	fileRow, err := dal.GetFileByDate(txn, dt1)
	if err != nil {
		return err
	} else if fileRow == nil {
		fileRow = &models.File{
			Begin: dt1.Unix(),
			Count: 0,
		}
		if !isPresent {
			fileRow.End = dt2.Unix()
		}
		err = txn.Insert(fileRow)
		if err != nil {
			return err
		}
	} else if fileRow.End > 0 {
		return nil
	}

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

	err = ParseFile(reader, fileRow, dbm)
	if err != nil {
		return err
	}

	if !isPresent {
		fileRow.End = dt2.Unix()
	}
	txn.Update(fileRow)
	txn.Commit()
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
