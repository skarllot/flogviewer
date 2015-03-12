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

package common

import (
	"math"
	"strings"
)

type OnMissedCache func(keys ...string) (result interface{}, err error)

type RoutineCache struct {
	cache  map[string]interface{}
	hits   map[string]int
	maxlen int
	missed OnMissedCache
}

func NewRoutineCache(size int, onMiss OnMissedCache) *RoutineCache {
	return &RoutineCache{
		cache:  make(map[string]interface{}, size),
		hits:   make(map[string]int, size),
		maxlen: size,
		missed: onMiss,
	}
}

func (self *RoutineCache) leastUsedKey() string {
	minHits := math.MaxUint32
	var minKey string

	for k, v := range self.hits {
		if v < minHits {
			minHits = v
			minKey = k
		}
	}
	return minKey
}

func (self *RoutineCache) removeLeastUsed() {
	key := self.leastUsedKey()
	delete(self.cache, key)
	delete(self.hits, key)
}

func (self *RoutineCache) Value(keys ...string) (interface{}, error) {
	var err error
	key := strings.Join(keys, "|")

	result, ok := self.cache[key]
	if ok {
		self.hits[key]++
		return result, nil
	} else {
		result, err = self.missed(keys...)
		if err != nil {
			return nil, err
		}

		if len(self.cache) == self.maxlen {
			self.removeLeastUsed()
		}
		self.cache[key] = result
		self.hits[key] = 1
		return result, nil
	}
}
