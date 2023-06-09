package extension

import (
	"blogrpc/core/util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"blogrpc/core/extension/bson"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

const (
	CACHE_POOL          = "cache"
	RESQUE_POOL         = "resque"
	GLOBAL_QUEUE        = "global"
	RESPONSE_CACHE_POOL = "response_cache"
)

var (
	// RedisClient is the proxy instance for actaul request instance
	RedisClient ReidsClientProxy
	redisClient *RedisManager
)

func init() {
	redisClient := NewRedisClient()
	RedisClient = redisClient
	RegisterExtension(redisClient)
}

type RedisManager struct {
	name    string
	conf    map[string]interface{}
	poolMap map[string]*redis.Pool
}

type ReidsClientProxy interface {
	Keys(key string) ([]string, error)
	Exists(key string) (bool, error)
	Del(key string) (bool, error)
	Expire(key string, time int64) (bool, error)
	Ttl(key string) (int64, error)
	Set(key string, val string) error
	SetNX(key string, value string, ttl int) (bool, error)
	SetEx(key string, time int64, val string) (bool, error)
	SAdd(key string, val string) (bool, error)
	SRem(key string, val string) (bool, error)
	SPop(key string) (string, error)
	SCard(key string) (int, error)
	SMembers(key string) ([]interface{}, error)
	Get(key string) (string, error)
	Incr(key string) (int, error)
	IncrBy(key string, increment int) (int, error)
	Hget(hashName string, key string) (string, error)
	Hset(hashName string, key string, val string) error
	Hdel(hashName string) (bool, error)
	Hincrby(hashName string, key string, val int) (int, error)
	Hmset(hashName string, data map[string]string) error
	Hgetall(hashName string) (map[string]string, error)
	ZRangeByScore(key string, args ...interface{}) (interface{}, error)
	Enqueue(className string, args interface{}) (string, error)
	EnqueueWithName(className string, args interface{}, name string) (string, error)
	AddJobToResque(module, jobName string, jobArg interface{}) (string, error)
	AddDelayedJobToResque(module, jobName string, jobArg interface{}, executeTime int64) error
	Blpop(listKey string, timeout int64) ([]string, error)
	Lpush(listKey string, val string) error
	GetSet(key string, val string) (string, error)
	BLPeek(key string, timeout int64) (string, error)
	Ping() (string, error)
	IScan(idx int, cursor, key string, limit int) (string, []string, error)
	Scan(cursor, key string, limit int) (string, []string, error)
	GetResponseCache(key string) (string, error)
	SetResponseCache(key string, time int64, val string) (bool, error)
	IScanResponseCache(idx int, cursor, key string, limit int) (string, []string, error)
	ScanResponseCache(cursor, key string, limit int) (string, []string, error)
	DelResponseCache(key string) (bool, error)
}

type job struct {
	Id    string      `json:"id"`
	Class string      `json:"class"`
	Args  interface{} `json:"args"`
}

type delayedJob struct {
	Class string      `json:"class"`
	Args  interface{} `json:"args"`
	Queue string      `json:"queue"`
}

type jobStatus struct {
	Status  int64 `json:"status"`
	Updated int64 `json:"updated"`
	Started int64 `json:"started"`
}

func NewRedisClient() *RedisManager {
	if redisClient == nil {
		return &RedisManager{
			name: "redis",
			conf: make(map[string]interface{}),
			poolMap: map[string]*redis.Pool{
				CACHE_POOL:  nil,
				RESQUE_POOL: nil,
			},
		}
	}
	return redisClient
}

func (client *RedisManager) Name() string {
	return client.name
}

func (client *RedisManager) InitWithConf(conf map[string]interface{}, debug bool) error {
	var err error
	client.conf = conf
	client.conf["debug"] = debug

	host := util.GetCacheHost()
	port := util.GetCachePort()
	pass := util.GetCachePassword()
	db := cast.ToString(conf["db"])
	responseDb := cast.ToString(conf["response-cache-db"])

	rhost := fmt.Sprintf("%s:%s", host, port)
	log.Infof("Init Redis with host: %s", rhost)

	client.poolMap[CACHE_POOL], err = InitRedisPool(rhost, pass, db)
	if err != nil {
		return err
	}

	if responseDb != "" {
		client.poolMap[RESPONSE_CACHE_POOL], err = InitRedisPool(rhost, pass, responseDb)
		if err != nil {
			return err
		}
	}

	host = util.GetResqueHost()
	port = util.GetResquePort()
	pass = util.GetResquePassword()
	db = cast.ToString(conf["resque-db"])
	if db == "" {
		return nil
	}

	rhost = fmt.Sprintf("%s:%s", host, port)
	log.Infof("Init Resque with host: %s", rhost)

	client.poolMap[RESQUE_POOL], err = InitRedisPool(rhost, pass, db)
	if err != nil {
		return err
	}

	log.Info("Successfully loaded Redis extension")
	return nil
}

func (client *RedisManager) GetPool(category string) *redis.Pool {
	if pool, ok := client.poolMap[category]; ok {
		return pool
	}
	return nil
}

func (client *RedisManager) Close() {
	for _, pool := range client.poolMap {
		pool.Close()
	}
}

func (client *RedisManager) Conn(category string) redis.Conn {
	if debug, ok := client.conf["debug"].(bool); ok && debug {
		conn := client.poolMap[category].Get()

		return redis.NewLoggingConn(conn, getLogger("[redis]  "), "REDIS-DEBUG")
	}

	return client.poolMap[category].Get()
}

func (client *RedisManager) Keys(key string) ([]string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Strings(conn.Do("KEYS", key))
}

func (client *RedisManager) Exists(key string) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

func (client *RedisManager) Del(key string) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
}

