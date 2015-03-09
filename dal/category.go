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
	SQL_CATEGORY_BYNAME = `SELECT id, description
	FROM category
	WHERE description = :desc`
)

func GetCategoryByDescription(
	txn *gorp.Transaction,
	desc string) (*models.Category, error) {

	qrows := make([]models.Category, 0)

	_, err := txn.Select(&qrows, SQL_CATEGORY_BYNAME, map[string]interface{}{
		"desc": desc,
	})
	if err != nil {
		return nil, err
	}
	if len(qrows) != 1 {
		return nil, nil
	}

	return &qrows[0], nil
}

func GetOrInsertCategoryByDescription(
	txn *gorp.Transaction,
	desc string) (*models.Category, error) {

	row, err := GetOrInsertByUnique(txn, SQL_CATEGORY_BYNAME,
		&[]*models.Category{},
		map[string]interface{}{
			"desc": desc,
		}, func() interface{} {
			return &models.Category{
				Description: desc,
			}
		})

	if err != nil {
		return nil, err
	}
	switch v := row.(type) {
	case *models.Category:
		return v, nil
	default:
		return nil, errors.New("Invalid type")
	}
}
