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
	"reflect"
)

func GetRowByUnique(
	txn *gorp.Transaction,
	query string,
	qrows interface{},
	unique map[string]interface{}) (interface{}, error) {

	_, err := txn.Select(qrows, query, unique)
	if err != nil {
		return nil, err
	}

	rowsSlice := MakeGenericSlice(qrows)
	if len(rowsSlice) != 1 {
		return nil, nil
	}

	return rowsSlice[0], nil
}

func GetOrInsertByUnique(
	txn *gorp.Transaction,
	query string,
	qrows interface{},
	unique map[string]interface{},
	insert func() interface{}) (interface{}, error) {

	row, err := GetRowByUnique(txn, query, qrows, unique)
	if err != nil {
		return nil, err
	}

	if row == nil {
		row = insert()
		err = txn.Insert(row)
		if err != nil {
			return nil, err
		}
	}

	return row, nil
}

func MakeGenericSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
