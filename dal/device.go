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
	SQL_DEVICE_BYSERIAL = `SELECT id, name, serial
	FROM device
	WHERE serial = :serial`
)

func GetDeviceBySerial(txn *gorp.Transaction, serial string) (*models.Device, error) {
	qrows := make([]models.Device, 0)

	_, err := txn.Select(&qrows, SQL_DEVICE_BYSERIAL, map[string]interface{}{
		"serial": serial,
	})
	if err != nil {
		return nil, err
	}

	if len(qrows) != 1 {
		return nil, nil
	}

	return &qrows[0], nil
}
