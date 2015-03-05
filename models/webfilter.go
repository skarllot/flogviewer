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

type WebFilter struct {
	Id         int64  `db:"id"`
	ProfileId  int64  `db:"profile"`
	StatusId   int64  `db:"status"`
	CategoryId int64  `db:"category"`
	Host       string `db:"host"`
	Url        string `db:"url"`
	Message    string `db:"message"`

	Log      *Log       `db:"-"`
	Profile  *Profile   `db:"-"`
	Status   *UtmStatus `db:"-"`
	Category *Category  `db:"-"`
}

func DefineWebfilterTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(WebFilter{}, "webfilter")
	t.SetKeys(false, "id")
	t.ColMap("host").
		SetMaxSize(255).
		SetNotNull(true)
	t.ColMap("url").
		SetMaxSize(2048).
		SetNotNull(true)
	t.ColMap("message").
		SetNotNull(true)
}

func (self *WebFilter) PreInsert(gorp.SqlExecutor) error {
	if self.Log != nil {
		self.Id = self.Log.Id
	}
	if self.Profile != nil {
		self.ProfileId = self.Profile.Id
	}
	if self.Status != nil {
		self.StatusId = self.Status.Id
	}
	if self.Category != nil {
		self.CategoryId = self.Category.Id
	}

	return nil
}

func (self *WebFilter) PostGet(exe gorp.SqlExecutor) error {
	obj, err := exe.Get(Log{}, self.Id)
	if err != nil {
		return err
	}
	self.Log = obj.(*Log)

	obj, err = exe.Get(Profile{}, self.ProfileId)
	if err != nil {
		return err
	}
	self.Profile = obj.(*Profile)

	obj, err = exe.Get(UtmStatus{}, self.StatusId)
	if err != nil {
		return err
	}
	self.Status = obj.(*UtmStatus)

	obj, err = exe.Get(Category{}, self.CategoryId)
	if err != nil {
		return err
	}
	self.Category = obj.(*Category)
	return nil
}
