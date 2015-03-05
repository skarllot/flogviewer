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
)

type LogType struct {
	Id        int64  `db:"id"`
	Name      string `db:"name"`
	SubTypeId *int64 `db:"subtype"`

	SubType *LogType `db:"-"`
}

func DefineLogtypeTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(LogType{}, "logtype")
	t.SetKeys(true, "id")
	t.ColMap("name").
		SetMaxSize(45).
		SetNotNull(true)
}

func (self *LogType) PreInsert(gorp.SqlExecutor) error {
	if self.SubType != nil {
		id := self.SubType.Id
		self.SubTypeId = &id
	}

	return nil
}

func (self *LogType) PostGet(exe gorp.SqlExecutor) error {
	if self.SubTypeId != nil {
		obj, err := exe.Get(LogType{}, self.SubTypeId)
		if err != nil {
			return err
		}
		self.SubType = obj.(*LogType)
	}

	return nil
}
