package extension

import (
	rpc_errors "blogrpc/core/errors"
	"blogrpc/core/extension/bson"
	"blogrpc/core/log"
	"blogrpc/core/util"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/qiniu/qmgo"
	qmgo_options "github.com/qiniu/qmgo/options"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/tag"
	"golang.org/x/net/context"
)

const (
	REFRESH_DSN_INTERVAL             = 10 // seconds，刷新每个租户数据库连接 URI 的时间间隔
	UPDATE_TENANT_BAD_HOSTS_INTERVAL = 10 // seconds，更新 badHosts 的时间间隔（释放早期失败的记录）
)

var (
	TenantDBConnector *tenantDBConnector = nil
	MasterDBConnector *masterDBConnector = nil
	InDebug                              = false
)

func init() {
	MasterDBConnector = &masterDBConnector{
		conf:   make(map[string]interface{}),
		client: nil,
	}

	TenantDBConnector = &tenantDBConnector{
		hosts: sync.Map{},
		stats: make(map[string]*clientStats),
		conf:  make(map[string]interface{}),
	}
}

type clientStats struct {
	client       *qmgo.QmgoClient
	latestUsedAt time.Time
}

type DBConnector interface {
	GetClient(ctx context.Context) *qmgo.QmgoClient
}

type masterDBConnector struct {
	conf   map[string]interface{}
	client *qmgo.QmgoClient
}

func (*masterDBConnector) Name() string {
	return "mgo"
}

func (m *masterDBConnector) InitWithConf(conf map[string]interface{}, debug bool) error {
	slowQuery := cast.ToInt64(conf["slow-query-threshold-in-ms"])
	if slowQuery > 0 {
		slowQueryThresholdInMS = slowQuery
	}

	InDebug = debug
	m.conf = conf
	err := m.connect(context.Background()) // todo
	if nil != err {
		return err
	}

	// todo
	//mgo.SetDebug(debug)
	registerDBRepository(cast.ToString(conf["strategy"]))

	TenantDBConnector.triggerRefreshHosts()
	TenantDBConnector.triggerUpdateBadHosts()
	TenantDBConnector.conf = conf
	return nil
}

func (m *masterDBConnector) Close() {
	// todo
	m.client.Close(context.Background())
}

func (m *masterDBConnector) Ping() error {
	// todo
	return m.client.Ping(1000)
}

func (m *masterDBConnector) Host() string {
	host := buildMongoHosts(m.conf)
	m.conf["replset"] = util.GetMongoMasterReplset()
	return getMgoConnectString(host, m.conf)
}

func (m *masterDBConnector) connect(ctx context.Context) error {
	client, err := dialClient(ctx, m.Host())
	if err != nil {
		return err
	}

	m.client = client

	return nil
}

func (m *masterDBConnector) GetClient(ctx context.Context) *qmgo.QmgoClient {
	return m.client
}

type tenantDBConnector struct {
	stats    map[string]*clientStats
	hosts    sync.Map
	lock     sync.Mutex
	conf     map[string]interface{}
	badHosts sync.Map
}

// func (tenant *tenantDBConnector) GetSession(ctx context.Context) *qmgo.QmgoClient {
func (tenant *tenantDBConnector) GetClient(ctx context.Context) *qmgo.QmgoClient {
	accountId := util.MustGetAccountId(ctx)
	return tenant.getClient(ctx, accountId)
}

func (tenant *tenantDBConnector) triggerRefreshHosts() {
	ctx := context.Background()

	go func() {
		// due to import cycle, we copy/paste this
		// recover code from goroutine_wrapper.go
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, log.MaxStackSize)
				stack = stack[:runtime.Stack(stack, false)]
				// if panic, set custom error to 'err', in order that client and sense it.
				err := rpc_errors.ConvertRecoveryError(r)
				log.ErrorTrace(ctx, "Panic in goroutine", log.Fields{
					"error": err.Error(),
				}, stack)
			}
		}()

		for {
			time.Sleep(REFRESH_DSN_INTERVAL * time.Second)
			TenantDBConnector.refreshAllHosts(ctx)
		}
	}()
}

func (tenant *tenantDBConnector) triggerUpdateBadHosts() {
	ctx := context.Background()

	go func() {
		// due to import cycle, we copy/paste this
		// recover code from goroutine_wrapper.go
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, log.MaxStackSize)
				stack = stack[:runtime.Stack(stack, false)]
				// if panic, set custom error to 'err', in order that client and sense it.
				err := rpc_errors.ConvertRecoveryError(r)
				log.ErrorTrace(ctx, "Panic in goroutine", log.Fields{
					"error": err.Error(),
				}, stack)
			}
		}()

		for {
			time.Sleep(UPDATE_TENANT_BAD_HOSTS_INTERVAL * time.Second)
			TenantDBConnector.updateBadHosts()
		}
	}()
}

func (tenant *tenantDBConnector) updateBadHosts() {
	ignoredHosts := []string{}
	tenant.badHosts.Range(func(host, lastFailedTime interface{}) bool {
		ignoredHosts = append(ignoredHosts, host.(string))
		return true
	})
	for _, ignoredHost := range ignoredHosts {
		tenant.badHosts.Delete(ignoredHost)
	}
}

