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
	"fmt"
	"github.com/go-gorp/gorp"
)

const (
	SQL_MODEL_COUNT = `SELECT count(%s) AS count FROM %s`
)

func CountRows(tableName, columnName string, txn *gorp.Transaction) (int64, error) {
	count, err := txn.SelectInt(fmt.Sprintf(SQL_MODEL_COUNT, columnName, tableName))
	if err != nil {
		return -1, err
	}

	return count, nil
}
