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

type File struct {
	Id    int64 `db:"id"`
	Begin int64 `db:"begin_dt"`
	End   int64 `db:"end_dt"`
	Count int64 `db:"count_lines"`
}

func DefineFileTable(dbm *gorp.DbMap) {
	t := dbm.AddTableWithName(File{}, "file")
	t.SetKeys(true, "id")
	t.ColMap("begin_dt").
		SetUnique(true).
		SetNotNull(true)
	t.ColMap("end_dt").
		SetNotNull(true)
	t.ColMap("count_lines").
		SetNotNull(true)
}
