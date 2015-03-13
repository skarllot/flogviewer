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
	"errors"
	"github.com/go-gorp/gorp"
	"github.com/skarllot/flogviewer/models"
)

const (
	SQL_HOST_BYNAME = `SELECT id, name, category
	FROM host
	WHERE name = :name`
)

func GetHostByName(
	txn *gorp.Transaction,
	name string) (*models.Host, error) {

	qrows := make([]models.Host, 0)

	_, err := txn.Select(&qrows, SQL_HOST_BYNAME, map[string]interface{}{
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

func GetOrInsertHostByName(
	txn *gorp.Transaction,
	name string,
	category *models.Category) (*models.Host, error) {

	row, err := GetOrInsertByUnique(txn, SQL_HOST_BYNAME,
		&[]*models.Host{},
		map[string]interface{}{
			"name": name,
		}, func() interface{} {
			return &models.Host{
				Name:     name,
				Category: category,
			}
		})

	if err != nil {
		return nil, err
	}
	switch v := row.(type) {
	case *models.Host:
		return v, nil
	default:
		return nil, errors.New("Invalid type")
	}
}
