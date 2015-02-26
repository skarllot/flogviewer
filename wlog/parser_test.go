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

package wlog

import (
	"encoding/json"
	"testing"
	"time"
)

const (
	TEST_LINE = `date=2015-02-21 time=23:41:33 devname=TST-FGT01 devid=FG200D4614802240 logid=0317013312 type=utm subtype=webfilter eventtype=ftgd_allow level=notice vd="root" policyid=33 identidx=2 sessionid=623540870 user="HFONSECA" srcip=10.1.1.183 srcport=3204 srcintf="port6" dstip=216.58.222.2 dstport=80 dstintf="wan1" service="http" hostname="pagead2.googlesyndication.com" profiletype="Webfilter_Profile" profile="G_Web_Medium" status="passthrough" reqtype="referral" url="/activeview?id=osdim&avi=BXh_6zUHpVJebH4LNfP3xgpgHAAAAABABOAHIAQTAAgLgAgDgBAGgBhXCEwMQgAE&ti=1&adk" sentbyte=1310 rcvdbyte=506 msg="URL belongs to an allowed category in policy" method=domain class=0 cat=17 catdesc="Advertising"`
)

func TestParseLine(t *testing.T) {
	c := make(chan *WebFilter)
	go ParseLine(TEST_LINE, c)
	result := <-c
	if result == nil {
		t.Fatal("Could not parse log line")
		return
	}

	dateTest, _ := time.Parse("2006-01-02 15:04:05", "2015-02-21 23:41:33")
	if result.Date != dateTest {
		t.Error("The parsed date do not match")
	}
	if result.Device != "TST-FGT01" {
		t.Error("The parsed device name do not match")
	}
	if result.LogLevel != "notice" {
		t.Error("The parsed log level do not match")
	}
	if result.PolicyId != 33 {
		t.Error("The parsed policy ID do not match")
	}
	if result.CategoryId != 17 {
		t.Error("The parsed category ID do not match")
	}
	if result.CategoryDesc != "Advertising" {
		t.Error("The parsed category description do not match")
	}

	out, _ := json.Marshal(result)
	t.Logf("Parsed line: %s", out)
}
