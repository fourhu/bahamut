// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"fmt"

	"gopkg.in/redis.v3"
)

// RedisInfo represents Redis connection information
type RedisInfo struct {
	Addresses   []string
	Password    string
	DBNumber    int64
	ClusterName string
}

// NewRedisInfo returns a new RedisInfo
func NewRedisInfo(addresses []string, password string, db int64, clusterName string) *RedisInfo {

	return &RedisInfo{
		Addresses:   addresses,
		Password:    password,
		ClusterName: clusterName,
		DBNumber:    db,
	}
}

// IsSentinelModeActive checks if the current redis info want to use redis Sentinel
func (r *RedisInfo) IsSentinelModeActive() bool {

	return r.ClusterName != ""
}

func (r *RedisInfo) makeRedisClient() *redis.Client {

	if r.IsSentinelModeActive() {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    r.ClusterName,
			SentinelAddrs: r.Addresses,
			Password:      r.Password,
			DB:            r.DBNumber,
		})
	}

	return redis.NewClient(&redis.Options{
		Addr:     r.Addresses[0],
		Password: r.Password,
		DB:       r.DBNumber,
	})
}

func (r *RedisInfo) String() string {

	if r.IsSentinelModeActive() {
		return fmt.Sprintf("<redis clusterName: %s addresses: %v db: %d>", r.ClusterName, r.Addresses, r.DBNumber)
	}

	return fmt.Sprintf("<redis address: %s db: %d>", r.Addresses[0], r.DBNumber)
}