func (client *RedisManager) Expire(key string, time int64) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Bool(conn.Do("EXPIRE", key, time))
}

func (client *RedisManager) Ttl(key string) (int64, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Int64(conn.Do("TTL", key))
}

func (client *RedisManager) Set(key string, val string) error {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	_, err := conn.Do("SET", key, val)

	return err
}

func (client *RedisManager) SetNX(key string, val string, ttl int) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	ok, err := redis.String(conn.Do("SET", key, val, "EX", ttl, "NX"))
	if ok == "" {
		return false, err
	}

	return true, err
}

func (client *RedisManager) SetEx(key string, expire int64, val string) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	ok, err := redis.String(conn.Do("SETEX", key, expire, val))
	if ok == "" {
		return false, err
	}

	return true, err
}

func (client *RedisManager) SAdd(key string, val string) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Bool(conn.Do("SADD", key, val))
}

func (client *RedisManager) SPop(key string) (string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.String(conn.Do("SPOP", key))
}

func (client *RedisManager) SCard(key string) (int, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Int(conn.Do("SCARD", key))
}

func (client *RedisManager) SMembers(key string) ([]interface{}, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Values(conn.Do("SMEMBERS", key))
}

func (client *RedisManager) SRem(key string, val string) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Bool(conn.Do("SREM", key, val))
}

