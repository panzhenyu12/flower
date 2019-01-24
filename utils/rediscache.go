package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc32"

	"github.com/go-redis/redis"
	"github.com/go-xorm/core"
	// "log"
	"reflect"
	// "strconv"
	"time"
	"unsafe"
)

const (
	DEFAULT_EXPIRATION = time.Duration(0)
	FOREVER_EXPIRATION = time.Duration(-1)

	LOGGING_PREFIX = "[redis_cacher]"
)

type RedisCacher struct {
	redisclient       *redis.Client
	defaultExpiration time.Duration
	Logger            core.ILogger
}

func NewRedisCacher(defaultExpiration time.Duration, logger core.ILogger) *RedisCacher {
	return MakeRedisCacher(GetRedisClient(), defaultExpiration, logger)
}

func MakeRedisCacher(client *redis.Client, defaultExpiration time.Duration, logger core.ILogger) *RedisCacher {
	return &RedisCacher{redisclient: client, defaultExpiration: defaultExpiration, Logger: logger}
}

func exists(conn *redis.Client, key string) bool {
	existed, _ := conn.Do("EXISTS", key).Bool()
	fmt.Println(key, existed)
	//conn.Exists(key).Result()
	return existed
}

func (c *RedisCacher) logErrf(format string, contents ...interface{}) {
	if c.Logger != nil {
		c.Logger.Errorf(fmt.Sprintf("%s %s", LOGGING_PREFIX, format), contents...)
	}
}

func (c *RedisCacher) logDebugf(format string, contents ...interface{}) {
	if c.Logger != nil {
		c.Logger.Debugf(fmt.Sprintf("%s %s", LOGGING_PREFIX, format), contents...)
	}
}

func (c *RedisCacher) getBeanKey(tableName string, id string) string {
	return fmt.Sprintf("xorm:bean:%s:%s", tableName, id)
}

func (c *RedisCacher) getSqlKey(tableName string, sql string) string {
	crc := crc32.ChecksumIEEE([]byte(sql))
	return fmt.Sprintf("xorm:sql:%s:%d", tableName, crc)
}

// Delete all xorm cached objects
func (c *RedisCacher) Flush() error {
	// conn := c.pool.Get()
	// defer conn.Close()
	// _, err := conn.Do("FLUSHALL")
	// return err
	return c.delObject("xorm:*")
}

func (c *RedisCacher) getObject(key string) interface{} {
	conn := c.redisclient
	item, err := conn.Get(key).Bytes()
	if err != nil {
		c.logErrf("redis.Bytes failed: %s", err)
		return nil
	}
	value, err := c.deserialize(item)
	return value
}

func (c *RedisCacher) GetIds(tableName, sql string) interface{} {
	sqlKey := c.getSqlKey(tableName, sql)
	c.logDebugf(" GetIds|tableName:%s|sql:%s|key:%s", tableName, sql, sqlKey)
	return c.getObject(sqlKey)
}

func (c *RedisCacher) GetBean(tableName string, id string) interface{} {
	beanKey := c.getBeanKey(tableName, id)
	c.logDebugf("[xorm/redis_cacher] GetBean|tableName:%s|id:%s|key:%s", tableName, id, beanKey)
	return c.getObject(beanKey)
}

func (c *RedisCacher) putObject(key string, value interface{}) {
	c.invoke(c.redisclient.Do, key, value, c.defaultExpiration)
}

func (c *RedisCacher) PutIds(tableName, sql string, ids interface{}) {
	sqlKey := c.getSqlKey(tableName, sql)
	c.logDebugf("PutIds|tableName:%s|sql:%s|key:%s|obj:%s|type:%v", tableName, sql, sqlKey, ids, reflect.TypeOf(ids))
	c.putObject(sqlKey, ids)
}

func (c *RedisCacher) PutBean(tableName string, id string, obj interface{}) {
	beanKey := c.getBeanKey(tableName, id)
	c.logDebugf("PutBean|tableName:%s|id:%s|key:%s|type:%v", tableName, id, beanKey, reflect.TypeOf(obj))
	c.putObject(beanKey, obj)
}

