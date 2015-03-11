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

	countRecord := int64(0)
	switch args[0] {
	case "dir":
		files, err := ioutil.ReadDir(args[1])
		if err != nil {
			fError(err)
			return
		}
		for _, f := range files {
			fname := path.Join(args[1], f.Name())
			count, err := wfc.LoadFile(fname, wfc.Dbm)
			if err != nil {
				fError(err)
				return
			}
			countRecord += count
		}
	case "file":
		count, err := wfc.LoadFile(args[1], wfc.Dbm)
		if err != nil {
			fError(err)
			return
		}
		countRecord += count
	default:
		fError(errors.New("You must choose between dir and file"))
		return
	}

	fmt.Printf(
		"Parsing log files done [%v  %d items]\n",
		time.Now().Sub(dt1), countRecord)
}

func (wfc *WebFilterCommand) LoadFile(
	fname string,
	dbm *gorp.DbMap) (int64, error) {

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
		return 0, nil
	}

	txn, err := dbm.Begin()
	if err != nil {
		return 0, err
	}
	defer txn.Rollback()

	fileRow, err := dal.GetFileByDate(txn, dt1)
	if err != nil {
		return 0, err
	} else if fileRow == nil {
		fileRow = &models.File{
			Begin: dt1,
			Count: 0,
		}
		if !isPresent {
			fileRow.End = dt2
		}
		err = txn.Insert(fileRow)
		if err != nil {
			return 0, err
		}
	} else if fileRow.End.Year() > 1 {
		return 0, nil
	}

	var reader io.Reader
	file, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return 0, err
	}
	defer gz.Close()
	reader = gz

	countInitial := fileRow.Count
	err = ParseFile(reader, fileRow, dbm)
	if err != nil {
		return 0, err
	}

	if !isPresent {
		fileRow.End = dt2
	}
	txn.Update(fileRow)
	txn.Commit()
	return fileRow.Count - countInitial, nil
}
