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
	"time"
)

type Log struct {
	Id           int64     `db:"id"`
	FileId       int64     `db:"file"`
	LogtypeId    int64     `db:"logtype"`
	DeviceId     int64     `db:"device"`
	LevelId      int64     `db:"level"`
	UserId       int64     `db:"user"`
	ServiceId    int64     `db:"service"`
	LogId        int64     `db:"log_id"`
	Date         time.Time `db:"date"`
	SessionId    int64     `db:"session_id"`
	PolicyId     int64     `db:"policy_id"`
	SourceIp     string    `db:"source_ip"`
	SourceIf     string    `db:"source_if"`
	DestIp       string    `db:"dest_ip"`
	DestPort     uint16    `db:"dest_port"`
	DestIf       string    `db:"dest_if"`
	SentByte     uint64    `db:"sent_byte"`
	ReceivedByte uint64    `db:"received_byte"`
	Message      *string   `db:"message"`

	File    *File     `db:"-"`
	LogType *LogType  `db:"-"`
	Device  *Device   `db:"-"`
	Level   *LogLevel `db:"-"`
	User    *User     `db:"-"`
	Service *Service  `db:"-"`
}

func DefineLogTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(Log{}, "log")
	t.SetKeys(true, "id")
	SetNotNull(t,
		"file", "logtype", "device", "level", "user", "service", "date",
		"policy_id", "source_if", "dest_port", "dest_if", "sent_byte",
		"received_byte")
	t.ColMap("source_ip").SetMaxSize(45).SetNotNull(true)
	t.ColMap("dest_ip").SetMaxSize(45).SetNotNull(true)
	t.ColMap("message").SetMaxSize(255)
}

func (self *Log) PreInsert(gorp.SqlExecutor) error {
	if self.File != nil {
		self.FileId = self.File.Id
	}
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
	obj, err := exe.Get(File{}, self.FileId)
	if err != nil {
		return err
	}
	self.File = obj.(*File)

	obj, err = exe.Get(LogType{}, self.LogtypeId)
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
