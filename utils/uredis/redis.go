package uredis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	redisPool *redis.Pool
)

// empty password is ""
func redisConn(ip, port, passwd string) (redis.Conn, error) {
	c, err := redis.Dial("tcp",
		ip+":"+port,
		redis.DialConnectTimeout(5*time.Second),
		redis.DialReadTimeout(1*time.Second),
		redis.DialWriteTimeout(1*time.Second),
		redis.DialPassword(passwd),
		redis.DialKeepAlive(1*time.Second),
	)
	return c, err
}

//pool
//ip,port,passwd
func newPool(ip, port, passwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 5, //
		// MaxActive:       18, //
		IdleTimeout:     240 * time.Second,
		MaxConnLifetime: 300 * time.Second,
		Dial:            func() (redis.Conn, error) { return redisConn(ip, port, passwd) },
	}
}

// InitPool ...
func InitPool(ip, port, passwd string) {
	redisPool = newPool(ip, port, passwd)
}

// GetRedisDB ...
func GetRedisDB() redis.Conn {
	return redisPool.Get()
}

// RedisScript ...
func RedisScript(s string, keynum int) *redis.Script {
	var rScript = redis.NewScript(keynum, s)
	r := GetRedisDB()
	rScript.Load(r)
	return rScript
}