func (tenant *tenantDBConnector) refreshAllHosts(ctx context.Context) {
	var accountIds []bson.ObjectId
	tenant.lock.Lock()
	tenant.hosts.Range(func(k, v interface{}) bool {
		accountIds = append(accountIds, bson.ObjectIdHex(k.(string)))
		return true
	})
	tenant.lock.Unlock()

	// no cached host, then need not update
	if len(accountIds) == 0 {
		return
	}

	accountDBConfigs := CAccountDBConfig.GetByAccountIds(ctx, accountIds)
	for _, config := range accountDBConfigs {
		strAccountId := config.AccountId.Hex()
		host := config.Host()

		tenant.lock.Lock()
		cachedHost, ok := tenant.hosts.Load(strAccountId)

		// if host no changes, do nothing
		if ok && host == cachedHost.(string) {
			tenant.lock.Unlock()
			continue
		}

		if _, ok := tenant.badHosts.Load(host); ok {
			tenant.lock.Unlock()
			continue
		}

		// host changed, so we cache the new host
		tenant.hosts.Store(strAccountId, host)
		log.Warn(ctx, "The dsn changed", log.Fields{
			"prev":      cachedHost.(string),
			"current":   host,
			"accountId": strAccountId,
		})
		tenant.lock.Unlock()
	}
}

// getHost will return host for given accountId. the accountId must exists
// in ctx. This function have no thread lock, so make sure you locked before calling it
func (tenant *tenantDBConnector) GetHost(ctx context.Context, accountId string) string {
	if host, exists := tenant.hosts.Load(accountId); exists {
		return host.(string)
	}

	// the logic here is for the first concurrency request to get host in db
	accountDBConfig := CAccountDBConfig.Get(ctx, accountId)
	if !accountDBConfig.Id.Valid() {
		msg := fmt.Sprintf("Account: %s no configuration of hosts", accountId)
		log.Warn(ctx, msg, nil)
		panic(msg)
	}

	host := accountDBConfig.Host()
	tenant.hosts.Store(accountId, host)

	return host
}

// getSession will return a session for given accountId. the accountId must exists
// in ctx
func (tenant *tenantDBConnector) getClient(ctx context.Context, accountId string) *qmgo.QmgoClient {
	// add lock to avoid get session multiple times when encountering
	// high concurrency requests
	tenant.lock.Lock()
	defer tenant.lock.Unlock()

	host := tenant.GetHost(ctx, accountId)

	if client, exists := tenant.hitClient(host); exists {
		return client
	}

	if t, ok := tenant.badHosts.Load(host); ok {
		log.Error(ctx, "Bad host", log.Fields{
			"accountId":     accountId,
			"host":          host,
			"lastErrorTime": t,
		})
		panic(fmt.Errorf("Cached host of account %s is bad", accountId))
	}

	// the logic here is for the first concurrency request to dial db and create session
	//	session, err := dialSession(host, getSessionMode(tenant.conf))
	client, err := dialClient(ctx, host)
	if err != nil {
		tenant.badHosts.Store(host, time.Now())
		log.Error(ctx, "failed to dial session", log.Fields{
			"accountId": accountId,
			"error":     err,
		})
		panic(err)
	}
	tenant.addClient(host, client)
	return client
}

func (tenant *tenantDBConnector) addClient(host string, client *qmgo.QmgoClient) {
	tenant.stats[host] = &clientStats{
		client:       client,
		latestUsedAt: time.Now(),
	}
}

func (tenant *tenantDBConnector) hitClient(host string) (*qmgo.QmgoClient, bool) {
	if stats, ok := tenant.stats[host]; ok {
		stats.latestUsedAt = time.Now()
		return stats.client, true
	}
	return nil, false
}

// Remove will delete a session by given accountId. Since the session
// will be deleted, all the accountId that using this session will be
// deleted too
func (tenant *tenantDBConnector) Remove(ctx context.Context, accountId string) {
	tenant.lock.Lock()
	defer tenant.lock.Unlock()

	host, _ := tenant.hosts.Load(accountId)

	client, exists := tenant.stats[host.(string)]
	if exists {
		client.client.Close(ctx)
		delete(tenant.stats, host.(string))
	}
}

func buildMongoHosts(conf map[string]interface{}) string {
	return util.GetMongoMasterDsn()
}

func getMgoConnectString(hosts string, options map[string]interface{}) string {
	replSet := cast.ToString(options["replset"])
	if replSet != "" && replSet != "none" {
		hosts += "?replicaSet=" + replSet
	}
	return hosts
}