func (client *RedisManager) Get(key string) (string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (client *RedisManager) Incr(key string) (int, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Int(conn.Do("INCR", key))
}

func (client *RedisManager) IncrBy(key string, increment int) (int, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Int(conn.Do("INCRBY", key, increment))
}

func (client *RedisManager) Hget(hashName string, key string) (string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.String(conn.Do("HGET", hashName, key))
}

func (client *RedisManager) Hset(hashName string, key string, val string) error {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	_, err := conn.Do("HSET", hashName, key, val)
	return err
}

func (client *RedisManager) Hdel(hashName string) (bool, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Bool(conn.Do("HDEL", hashName))
}

func (client *RedisManager) Hincrby(hashName string, key string, val int) (int, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Int(conn.Do("HINCRBY", hashName, key, val))
}

func (client *RedisManager) Hmset(hashName string, data map[string]string) error {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()

	setContent := []interface{}{hashName}
	for k, v := range data {
		setContent = append(setContent, k, v)
	}
	_, err := conn.Do("HMSET", setContent...)
	return err
}

func (client *RedisManager) Hgetall(hashName string) (map[string]string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()

	reply, err := redis.Strings(conn.Do("HGETALL", hashName))
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for index := 1; index <= len(reply); index += 2 {
		key := reply[index-1]
		value := reply[index]
		result[key] = value
	}
	return result, nil
}

func (client *RedisManager) Blpop(listKey string, timeout int64) ([]string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.Strings(conn.Do("BLPOP", listKey, timeout))
}

func (client *RedisManager) Lpush(listKey string, val string) error {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()

	_, err := conn.Do("LPUSH", listKey, val)
	return err
}

func (client *RedisManager) BLPeek(key string, timeout int64) (string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()

	return redis.String(conn.Do("BRPOPLPUSH", key, key, timeout))
}

func (client *RedisManager) GetSet(key string, val string) (string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.String(conn.Do("GETSET", key, val))
}

func (client *RedisManager) ZRangeByScore(key string, args ...interface{}) (interface{}, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	zArgs := []interface{}{}
	zArgs = append(zArgs, key)
	zArgs = append(zArgs, args...)
	return conn.Do("ZRANGEBYSCORE", zArgs...)
}

func (client *RedisManager) Ping() (string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	return redis.String(conn.Do("PING"))
}

func (client *RedisManager) IScan(idx int, cursor string, key string, limit int) (string, []string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	if limit > 1000 {
		limit = 1000
	}
	if cursor == "" {
		cursor = "0"
	}

	result, err := conn.Do("ISCAN", idx, cursor, "MATCH", key, "COUNT", limit)
	if err != nil {
		return "", nil, err
	}
	scanResults := result.([]interface{})

	keys := []string{}
	for _, k := range scanResults[1].([]interface{}) {
		keys = append(keys, string(k.([]byte)))
	}

	return cast.ToString(scanResults[0]), keys, nil
}

func (client *RedisManager) Scan(cursor string, key string, limit int) (string, []string, error) {
	conn := client.Conn(CACHE_POOL)
	defer conn.Close()
	if limit > 1000 {
		limit = 1000
	}
	if cursor == "" {
		cursor = "0"
	}

	result, err := conn.Do("SCAN", cursor, "MATCH", key, "COUNT", limit)
	if err != nil {
		return "", nil, err
	}
	scanResults := result.([]interface{})

	keys := []string{}
	for _, k := range scanResults[1].([]interface{}) {
		keys = append(keys, string(k.([]byte)))
	}

	return cast.ToString(scanResults[0]), keys, nil
}

func getQueue(args interface{}) string {
	argsMaps := cast.ToSlice(args)
	if argsMaps == nil {
		return GLOBAL_QUEUE
	}

	for _, argsMap := range argsMaps {
		argsM, ok := argsMap.(map[string]interface{})
		if ok {
			if queue, ok := argsM["queue"]; ok {
				return cast.ToString(queue)
			}
		}
	}

	return GLOBAL_QUEUE
}

func (client *RedisManager) Enqueue(className string, args interface{}) (string, error) {
	return client.EnqueueWithName(className, args, GLOBAL_QUEUE)
}

func (client *RedisManager) EnqueueWithName(className string, args interface{}, queueName string) (string, error) {
	if queueName == "" {
		queueName = GLOBAL_QUEUE
	}

	conn := client.Conn(RESQUE_POOL)
	defer conn.Close()

	// start transaction
	conn.Send("MULTI")

	conn.Send("SADD", "wmresque:queues", queueName)

	queue := fmt.Sprintf("wmresque:queue:%s", queueName)
	jobId := bson.NewObjectId().Hex()
	queueInfo := &job{Class: className, Args: args, Id: jobId}
	jobInfo, _ := json.Marshal(queueInfo)
	conn.Send("RPUSH", queue, jobInfo)

	jobStatus := &jobStatus{
		Status:  1,
		Updated: time.Now().Unix(),
		Started: time.Now().Unix(),
	}
	status, _ := json.Marshal(jobStatus)
	conn.Send("SET", "wmresque:job:"+jobId+":status", status)

	_, err := conn.Do("EXEC")
	if err != nil {
		return "", err
	}

	return jobId, nil
}

func (client *RedisManager) AddJobToResque(module, jobName string, jobArg interface{}) (string, error) {
	jobClass := genJobClassName(module, jobName)
	queueName := getQueue(jobArg)
	jobId, err := client.EnqueueWithName(jobClass, jobArg, queueName)

	return jobId, err
}

func (client *RedisManager) AddDelayedJobToResque(module, jobName string, jobArg interface{}, executeTime int64) error {
	className := genJobClassName(module, jobName)

	conn := client.Conn(RESQUE_POOL)
	defer conn.Close()

	conn.Send("MULTI")

	timestampStr := strconv.FormatInt(executeTime, 10)
	queue := fmt.Sprint("wmresque:delayed:", timestampStr)
	queueInfo := &delayedJob{Class: className, Args: jobArg, Queue: GLOBAL_QUEUE}
	jobInfo, _ := json.Marshal(queueInfo)
	conn.Send("RPUSH", queue, jobInfo)
	conn.Send("ZADD", "wmresque:delayed_queue_schedule", executeTime, executeTime)

	_, err := conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

func (client *RedisManager) GetResponseCache(key string) (string, error) {
	conn := client.Conn(RESPONSE_CACHE_POOL)
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (client *RedisManager) SetResponseCache(key string, expire int64, val string) (bool, error) {
	conn := client.Conn(RESPONSE_CACHE_POOL)
	defer conn.Close()
	ok, err := redis.String(conn.Do("SETEX", key, expire, val))
	if ok == "" {
		return false, err
	}

	return true, err
}

func (client *RedisManager) IScanResponseCache(idx int, cursor string, key string, limit int) (string, []string, error) {
	conn := client.Conn(RESPONSE_CACHE_POOL)
	defer conn.Close()
	if limit > 1000 {
		limit = 1000
	}
	if cursor == "" {
		cursor = "0"
	}

	result, err := conn.Do("ISCAN", idx, cursor, "MATCH", key, "COUNT", limit)
	if err != nil {
		return "", nil, err
	}
	scanResults := result.([]interface{})

	keys := []string{}
	for _, k := range scanResults[1].([]interface{}) {
		keys = append(keys, string(k.([]byte)))
	}

	return cast.ToString(scanResults[0]), keys, nil
}

func (client *RedisManager) ScanResponseCache(cursor string, key string, limit int) (string, []string, error) {
	conn := client.Conn(RESPONSE_CACHE_POOL)
	defer conn.Close()
	if limit > 1000 {
		limit = 1000
	}
	if cursor == "" {
		cursor = "0"
	}

	result, err := conn.Do("SCAN", cursor, "MATCH", key, "COUNT", limit)
	if err != nil {
		return "", nil, err
	}
	scanResults := result.([]interface{})

	keys := []string{}
	for _, k := range scanResults[1].([]interface{}) {
		keys = append(keys, string(k.([]byte)))
	}

	return cast.ToString(scanResults[0]), keys, nil
}

func (client *RedisManager) DelResponseCache(key string) (bool, error) {
	conn := client.Conn(RESPONSE_CACHE_POOL)
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
}

func genJobClassName(module, jobName string) string {
	return fmt.Sprintf("backend\\modules\\%s\\job\\%s", module, jobName)
}

func InitRedisPool(host, pass, db string) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", host, 2*time.Second, 0, 0)
			if err != nil {
				log.Errorf("Failed to connect to Redis with host: %s, error: %v", host, err)
				return nil, err
			}

			if "" != pass {
				if _, err := c.Do("AUTH", pass); err != nil {
					log.Errorf("Failed to auth Redis user with password %s, error: %v", pass, err)
					c.Close()
					return nil, err
				}
			}
			if _, err = c.Do("SELECT", db); err != nil {
				log.Errorf("Failed to select Redis db %s, error: %v", db, err)
				c.Close()
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	//mongo Redis connection
	conn := pool.Get()
	defer conn.Close() //close it to pool after tested

	_, err := conn.Do("PING")
	if nil != err {
		log.Errorf("Failed to send Ping to Redis server, with error: %v", err)
		return nil, err
	}
	return pool, nil
}
