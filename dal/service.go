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
	SQL_SERVICE_BYNAME = `SELECT id, name
	FROM service
	WHERE name = :name`
)

func GetServiceByName(txn *gorp.Transaction, name string) (*models.Service, error) {
	qrows := make([]models.Service, 0)

	_, err := txn.Select(&qrows, SQL_SERVICE_BYNAME, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, err
	}
	if len(qrows) != 1 {
		return nil, nil
	}

	return &qrows[0], nil
}

func GetOrInsertServiceByName(
	txn *gorp.Transaction,
	name string) (*models.Service, error) {

	row, err := GetOrInsertByUnique(txn, SQL_SERVICE_BYNAME,
		&[]*models.Service{},
		map[string]interface{}{
			"name": name,
		}, func() interface{} {
			return &models.Service{
				Name: name,
			}
		})

	if err != nil {
		return nil, err
	}
	return row.(*models.Service), nil
}