func dialClient(ctx context.Context, hosts string) (*qmgo.QmgoClient, error) {
	dialConfig, diaOptions, err := ParseURL(hosts)
	if err != nil {
		return nil, err
	}
	diaOptions.SetSocketTimeout(1 * time.Minute)
	diaOptions.SetMaxConnIdleTime(1 * time.Minute)
	diaOptions.SetConnectTimeout(5 * time.Second)
	diaOptions.SetServerSelectionTimeout(2 * time.Second)
	diaOptions.SetMaxPoolSize(4096) // copy from mgo
	diaOptions.ReadPreference, err = readpref.New(readpref.PrimaryMode)
	if err != nil {
		return nil, err
	}

	diaOptions.Registry = bson.DefaultRegistry

	qmoClient, err := qmgo.Open(ctx, dialConfig, qmgo_options.ClientOptions{
		ClientOptions: diaOptions,
	})
	if err != nil {
		return nil, err
	}

	return qmoClient, err
}

func ParseURL(url string) (*qmgo.Config, *options.ClientOptions, error) {
	info, err := extractURL(url)
	if err != nil {
		return nil, nil, err
	}

	config := &qmgo.Config{
		Uri:      info.uri,
		Database: info.db,
	}

	poolLimit := 0
	opts := options.Client()
	opts.SetAppName(util.GetMongoAppName())
	opts.SetDirect(false)
	var readPreferenceTags []tag.Set
	readPreferenceMode := readpref.PrimaryMode
	var preference *readpref.ReadPref
	for _, opt := range info.options {
		switch opt.key {
		case "replicaSet":
			opts.SetReplicaSet(opt.value)
		case "maxPoolSize":
			poolLimit, err = strconv.Atoi(opt.value)
			if err != nil {
				return nil, nil, errors.New("bad value for maxPoolSize: " + opt.value)
			}
			opts.SetMaxPoolSize(uint64(poolLimit))
		case "appName":
			if len(opt.value) > 128 {
				return nil, nil, errors.New("appName too long, must be < 128 bytes: " + opt.value)
			}
			opts.SetAppName(opt.value)
		case "readPreferenceTags":
			preferenceTags := strings.Split(opt.value, ",")
			var set tag.Set
			for _, preferenceTag := range preferenceTags {
				kvp := strings.Split(preferenceTag, ":")
				if len(kvp) != 2 {
					return nil, nil, errors.New("bad value for readPreferenceTags: " + opt.value)
				}
				set = append(set, tag.Tag{Name: kvp[0], Value: kvp[1]})
			}
			readPreferenceTags = append(readPreferenceTags, set)
		case "minPoolSize":
			minPoolSize, err := strconv.Atoi(opt.value)
			if err != nil {
				return nil, nil, errors.New("bad value for minPoolSize: " + opt.value)
			}
			if minPoolSize < 0 {
				return nil, nil, errors.New("bad value (negative) for minPoolSize: " + opt.value)
			}
			opts.SetMinPoolSize(uint64(minPoolSize))
		case "maxIdleTimeMS":
			maxIdleTimeMS, err := strconv.Atoi(opt.value)
			if err != nil {
				return nil, nil, errors.New("bad value for maxIdleTimeMS: " + opt.value)
			}
			if maxIdleTimeMS < 0 {
				return nil, nil, errors.New("bad value (negative) for maxIdleTimeMS: " + opt.value)
			}
			opts.SetMaxConnIdleTime(time.Duration(maxIdleTimeMS))
		case "connect":
			if opt.value == "direct" {
				opts.SetDirect(true)
				break
			}
			if opt.value == "replicaSet" {
				break
			}
			fallthrough
		case "authSource":
			config.Uri = fmt.Sprintf("%s?authSource=%s", config.Uri, opt.value)
		default:
			return nil, nil, errors.New("unsupported connection URL option: " + opt.key + "=" + opt.value)
		}
	}
	preference, err = readpref.New(readPreferenceMode, func(readPref *readpref.ReadPref) error {
		// readPreferenceTags
		return nil
	})
	opts.SetReadPreference(preference)

	return config, opts, nil
}

func extractURL(s string) (*urlInfo, error) {
	var opts []urlInfoOption
	var db string
	if c := strings.Index(s, "?"); c != -1 {
		for _, pair := range strings.FieldsFunc(s[c+1:], isOptSep) {
			l := strings.SplitN(pair, "=", 2)
			if len(l) != 2 || l[0] == "" || l[1] == "" {
				return nil, errors.New("connection option must be key=value: " + pair)
			}
			opts = append(opts, urlInfoOption{key: l[0], value: l[1]})
		}
		s = s[:c]
	}
	if c := util.LastIndexOf(s, "/"); c != -1 {
		db = s[c+1:]
	}
	info := urlInfo{
		uri:     s,
		db:      db,
		options: opts,
	}
	return &info, nil
}

type urlInfo struct {
	uri     string
	db      string
	options []urlInfoOption
}

type urlInfoOption struct {
	key   string
	value string
}

func isOptSep(c rune) bool {
	return c == ';' || c == '&'
}

//func (*tenantDBConnector) GetClient(ctx context.Context, collection string) *qmgo.QmgoClient {
//	// 如果是在容器内链接另一个 mongo 容器，需要使用 mongo 容器的内部端口
//	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://mongo:27017", Database: "class", Coll: collection, Auth: &qmgo.Credential{
//		Username: "root",
//		Password: "root",
//	},
//	})
//	if err != nil {
//		panic(err)
//	}
//	return cli
//}
