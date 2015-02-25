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
	"os"
	"reflect"
	"regexp"
	"strconv"
)

const (
	LINE_KEY_PATTERN = `(\w+)=(?:([0-9a-zA-Z-:_.]+)|(?:["]([^"]+)["]))`
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
	r, _ := regexp.Compile(LINE_KEY_PATTERN)
	count := 0
	for scanner.Scan() {
		go ParseLine(scanner.Text(), r, c)
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

func ParseLine(line string, r *regexp.Regexp, c chan *WebFilter) {
	keys := r.FindAllStringSubmatch(line, -1)
	if keys == nil {
		c <- nil
	}

	list := make(map[string]string, 0)
	for _, k := range keys {
		name := k[1]
		value := k[2]
		if k[3] != "" {
			value = k[3]
		}
		list[name] = value
	}

	wf := &WebFilter{}
	wftype := reflect.TypeOf(*wf)
	wfref := reflect.ValueOf(wf).Elem()
	for i := 0; i < wftype.NumField(); i++ {
		field := wftype.Field(i)
		fieldref := wfref.Field(i)
		tag := field.Tag.Get("log")
		if lVal, ok := list[tag]; ok {
			switch fieldref.Kind() {
			case reflect.String:
				fieldref.SetString(lVal)
			case reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64:
				intVal, _ := strconv.Atoi(lVal)
				fieldref.SetInt(int64(intVal))
			case reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64:
				intVal, _ := strconv.Atoi(lVal)
				fieldref.SetUint(uint64(intVal))
			}
		}
	}

	wf.ConvertFields()
	c <- wf
}
