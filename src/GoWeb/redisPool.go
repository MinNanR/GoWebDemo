package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisConnection struct {
	rs redis.Conn
}

/*
获取值
*/
func (rc RedisConnection) getValue(key string) string {
	rs := rc.rs
	value, err := redis.String(rs.Do("get", key))
	if err != nil {
		return ""
	}
	return value
}

/*
设置值
*/
func (rc RedisConnection) setValue(key string, value string) {
	rs := rc.rs
	rs.Do("set", key, value)
}

/*
将对象存入redis
*/
func (rc RedisConnection) setObjectValue(key string, obj interface{}) {
	rs := rc.rs
	jsonByte, _ := json.Marshal(obj)
	rs.Do("set", key, jsonByte)
}

/*
从redis中取出对象
*/
func (rc RedisConnection) getObjectValue(key string, obj *interface{}) {
	rs := rc.rs
	jsonByte, _ := redis.Bytes(rs.Do("get", key))
	json.Unmarshal(jsonByte, obj)
}

/*
设置key的过期时间
*/
func (rc RedisConnection) expire(key string, duration time.Duration) {
	rs := rc.rs
	rs.Do("expire", key, duration.Seconds())
}

/*
存入key并设置过期时间
*/
func (rc RedisConnection) setValueExpire(key string, value string, duration time.Duration) {
	rc.setValue(key, value)
	rc.expire(key, duration)
}

/*
存入对象并设置过期时间
*/
func (rc RedisConnection) setObjectValueExpire(key string, value interface{}, duration time.Duration) {
	rc.setObjectValue(key, value)
	rc.expire(key, duration)
}

var rc = RedisConnection{rs: connectionRedis()}

func connectionRedis() redis.Conn {
	host, port, database, password := redisConfig.Host, redisConfig.Port, redisConfig.Database, redisConfig.Password
	rs, _ := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	rs.Do("auth", password)
	rs.Do("select", database)
	return rs
}
