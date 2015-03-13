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
	"github.com/skarllot/flogviewer/common"
)

type Host struct {
	Id         int64  `db:"id"`
	Name       string `db:"name"`
	CategoryId *int64 `db:"category"`

	Category *Category `db:"-"`
}

func DefineHostTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(Host{}, "host")
	t.SetKeys(true, "id")
	t.ColMap("name").
		SetUnique(true).
		SetNotNull(true)
}

func (self *Host) PreInsert(gorp.SqlExecutor) error {
	if self.Category != nil {
		self.CategoryId = common.NInt64(self.Category.Id)
	}

	return nil
}

func (self *Host) PreUpdate(exe gorp.SqlExecutor) error {
	return self.PreInsert(exe)
}

func (self *Host) PostGet(exe gorp.SqlExecutor) error {
	if self.CategoryId != nil {
		obj, err := exe.Get(Category{}, *self.CategoryId)
		if err != nil {
			return err
		}
		self.Category = obj.(*Category)
	}
	return nil
}