func (c *RedisCacher) delObject(key string) error {
	c.logDebugf("delObject key:[%s]", key)

	conn := c.redisclient
	if !exists(conn, key) {
		c.logErrf("delObject key:[%s] err: %v", key, core.ErrCacheMiss)
		return core.ErrCacheMiss
	}
	_, err := conn.Del(key).Result()
	//_, err := conn.Do("DEL", key)
	return err
}

func (c *RedisCacher) delObjects(key string) error {

	c.logDebugf("delObjects key:[%s]", key)

	conn := c.redisclient
	//defer conn.Close()

	keys, err := conn.Keys(key).Result()
	c.logDebugf("delObjects keys: %v", keys)

	_, err = conn.Del(keys...).Result()
	return err
}

func (c *RedisCacher) DelIds(tableName, sql string) {
	c.delObject(c.getSqlKey(tableName, sql))
}

func (c *RedisCacher) DelBean(tableName string, id string) {
	c.delObject(c.getBeanKey(tableName, id))
}

func (c *RedisCacher) ClearIds(tableName string) {
	c.delObjects(fmt.Sprintf("xorm:sql:%s:*", tableName))
}

func (c *RedisCacher) ClearBeans(tableName string) {
	c.delObjects(c.getBeanKey(tableName, "*"))
}

func (c *RedisCacher) invoke(f func(...interface{}) *redis.Cmd,
	key string, value interface{}, expires time.Duration) error {

	switch expires {
	case DEFAULT_EXPIRATION:
		expires = c.defaultExpiration
	case FOREVER_EXPIRATION:
		expires = time.Duration(0)
	}

	b, err := c.serialize(value)
	if err != nil {
		return err
	}
	conn := c.redisclient
	//defer conn.Close()
	if expires > 0 {
		_, err := conn.Set(key, b, expires/time.Second).Result()
		//f("SETEX", key, int32(expires/time.Second), b).Result()
		return err
	} else {
		_, err := conn.Set(key, b, 0).Result()
		return err
	}
}

func (c *RedisCacher) serialize(value interface{}) ([]byte, error) {

	err := c.registerGobConcreteType(value)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(value).Kind() == reflect.Struct {
		return nil, fmt.Errorf("serialize func only take pointer of a struct")
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)

	c.logDebugf("serialize type:%v", reflect.TypeOf(value))
	err = encoder.Encode(&value)
	if err != nil {
		c.logErrf("gob encoding '%s' failed: %s|value:%v", value, err, value)
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *RedisCacher) deserialize(byt []byte) (ptr interface{}, err error) {
	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)

	var p interface{}
	err = decoder.Decode(&p)
	if err != nil {
		c.logErrf("decode failed: %v", err)
		return
	}

	v := reflect.ValueOf(p)
	c.logDebugf("deserialize type:%v", v.Type())
	if v.Kind() == reflect.Struct {

		var pp interface{} = &p
		datas := reflect.ValueOf(pp).Elem().InterfaceData()

		sp := reflect.NewAt(v.Type(),
			unsafe.Pointer(datas[1])).Interface()
		ptr = sp
		vv := reflect.ValueOf(ptr)
		c.logDebugf("deserialize convert ptr type:%v | CanAddr:%t", vv.Type(), vv.CanAddr())
	} else {
		ptr = p
	}
	return
}

func (c *RedisCacher) registerGobConcreteType(value interface{}) error {

	t := reflect.TypeOf(value)

	c.logDebugf("registerGobConcreteType:%v", t)

	switch t.Kind() {
	case reflect.Ptr:
		v := reflect.ValueOf(value)
		i := v.Elem().Interface()
		gob.Register(&i)
	case reflect.Struct, reflect.Map, reflect.Slice:
		gob.Register(value)
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		// do nothing since already registered known type
	default:
		return fmt.Errorf("unhandled type: %v", t)
	}
	return nil
}

func (c *RedisCacher) GetPool() (*redis.Client, error) {
	return c.redisclient, nil
}

func (c *RedisCacher) SetPool(client *redis.Client) {
	c.redisclient = client
}
