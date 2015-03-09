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
	"bufio"
	"errors"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/skarllot/flogviewer/common"
	"github.com/skarllot/flogviewer/dal"
	"github.com/skarllot/flogviewer/models"
	"io"
)

func ParseFile(r io.Reader, dbm *gorp.DbMap) error {
	scanner := bufio.NewScanner(r)
	c := make(chan *WebFilter)
	count := 0
	defer close(c)
	for scanner.Scan() {
		go ParseLine(scanner.Text(), c)
		count++
	}

	txn, err := dbm.Begin()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		item := <-c
		if item != nil {
			err = InsertWebFilter(*item, txn)
			if err != nil {
				fmt.Println("Error inserting new record:", err)
			}
		}
	}
	txn.Commit()

	return nil
}

func ParseLine(line string, c chan *WebFilter) {
	wf := &WebFilter{}
	if err := common.ParseKeyValueLog(line, wf); err != nil {
		c <- nil
		return
	}

	wf.ConvertFields()
	c <- wf
}

func InsertWebFilter(wf WebFilter, txn *gorp.Transaction) error {
	dev, err := dal.GetOrInsertDeviceBySerial(txn, wf.DeviceSerial, wf.Device)
	if err != nil {
		return err
	}

	logtype, err := dal.GetOrInsertLogtypeByNames(txn, wf.LogType, wf.LogSubType)
	if err != nil {
		return err
	}

	loglevel, err := dal.GetLoglevelByName(txn, wf.LogLevel)
	if err != nil {
		return err
	}
	if loglevel == nil {
		return errors.New("Invalid log level value: " + wf.LogLevel)
	}

	user, err := dal.GetOrInsertUserByName(txn, wf.User)
	if err != nil {
		return err
	}

	service, err := dal.GetOrInsertServiceByName(txn, wf.Service)
	if err != nil {
		return err
	}

	message := wf.Message
	log := &models.Log{
		Date:         wf.Date,
		PolicyId:     wf.PolicyId,
		SourceIp:     wf.SourceIPStr,
		SourceIf:     wf.SourceIf,
		DestIp:       wf.DestIPStr,
		DestPort:     wf.DestPort,
		DestIf:       wf.DestIf,
		SentByte:     wf.TrafficOut,
		ReceivedByte: wf.TrafficIn,
		Message:      &message,
		LogType:      logtype,
		Device:       dev,
		Level:        loglevel,
		User:         user,
		Service:      service,
	}
	err = txn.Insert(log)
	if err != nil {
		return err
	}

	profile, err := dal.GetOrInsertProfileByName(txn, wf.Profile)
	if err != nil {
		return err
	}

	fmt.Println("Done profile")
	status, err := dal.GetUtmstatusByName(txn, wf.Status)
	if err != nil {
		return err
	}
	if status == nil {
		return errors.New("Invalid UTM status value: " + wf.Status)
	}

	fmt.Println("Done utm status")
	category, err := dal.GetOrInsertCategoryByDescription(txn, wf.CategoryDesc)
	if err == nil {
		return err
	}

	fmt.Println("Done category")
	webfilter := &models.WebFilter{
		Host:     wf.Hostname,
		Url:      wf.Url,
		Log:      log,
		Profile:  profile,
		Status:   status,
		Category: category,
	}
	fmt.Println("WebFilter:", webfilter)
	err = txn.Insert(webfilter)
	fmt.Println("Error:", err)
	if err != nil {
		return err
	}

	return nil
}
