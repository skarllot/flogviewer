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

type LogLevel struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

func DefineLoglevelTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(LogLevel{}, "loglevel")
	t.SetKeys(true, "id")
	t.ColMap("name").
		SetMaxSize(15).
		SetNotNull(true)
}

func InitLoglevelTable(txn *gorp.Transaction) {
	rows := []*LogLevel{
		&LogLevel{0, "emergency"},
		&LogLevel{0, "alert"},
		&LogLevel{0, "critical"},
		&LogLevel{0, "error"},
		&LogLevel{0, "warning"},
		&LogLevel{0, "notice"},
		&LogLevel{0, "information"},
		&LogLevel{0, "debug"},
	}

	for _, r := range rows {
		if err := txn.Insert(r); err != nil {
			panic(err)
		}
	}
}
