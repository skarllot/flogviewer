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
	defer txn.Rollback()
	for i := 0; i < count; i++ {
		item := <-c
		if item != nil {
			err = InsertWebFilter(*item, txn)
			if err != nil {
				return err
			}
		}
	}

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
	dev, err := dal.GetDeviceBySerial(txn, wf.DeviceSerial)
	if err != nil {
		return err
	}
	if dev == nil {
		dev = &models.Device{}
		dev.Name = wf.Device
		dev.Serial = wf.DeviceSerial
		err = txn.Insert(dev)
		if err != nil {
			return err
		}
	}

	logtype, err := dal.GetLogtypeByNames(txn, wf.LogType, wf.LogSubType)
	if err != nil {
		return err
	}
	if logtype == nil {
		logtype = &models.LogType{}
		logtype.Level1 = wf.LogType
		logtype.Level2 = wf.LogSubType
		err = txn.Insert(logtype)
		if err != nil {
			return err
		}
	}

	loglevel, err := dal.GetLoglevelByDesc(txn, wf.LogLevel)
	if err != nil {
		return err
	}
	if loglevel == nil {
		return errors.New("Invalid log level value: " + wf.LogLevel)
	}

	user, err := dal.GetUserByName(txn, wf.User)
	if err != nil {
		return err
	}
	if user == nil {
		user = &models.User{}
		user.Name = wf.User
		err = txn.Insert(user)
		if err != nil {
			return err
		}
	}

	service, err := dal.GetServiceByName(txn, wf.Service)
	if err != nil {
		return err
	}
	if service == nil {
		service = &models.Service{}
		service.Name = wf.Service
		err = txn.Insert(service)
		if err != nil {
			return err
		}
	}

	profile, err := dal.GetProfileByName(txn, wf.Profile)
	if err != nil {
		return err
	}
	if profile == nil {
		profile = &models.Profile{}
		profile.Name = wf.Profile
		err = txn.Insert(profile)
		if err != nil {
			return err
		}
	}

	status, err := dal.GetOrInsertUtmstatusByName(txn, wf.Status, func() *models.UtmStatus {
		return &models.UtmStatus{
			Name: wf.Status,
		}
	})
	if err != nil {
		return err
	}
	if status != status {
		return nil
	}

	return nil
}
