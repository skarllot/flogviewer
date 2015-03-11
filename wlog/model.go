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
	"net"
	"time"
)

type WebFilter struct {
	Date         time.Time
	DateStr      string `log:"date"`
	TimeStr      string `log:"time"`
	Device       string `log:"devname"`
	DeviceSerial string `log:"devid"`
	LogId        int64  `log:"logid"`
	LogType      string `log:"type"`
	LogSubType   string `log:"subtype"`
	LogLevel     string `log:"level"`
	PolicyId     int64  `log:"policyid"`
	SessionId    int64  `log:"sessionid"`
	User         string `log:"user"`
	SourceIP     net.IP
	SourceIPStr  string `log:"srcip"`
	SourceIf     string `log:"srcintf"`
	DestIP       net.IP
	DestIPStr    string `log:"dstip"`
	DestPort     uint16 `log:"dstport"`
	DestIf       string `log:"dstintf"`
	Service      string `log:"service"`
	Hostname     string `log:"hostname"`
	Profile      string `log:"profile"`
	Status       string `log:"status"`
	Url          string `log:"url"`
	TrafficOut   uint64 `log:"sentbyte"`
	TrafficIn    uint64 `log:"rcvdbyte"`
	Message      string `log:"msg"`
	CategoryId   int    `log:"cat"`
	CategoryDesc string `log:"catdesc"`
}

type WebFilterList []WebFilter

func (wf *WebFilter) ConvertFields() {
	wf.Date, _ = time.Parse("2006-01-02 15:04:05", wf.DateStr+" "+wf.TimeStr)
	wf.SourceIP = net.ParseIP(wf.SourceIPStr)
	wf.DestIP = net.ParseIP(wf.DestIPStr)
}
