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
	"bufio"
	"github.com/skarllot/flogviewer/common"
	"os"
)

func ParseFile(fname string) WebFilterList {
	result := make(WebFilterList, 0)

	file, err := os.Open(fname)
	if err != nil {
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	c := make(chan *WebFilter)
	count := 0
	for scanner.Scan() {
		go ParseLine(scanner.Text(), c)
		count++
	}

	for i := 0; i < count; i++ {
		item := <-c
		if item != nil {
			result = append(result, *item)
		}
	}

	return result
}

func ParseLine(line string, c chan *WebFilter) {
	wf := &WebFilter{}
	if err := common.ParseKeyValueLog(line, wf); err != nil {
		c <- nil
		return
	}

	wf.ConvertFields()
	c <- wf
}
