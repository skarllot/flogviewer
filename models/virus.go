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

type Virus struct {
	Id         int64  `db:"id"`
	StatusId   int64  `db:"status"`
	VirusDefId int64  `db:"virusdef"`
	Url        string `db:"url"`

	Log      *Log      `db:"-"`
	VirusDef *VirusDef `db:"-"`
}

func DefineVirusTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(Virus{}, "virus")
	t.SetKeys(false, "id")
	SetNotNull(t, "status", "virusdef")
	t.ColMap("url").
		SetMaxSize(2048).
		SetNotNull(true)
}

func (self *Virus) PreInsert(gorp.SqlExecutor) error {
	if self.Log != nil {
		self.Id = self.Log.Id
	}
	if self.VirusDef != nil {
		self.VirusDefId = self.VirusDef.Id
	}

	return nil
}

func (self *Virus) PostGet(exe gorp.SqlExecutor) error {
	obj, err := exe.Get(Log{}, self.Id)
	if err != nil {
		return err
	}
	self.Log = obj.(*Log)

	obj, err = exe.Get(VirusDef{}, self.VirusDefId)
	if err != nil {
		return err
	}
	self.VirusDef = obj.(*VirusDef)

	return nil
}
