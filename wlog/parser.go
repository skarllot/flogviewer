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
	"strings"
)

const (
	PARSE_BATCH_SIZE                 = 100
	PARSE_CONNECTION_POOL            = 15
	MYSQL_ERROR_TOO_MANY_CONNECTIONS = "1040"
	MYSQL_ERROR_DUPLICATE_ENTRY      = "1062"
	MYSQL_ERROR_DEADLOCK             = "1213"
)

type ParserBatch struct {
	dbm      *gorp.DbMap
	file     *models.File
	get      chan []string
	response chan interface{}
	quit     chan bool
}

func ParseFile(r io.Reader, fileRow *models.File, dbm *gorp.DbMap) error {
	chanLine := make(chan *WebFilter)
	jobScan := make(chan bool, PARSE_BATCH_SIZE)

	scanner := bufio.NewScanner(r)
	defer close(chanLine)
	go func() {
		var linesCount int64 = 0
		for scanner.Scan() {
			linesCount++
			if linesCount <= fileRow.Count {
				continue
			}

			go ParseLine(scanner.Text(), chanLine)
			jobScan <- true
		}
		fileRow.Count = linesCount
		close(jobScan)
	}()

	batch := make([]*WebFilter, 0, PARSE_BATCH_SIZE)
	chanBatch := make(chan int)
	batcher := &ParserBatch{
		dbm:      dbm,
		file:     fileRow,
		get:      make(chan []string),
		response: make(chan interface{}),
		quit:     make(chan bool),
	}
	go batcher.ForeignTableGet()

	jobBatch := make(chan bool, PARSE_CONNECTION_POOL)
	go func() {
		for _ = range jobScan {
			item := <-chanLine
			if item != nil {
				batch = append(batch, item)
				if len(batch) == PARSE_BATCH_SIZE {
					go batcher.InsertWebFilterList(batch, chanBatch)
					jobBatch <- true
					batch = make([]*WebFilter, 0, PARSE_BATCH_SIZE)
				}
			}
		}
		if len(batch) > 0 {
			go batcher.InsertWebFilterList(batch, chanBatch)
			jobBatch <- true
		}
		close(jobBatch)
	}()

	fmt.Print("Records inserted: ")
	countRecords, lastPrintedCount := 0, 0
	for _ = range jobBatch {
		batchResult := <-chanBatch
		if batchResult > 0 {
			countRecords += batchResult
		}
		if countRecords >= lastPrintedCount+10000 {
			fmt.Printf("%d ", countRecords)
			lastPrintedCount = countRecords
		}
	}
	if lastPrintedCount != countRecords {
		fmt.Printf("%d", countRecords)
	}
	fmt.Println()

	batcher.quit <- true
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

func (self *ParserBatch) InsertWebFilterList(list []*WebFilter, c chan<- int) {
	txn, err := self.dbm.Begin()
	if err != nil {
		fmt.Println("Error begining transaction:", err)
		c <- 0
		return
	}
	defer txn.Rollback()

	insertCount := 0
	for _, i := range list {
		err = self.InsertWebFilter(*i, txn)
		if err != nil {
			fmt.Println("Error inserting new record:", err)
		} else {
			insertCount++
		}
	}
	txn.Commit()
	c <- insertCount
}

func (self *ParserBatch) InsertWebFilter(wf WebFilter, txn *gorp.Transaction) error {
	log := &models.Log{
		LogId:        wf.LogId,
		Date:         wf.Date,
		SessionId:    wf.SessionId,
		PolicyId:     wf.PolicyId,
		SourceIp:     wf.SourceIPStr,
		SourceIf:     wf.SourceIf,
		DestIp:       wf.DestIPStr,
		DestPort:     wf.DestPort,
		DestIf:       wf.DestIf,
		SentByte:     wf.TrafficOut,
		ReceivedByte: wf.TrafficIn,
		Message:      common.NString(wf.Message),
		File:         self.file,
	}

	self.get <- []string{"device", wf.DeviceSerial, wf.Device}
	fResult := <-self.response
	switch t := fResult.(type) {
	case *models.Device:
		log.Device = t
	case error:
		return t
	}

	self.get <- []string{"logtype", wf.LogType, wf.LogSubType}
	fResult = <-self.response
	switch t := fResult.(type) {
	case *models.LogType:
		log.LogType = t
	case error:
		return t
	}

	if loglevel, err := dal.GetLoglevelByName(txn, wf.LogLevel); err != nil {
		return err
	} else if loglevel == nil {
		return errors.New("Invalid log level value: " + wf.LogLevel)
	} else {
		log.Level = loglevel
	}

	self.get <- []string{"user", strings.ToLower(wf.User)}
	fResult = <-self.response
	switch t := fResult.(type) {
	case *models.User:
		log.User = t
	case error:
		return t
	}

	self.get <- []string{"service", wf.Service}
	fResult = <-self.response
	switch t := fResult.(type) {
	case *models.Service:
		log.Service = t
	case error:
		return t
	}

	if err := txn.Insert(log); err != nil {
		return err
	}

	webfilter := &models.WebFilter{
		Log: log,
		Url: wf.Url,
	}

	self.get <- []string{"profile", wf.Profile}
	fResult = <-self.response
	switch t := fResult.(type) {
	case *models.Profile:
		webfilter.Profile = t
	case error:
		return t
	}

	if status, err := dal.GetUtmstatusByName(txn, wf.Status); err != nil {
		return err
	} else if status == nil {
		return errors.New("Invalid UTM status value: " + wf.Status)
	} else {
		webfilter.Status = status
	}

	self.get <- []string{"host", wf.Hostname, wf.CategoryDesc}
	fResult = <-self.response
	switch t := fResult.(type) {
	case *models.Host:
		webfilter.Host = t
	case error:
		return t
	}

	if len(wf.CategoryDesc) > 0 {
		if wf.CategoryDesc == webfilter.Host.Category.Description {
			webfilter.Category = webfilter.Host.Category
		} else {
			self.get <- []string{"category", wf.CategoryDesc}
			fResult = <-self.response
			switch t := fResult.(type) {
			case *models.Category:
				webfilter.Category = t
			case error:
				return t
			}
		}
	}

	if err := txn.Insert(webfilter); err != nil {
		return err
	}

	return nil
}

func (self *ParserBatch) ForeignTableGet() {
	countOp := 0
	txn, err := self.dbm.Begin()
	if err != nil {
		fmt.Println("Error begining transaction:", err)
		return
	}

	deviceCache := common.NewRoutineCache(
		5, func(key []string, args []interface{}) (interface{}, error) {
			return dal.GetOrInsertDeviceBySerial(txn, key[0], args[0].(string))
		})
	logtypeCache := common.NewRoutineCache(
		5, func(key []string, args []interface{}) (interface{}, error) {
			return dal.GetOrInsertLogtypeByNames(txn, key[0], key[1])
		})
	userCache := common.NewRoutineCache(
		100, func(key []string, args []interface{}) (interface{}, error) {
			return dal.GetOrInsertUserByName(txn, key[0])
		})
	serviceCache := common.NewRoutineCache(
		5, func(key []string, args []interface{}) (interface{}, error) {
			return dal.GetOrInsertServiceByName(txn, key[0])
		})
	profileCache := common.NewRoutineCache(
		20, func(key []string, args []interface{}) (interface{}, error) {
			return dal.GetOrInsertProfileByName(txn, key[0])
		})
	categoryCache := common.NewRoutineCache(
		80, func(key []string, args []interface{}) (interface{}, error) {
			return dal.GetOrInsertCategoryByDescription(txn, key[0])
		})
	hostCache := common.NewRoutineCache(
		1000, func(key []string, args []interface{}) (interface{}, error) {
			return dal.GetOrInsertHostByName(
				txn, key[0], args[0].(*models.Category))
		})

	for {
		select {
		case req := <-self.get:
			var fResult interface{} = nil
			err = nil
			switch req[0] {
			case "device": // 1=serial, 2=device
				fResult, err = deviceCache.Key(req[1]).Value(req[2])
			case "logtype": // 1=level1, 2=level2
				fResult, err = logtypeCache.Key(req[1], req[2]).Value()
			case "user": // 1=name
				fResult, err = userCache.Key(req[1]).Value()
			case "service": // 1=name
				fResult, err = serviceCache.Key(req[1]).Value()
			case "profile": // 1=name
				fResult, err = profileCache.Key(req[1]).Value()
			case "category": // 1=description
				if len(req[1]) > 0 {
					fResult, err = categoryCache.Key(req[1]).Value()
				}
			case "host": // 1=host.name, 2=category.name
				var cat *models.Category = nil
				var host *models.Host = nil
				if len(req[2]) > 0 {
					fResult, err = categoryCache.Key(req[2]).Value()
					if err == nil {
						cat = fResult.(*models.Category)
					}
				}
				if err == nil {
					fResult, err = hostCache.Key(req[1]).Value(cat)
				}
				if err == nil && cat != nil {
					host = fResult.(*models.Host)
					if host.Category == nil {
						host.Category = cat
						_, err = txn.Update(host)
						if err != nil {
							hostCache.Key(req[2]).SetValue(host)
						}
					}
				}
			default:
				err = errors.New("Invalid foreign table name")
			}

			if err != nil {
				self.response <- err
			} else {
				self.response <- fResult
			}

			countOp++
			if countOp >= PARSE_BATCH_SIZE {
				txn.Commit()
				txn, err = self.dbm.Begin()
				if err != nil {
					fmt.Println("Error begining transaction:", err)
					return
				}
				countOp = 0
			}
		case <-self.quit:
			txn.Commit()
			return
		}
	}
}
