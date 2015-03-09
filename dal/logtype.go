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

package dal

import (
	"github.com/go-gorp/gorp"
	"github.com/skarllot/flogviewer/models"
)

const (
	SQL_LOGTYPE_BYNAME = `SELECT id, level1, level2
	FROM logtype
	WHERE level1 = :level1 AND level2 = :level2`
)

func GetLogtypeByNames(
	txn *gorp.Transaction,
	level1, level2 string) (*models.LogType, error) {
	qrows := make([]models.LogType, 0)

	_, err := txn.Select(&qrows, SQL_LOGTYPE_BYNAME, map[string]interface{}{
		"level1": level1,
		"level2": level2,
	})
	if err != nil {
		return nil, err
	}

	if len(qrows) != 1 {
		return nil, nil
	}

	return &qrows[0], nil
}

func GetOrInsertLogtypeByNames(
	txn *gorp.Transaction,
	level1, level2 string) (*models.LogType, error) {

	row, err := GetOrInsertByUnique(txn, SQL_LOGTYPE_BYNAME,
		&[]*models.LogType{},
		map[string]interface{}{
			"level1": level1,
			"level2": level2,
		}, func() interface{} {
			return &models.LogType{
				Level1: level1,
				Level2: level2,
			}
		})

	if err != nil {
		return nil, err
	}
	return row.(*models.LogType), nil
}
