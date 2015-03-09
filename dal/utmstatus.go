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
	SQL_UTMSTATUS_BYNAME = `SELECT id, name
	FROM utmstatus
	WHERE name = :name`
)

func GetUtmstatusByName(
	txn *gorp.Transaction,
	name string) (*models.UtmStatus, error) {

	row, err := GetRowByUnique(txn, SQL_UTMSTATUS_BYNAME,
		&[]*models.UtmStatus{},
		map[string]interface{}{
			"name": name,
		})

	if err != nil {
		return nil, err
	}
	return row.(*models.UtmStatus), err
}
