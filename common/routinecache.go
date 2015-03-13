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

type OnMissedCache func(key []string, args []interface{}) (result interface{}, err error)

type RoutineCache struct {
	cache  map[string]interface{}
	hits   map[string]int
	maxlen int
	missed OnMissedCache
}

type RoutineCacheQuery struct {
	rCache *RoutineCache
	args   []string
	key    string
}

func NewRoutineCache(cacheLen int, onMiss OnMissedCache) *RoutineCache {
	return &RoutineCache{
		cache:  make(map[string]interface{}, cacheLen),
		hits:   make(map[string]int, cacheLen),
		maxlen: cacheLen,
		missed: onMiss,
	}
}

func (self *RoutineCache) Key(args ...string) *RoutineCacheQuery {
	return &RoutineCacheQuery{
		rCache: self,
		args:   args,
		key:    strings.Join(args, "|"),
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

func (self *RoutineCacheQuery) SetValue(value interface{}) {
	self.rCache.cache[self.key] = value
}

func (self *RoutineCacheQuery) Value(args ...interface{}) (interface{}, error) {
	var err error

	result, ok := self.rCache.cache[self.key]
	if ok {
		self.rCache.hits[self.key]++
		return result, nil
	} else {
		result, err = self.rCache.missed(self.args, args)
		if err != nil {
			return nil, err
		}

		if len(self.rCache.cache) == self.rCache.maxlen {
			self.rCache.removeLeastUsed()
		}
		self.rCache.cache[self.key] = result
		self.rCache.hits[self.key] = 1
		return result, nil
	}
}
