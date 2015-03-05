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

package models

import (
	"github.com/go-gorp/gorp"
	"net"
	"time"
)

type Log struct {
	Id           int64     `db:"id"`
	LogtypeId    int64     `db:"logtype"`
	DeviceId     int64     `db:"device"`
	LevelId      int64     `db:"level"`
	UserId       int64     `db:"user"`
	ServiceId    int64     `db:"service"`
	Date         time.Time `db:"date"`
	PolicyId     int64     `db:"policy_id"`
	SourceIp     net.IP    `db:"source_ip"`
	SourceIf     string    `db:"source_if"`
	DestIp       net.IP    `db:"dest_ip"`
	DestPort     int16     `db:"dest_port"`
	DestIf       string    `db:"dest_if"`
	SentByte     int64     `db:"sent_byte"`
	ReceivedByte int64     `db:"received_byte"`
	Message      *string   `db:"message"`

	LogType *LogType  `db:"-"`
	Device  *Device   `db:"-"`
	Level   *LogLevel `db:"-"`
	User    *User     `db:"-"`
	Service *Service  `db:"-"`
}

func DefineLogTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(Log{}, "log")
	t.SetKeys(true, "id")
	t.ColMap("logtype").SetNotNull(true)
	t.ColMap("device").SetNotNull(true)
	t.ColMap("level").SetNotNull(true)
	t.ColMap("user").SetNotNull(true)
	t.ColMap("service").SetNotNull(true)
	t.ColMap("date").SetNotNull(true)
	t.ColMap("policy_id").SetNotNull(true)
	t.ColMap("source_ip").SetMaxSize(45).SetNotNull(true)
	t.ColMap("source_if").SetNotNull(true)
	t.ColMap("dest_ip").SetMaxSize(45).SetNotNull(true)
	t.ColMap("dest_port").SetNotNull(true)
	t.ColMap("dest_if").SetNotNull(true)
	t.ColMap("sent_byte").SetNotNull(true)
	t.ColMap("received_byte").SetNotNull(true)
	t.ColMap("message").SetMaxSize(255)
}

func (self *Log) PreInsert(gorp.SqlExecutor) error {
	if self.LogType != nil {
		self.LogtypeId = self.LogType.Id
	}
	if self.Device != nil {
		self.DeviceId = self.Device.Id
	}
	if self.Level != nil {
		self.LevelId = self.Level.Id
	}
	if self.User != nil {
		self.UserId = self.User.Id
	}
	if self.Service != nil {
		self.ServiceId = self.Service.Id
	}

	return nil
}

func (self *Log) PostGet(exe gorp.SqlExecutor) error {
	obj, err := exe.Get(LogType{}, self.LogtypeId)
	if err != nil {
		return err
	}
	self.LogType = obj.(*LogType)

	obj, err = exe.Get(Device{}, self.DeviceId)
	if err != nil {
		return err
	}
	self.Device = obj.(*Device)

	obj, err = exe.Get(LogLevel{}, self.LevelId)
	if err != nil {
		return err
	}
	self.Level = obj.(*LogLevel)

	obj, err = exe.Get(User{}, self.UserId)
	if err != nil {
		return err
	}
	self.User = obj.(*User)

	obj, err = exe.Get(Service{}, self.ServiceId)
	if err != nil {
		return err
	}
	self.Service = obj.(*Service)
	return nil
}
