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

package bll

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-gorp/gorp"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/skarllot/flogviewer/models"
	_ "github.com/ziutek/mymysql/godrv"
)

type Configuration struct {
	Database
}

type Database struct {
	Engine   string
	Host     string
	Port     uint16
	Name     string
	User     string
	Password string
	Protocol string
	DbArgs   string
}

func (c *Configuration) Load(path string) error {
	c.loadDefaults()
	err := gcfg.ReadFileInto(c, path)
	return err
}

func (c *Configuration) loadDefaults() {
	c.Database.Engine = "mysql"
	c.Database.Host = "localhost"
	c.Database.Port = 3306
	c.Database.Name = "flogviewer"
	c.Database.User = "flogviewer"
	c.Database.Password = ""
	c.Database.Protocol = "tcp"
	c.Database.DbArgs = ""
}

func (db *Database) GetConnectionString() (string, error) {
	switch db.Engine {
	case "mysql":
		return db.getMysqlConnectionString(), nil
	case "mymysql":
		return db.getMymysqlConnectionString(), nil
	default:
		if len(db.Engine) == 0 {
			return "", errors.New(
				"No engine name defined into configuration file")
		} else {
			return "", errors.New(
				fmt.Sprintf("The engine '%s' was not implemented", db.Engine))
		}

	}
}

func (db *Database) getMysqlConnectionString() string {
	if len(db.DbArgs) > 0 {
		db.DbArgs = "?" + db.DbArgs
	}

	return fmt.Sprintf(
		"%s:%s@%s([%s]:%d)/%s%s",
		db.User, db.Password, db.Protocol, db.Host, db.Port, db.Name, db.DbArgs)
}

func (db *Database) getMymysqlConnectionString() string {
	if len(db.DbArgs) > 0 {
		db.DbArgs = "," + db.DbArgs
	}

	return fmt.Sprintf(
		"%s:%s:%d%s*%s/%s/%s",
		db.Protocol, db.Host, db.Port, db.DbArgs, db.Name, db.User, db.Password)
}

func (self *Configuration) CreateDbMap() (*gorp.DbMap, error) {
	engine := self.Database.Engine
	cnxStr, err := self.GetConnectionString()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(engine, cnxStr)
	if err != nil {
		return nil, err
	}

	var dbm *gorp.DbMap
	switch engine {
	case "mysql":
	case "mymysql":
		dbm = &gorp.DbMap{
			Db:      db,
			Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
		}
	default:
		return nil, errors.New(fmt.Sprintf(
			"The engine '%s' is not implemented", engine))
	}

	if err := models.InitAllTables(dbm); err != nil {
		return nil, err
	}

	return dbm, nil
}
