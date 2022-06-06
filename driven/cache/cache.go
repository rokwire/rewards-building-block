// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cacheadapter

import (
	"github.com/patrickmn/go-cache"
	"rewards/core/model"
	"strconv"
	"time"
)

// CacheAdapter structure
type CacheAdapter struct {
	cache *cache.Cache
}

// NewCacheAdapter creates new instance
func NewCacheAdapter(defaultCacheExpirationSeconds string) *CacheAdapter {

	val, err := strconv.ParseInt(defaultCacheExpirationSeconds, 0, 64)
	var duration time.Duration
	if val == 0 || err != nil {
		duration = 1800 * time.Second
	} else {
		duration = time.Duration(val) * time.Second
	}

	cache := cache.New(duration, duration)

	return &CacheAdapter{
		cache: cache,
	}
}

// SetRewardTypes set reward types
func (s *CacheAdapter) SetRewardTypes(tips []model.RewardType) []model.RewardType {
	key := "reward_types"
	if tips == nil {
		s.cache.Delete(key)
	} else {
		s.cache.Set(key, tips, cache.DefaultExpiration)
	}
	return tips
}

// GetRewardTypes get all reward types
func (s *CacheAdapter) GetRewardTypes() []model.RewardType {
	obj, _ := s.cache.Get("reward_types")
	if obj != nil {
		return obj.([]model.RewardType)
	}
	return nil
}
