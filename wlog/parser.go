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
	pool     chan interface{}
	get      chan []string
	response chan interface{}
	quit     chan bool
}

func ParseFile(r io.Reader, dbm *gorp.DbMap) error {
	scanner := bufio.NewScanner(r)
	chanLine := make(chan *WebFilter)
	countLines := 0
	defer close(chanLine)
	for scanner.Scan() {
		go ParseLine(scanner.Text(), chanLine)
		countLines++
	}

	batch := make([]*WebFilter, 0, PARSE_BATCH_SIZE)
	countBatch := 0
	chanBatch := make(chan int)
	batcher := &ParserBatch{
		dbm:      dbm,
		pool:     make(chan interface{}, PARSE_CONNECTION_POOL),
		get:      make(chan []string),
		response: make(chan interface{}),
		quit:     make(chan bool),
	}
	for i := 0; i < PARSE_CONNECTION_POOL; i++ {
		batcher.pool <- 0
	}
	go batcher.ForeignTableGet()

	for i := 0; i < countLines; i++ {
		item := <-chanLine
		if item != nil {
			batch = append(batch, item)
			if len(batch) == PARSE_BATCH_SIZE {
				go batcher.InsertWebFilterList(batch, chanBatch)
				countBatch++
				batch = make([]*WebFilter, 0, PARSE_BATCH_SIZE)
			}
		}
	}
	if len(batch) > 0 {
		go batcher.InsertWebFilterList(batch, chanBatch)
		countBatch++
	}

	fmt.Print("Records inserted: ")
	countRecords, lastPrintedCount := 0, 0
	for i := 0; i < countBatch; i++ {
		batchResult := <-chanBatch
		if batchResult > 0 {
			countRecords += batchResult
		}
		if countRecords >= lastPrintedCount+1000 {
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
	lock := <-self.pool
	defer func() { self.pool <- lock }()

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
		Host: wf.Hostname,
		Url:  wf.Url,
		Log:  log,
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

	self.get <- []string{"category", wf.CategoryDesc}
	fResult = <-self.response
	switch t := fResult.(type) {
	case *models.Category:
		webfilter.Category = t
	case error:
		return t
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

	deviceCache := make(map[string]*models.Device, 0)
	logtypeCache := make(map[string]*models.LogType, 0)
	userCache := make(map[string]*models.User, 0)
	serviceCache := make(map[string]*models.Service, 0)
	profileCache := make(map[string]*models.Profile, 0)
	categoryCache := make(map[string]*models.Category, 0)

	for {
		select {
		case req := <-self.get:
			var fResult interface{}
			var ok bool
			err = nil
			switch req[0] {
			case "device":
				fResult, ok = deviceCache[req[1]]
				if !ok {
					fResult, err = dal.GetOrInsertDeviceBySerial(txn, req[1], req[2])
					deviceCache[req[1]] = fResult.(*models.Device)
				}
			case "logtype":
				fResult, ok = logtypeCache[req[1]+req[2]]
				if !ok {
					fResult, err = dal.GetOrInsertLogtypeByNames(txn, req[1], req[2])
					logtypeCache[req[1]+req[2]] = fResult.(*models.LogType)
				}
			case "user":
				fResult, ok = userCache[req[1]]
				if !ok {
					fResult, err = dal.GetOrInsertUserByName(txn, req[1])
					userCache[req[1]] = fResult.(*models.User)
				}
			case "service":
				fResult, ok = serviceCache[req[1]]
				if !ok {
					fResult, err = dal.GetOrInsertServiceByName(txn, req[1])
					serviceCache[req[1]] = fResult.(*models.Service)
				}
			case "profile":
				fResult, ok = profileCache[req[1]]
				if !ok {
					fResult, err = dal.GetOrInsertProfileByName(txn, req[1])
					profileCache[req[1]] = fResult.(*models.Profile)
				}
			case "category":
				fResult, ok = categoryCache[req[1]]
				if !ok {
					fResult, err = dal.GetOrInsertCategoryByDescription(txn, req[1])
					categoryCache[req[1]] = fResult.(*models.Category)
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
