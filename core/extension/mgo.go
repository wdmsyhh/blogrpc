package extension

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unsafe"

	qmgo_options "github.com/qiniu/qmgo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"blogrpc/core/log"
	"blogrpc/core/util"

	"blogrpc/core/extension/bson"

	"github.com/qiniu/qmgo"
)

const (
	// 1毫秒 = 1000微秒*1000纳秒
	NANOS_PER_MILLI_SECOND             = 1000 * 1000
	DEFAULT_SLOW_QUERY_THRESHOLD_IN_MS = 1000

	DB_CONNECT_STRATEGY_DIRECT = "direct"

	SERVER_UNREACHABLE             = "no reachable servers"
	SERVER_UNREACHABLE_RETRY_TIMES = 2 // this number count in the first try

	// error msg
	MGO_ERROR_DUPLICATE_KEY = "E11000 duplicate key error collection"

	TRANSACTION_ERROR_WRITE_CONFLICT = "WriteConflict"
)

var (
	slowQueryThresholdInMS int64              = DEFAULT_SLOW_QUERY_THRESHOLD_IN_MS
	DBRepository           DatabaseRepository = &mongoDBRepository{
		clientSelector: &proxyClientSelector{},
		clientRemover:  &proxyClientRemover{},
	}
)

var coreCollections = map[string]bool{
	"account":            true,
	"accountDBConfig":    true,
	"validation":         true,
	"user":               true,
	"channel":            true,
	"mobileApp":          true,
	"group":              true,
	"wechatVerification": true,
	"userMenu":           true,
	"userDataPermission": true,
	"app":                true,
	"identityProvider":   true,
	"dmsstore.openUser":  true,
	"helpDesk":           true,
}

func registerDBRepository(strategy string) {
	if strategy == DB_CONNECT_STRATEGY_DIRECT {
		DBRepository = &mongoDBRepository{
			clientSelector: &directClientSelector{},
			clientRemover:  &directClientRemover{},
		}
	}
}

type RunMongodbFunc func(collection *qmgo.Collection) error

func isCoreCollection(collection string) bool {
	_, ok := coreCollections[collection]
	return ok
}

func getClient(ctx context.Context, collection string) *qmgo.QmgoClient {
	var client *qmgo.QmgoClient
	if isCoreCollection(collection) {
		client = MasterDBConnector.GetClient(ctx)
	} else {
		client = TenantDBConnector.GetClient(ctx)
	}

	return client
}

func PingMgo() error {
	return MasterDBConnector.Ping()
}

func init() {
	RegisterExtension(MasterDBConnector)
}

type ClientSelector interface {
	GetClient(ctx context.Context, collectionName string) *qmgo.QmgoClient
}

type proxyClientSelector struct{}

func (*proxyClientSelector) GetClient(ctx context.Context, collectionName string) *qmgo.QmgoClient {
	return getClient(ctx, collectionName)
}

type directClientSelector struct{}

func (*directClientSelector) GetClient(ctx context.Context, collectionName string) *qmgo.QmgoClient {
	return MasterDBConnector.GetClient(ctx)
}

type ClientRemover interface {
	Remove(context.Context, string)
}

type proxyClientRemover struct{}

func (*proxyClientRemover) Remove(ctx context.Context, collectionName string) {
	if isCoreCollection(collectionName) {
		MasterDBConnector.connect(ctx)
	} else {
		accountId := util.MustGetAccountId(ctx)
		TenantDBConnector.Remove(ctx, accountId)
	}
}

type directClientRemover struct{}

func (*directClientRemover) Remove(ctx context.Context, collectionName string) {
	MasterDBConnector.connect(ctx)
}

type DatabaseRepository interface {
	FindOne(ctx context.Context, collectionName string, selector bson.M, result interface{}) error
	FindOneWithSortor(ctx context.Context, collectionName string, selector bson.M, sortor []string, result interface{}) error
	FindAll(ctx context.Context, collectionName string, selector bson.M, sortor []string, limit int, result interface{}) error
	FindAllWithFields(ctx context.Context, collectionName string, selector, fields bson.M, sortor []string, limit int, result interface{}) error
	FindByPK(ctx context.Context, collectionName string, id interface{}, result interface{}) error
	FindByPagination(ctx context.Context, collectionName string, page PagingCondition, result interface{}) (int, error)
	FindByPaginationWithoutCount(ctx context.Context, collectionName string, page PagingCondition, result interface{}) error
	FindByPaginationWithFields(ctx context.Context, collectionName string, page PagingCondition, result interface{}, fields bson.M) (int, error)
	FindAndApply(ctx context.Context, collectionName string, selector bson.M, sort []string, change qmgo.Change, result interface{}) error
	UpdateOne(ctx context.Context, collectionName string, selector bson.M, updator bson.M) error
	UpdateAll(ctx context.Context, collectionName string, selector bson.M, updator bson.M) (int, error)
	UpdateAllWithAggregation(ctx context.Context, collectionName string, selector bson.M, pipeline []bson.M) (int, error)
	Insert(ctx context.Context, collectionName string, docs ...interface{}) ([]bson.ObjectId, error)
	RemoveOne(ctx context.Context, collectionName string, selector bson.M) error
	RemoveAll(ctx context.Context, collectionName string, selector bson.M) (int, error)
	Count(ctx context.Context, collectionName string, selector bson.M) (int, error)
	Aggregate(ctx context.Context, collectionName string, pipeline interface{}, one bool, result interface{}) error
	Distinct(ctx context.Context, collectionName string, selector bson.M, key string, result interface{}) error
	Upsert(ctx context.Context, collectionName string, selector bson.M, updator bson.M) (interface{}, error)
	InsertUnordered(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException)
	BatchUpsert(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException)
	BatchUpdate(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException)
	BatchUpdateUnordered(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException)
	Iterate(ctx context.Context, collectionName string, selector bson.M, sortor []string) (IterWrapper, error)
	IterateWithOption(ctx context.Context, collectionName string, selector bson.M, opt IterateOption) (IterWrapper, error)
	FindAllWithHint(ctx context.Context, collectionName string, selector bson.M, sortor []string, limit int, hint string, result interface{}) error
	Transaction(ctx context.Context, transactionFunc func(sessCtx context.Context) (interface{}, error), opts ...*TransactionOption) (interface{}, error)
	CreateCollection(ctx context.Context, name string) error
}

type PagingCondition struct {
	Selector  bson.M
	PageIndex int
	PageSize  int
	Sortor    []string
}

type mongoDBRepository struct {
	clientSelector ClientSelector
	clientRemover  ClientRemover
}

func (repo *mongoDBRepository) FindOne(ctx context.Context, collectionName string, selector bson.M, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongoFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		q := collection.Find(ctx, selector)
		err := q.One(result)

		repo.logProfile(ctx, "FindOne", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongoFunc)

	return handleDatabaseError(ctx, err, "FindOne", collectionName, selector)
}

func (repo *mongoDBRepository) FindOneWithSortor(ctx context.Context, collectionName string, selector bson.M, sortor []string, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		q := collection.Find(ctx, selector).Sort(sortor...)
		err := q.One(result)
		repo.logProfile(ctx, "FindOneWithSortor", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"sortor":         sortor,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "FindOneWithSortor", collectionName, selector)
}

func (repo *mongoDBRepository) FindAllWithFields(ctx context.Context, collectionName string, selector, fields bson.M, sortor []string, limit int, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		formatSelectorIn(ctx, selector)
		q := collection.Find(ctx, selector).Select(fields)
		if len(sortor) > 0 {
			q = q.Sort(sortor...)
		}
		if limit > 0 {
			q = q.Limit(int64(limit))
		}
		startTime := time.Now()

		err := q.All(result)

		repo.logProfile(ctx, "FindAllWithFields", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"fields":         fields,
			"sortor":         sortor,
			"limit":          limit,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "FindAllWithFields", collectionName, selector, sortor, limit)
}

func (repo *mongoDBRepository) FindAll(ctx context.Context, collectionName string, selector bson.M, sortor []string, limit int, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		formatSelectorIn(ctx, selector)
		q := collection.Find(ctx, selector)
		if len(sortor) > 0 {
			q = q.Sort(sortor...)
		}
		if limit > 0 {
			q = q.Limit(int64(limit))
		}
		startTime := time.Now()

		err := q.All(result)

		repo.logProfile(ctx, "FindAll", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"sortor":         sortor,
			"limit":          limit,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "FindAll", collectionName, selector, sortor, limit)
}

func (repo *mongoDBRepository) FindByPK(ctx context.Context, collectionName string, id interface{}, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		q := collection.Find(ctx, bson.M{"_id": id})
		err := q.One(result)
		repo.logProfile(ctx, "FindByPK", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"id":             id,
		})
		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "FindByPK", collectionName, id)
}

func (repo *mongoDBRepository) FindByPaginationWithFields(ctx context.Context, collectionName string, page PagingCondition, result interface{}, fields bson.M) (int, error) {
	ctx = util.DuplicateContext(ctx)
	var total int
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		formatSelectorIn(ctx, page.Selector)
		skipCount, err := queryPageWhitFileds(ctx, collection, page, result, fields)
		if nil != err {
			return err
		}
		// Log pagination query profile
		repo.logProfile(ctx, "FindByPagination-Pagination", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"page":           page,
		})
		// Total count
		resultv := reflect.ValueOf(result)
		slicev := resultv.Elem()
		// Current queried count
		count := slicev.Len()
		if count > 0 && count < page.PageSize {
			// If current queried count < page size，the pagination reaches the last page and
			// no need to perform the Count() query again
			total = skipCount + count
		} else {
			countStartTime := time.Now()
			qCount := collection.Find(ctx, page.Selector)
			tmpCount, err := qCount.Count()
			total = int(tmpCount)
			if nil != err {
				return err
			}
			// Log total count profile
			repo.logProfile(ctx, "FindByPagination-Count", countStartTime, map[string]interface{}{
				"collectionName": collectionName,
				"page":           page,
			})
		}
		// Log pagination and total count query, q is for pagination query
		repo.logProfile(ctx, "FindByPagination", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"page":           page,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return total, handleDatabaseError(ctx, err, "FindByPaginationWithFields", collectionName, page)
}

func (repo *mongoDBRepository) FindByPagination(ctx context.Context, collectionName string, page PagingCondition, result interface{}) (int, error) {
	ctx = util.DuplicateContext(ctx)
	return repo.FindByPaginationWithFields(ctx, collectionName, page, result, nil)
}

// FindByPaginationWithoutCount is a simplified version of the FindByPagination
// It speeds up the query by not counting the total count
// For the situation: don't need the total count or get it by other maeans always(maybe use cache)
func (repo *mongoDBRepository) FindByPaginationWithoutCount(ctx context.Context, collectionName string, page PagingCondition, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		formatSelectorIn(ctx, page.Selector)
		_, err := queryPage(ctx, collection, page, result)
		if nil != err {
			return err
		}

		repo.logProfile(ctx, "FindByPaginationWithoutCount", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"page":           page,
		})

		return nil
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "FindByPagination", collectionName, page)
}

func queryPage(ctx context.Context, collection *qmgo.Collection, page PagingCondition, result interface{}) (int, error) {
	q := collection.Find(ctx, page.Selector)
	if len(page.Sortor) > 0 {
		q = q.Sort(page.Sortor...)
	}

	skipCount := int((page.PageIndex - 1) * page.PageSize)
	q = q.Skip(int64(skipCount)).Limit(int64(page.PageSize))
	return skipCount, q.All(result)
}

func queryPageWhitFileds(ctx context.Context, collection *qmgo.Collection, page PagingCondition, result interface{}, fields bson.M) (int, error) {
	var q qmgo.QueryI
	if fields != nil {
		q = collection.Find(ctx, page.Selector).Select(fields)
	} else {
		q = collection.Find(ctx, page.Selector)
	}
	if len(page.Sortor) > 0 {
		q = q.Sort(page.Sortor...)
	}

	skipCount := int((page.PageIndex - 1) * page.PageSize)
	q = q.Skip(int64(skipCount)).Limit(int64(page.PageSize))
	return skipCount, q.All(result)
}

func (repo *mongoDBRepository) FindAndApply(ctx context.Context, collectionName string, selector bson.M, sort []string, change qmgo.Change, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		q := collection.Find(ctx, selector)
		q = q.Sort(sort...)
		err = q.Apply(change, result)
		repo.logProfile(ctx, "FindAndApply", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"sort":           sort,
			"change":         change,
		})

		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	// todo
	return handleDatabaseError(ctx, err, "FindAndApply", collectionName, selector)
}

func (repo *mongoDBRepository) UpdateOne(ctx context.Context, collectionName string, selector bson.M, updator bson.M) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		err := collection.UpdateOne(ctx, selector, updator)
		repo.logProfile(ctx, "UpdateOne", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"updator":        updator,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "UpdateOne", collectionName, selector, updator)
}

func (repo *mongoDBRepository) UpdateAll(ctx context.Context, collectionName string, selector bson.M, updator bson.M) (int, error) {
	ctx = util.DuplicateContext(ctx)
	var info *qmgo.UpdateResult
	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		// todo @alomerry wu 数据库升级后支持更新数组中的多个元素 opition，ArrayFilters
		info, err = collection.UpdateAll(ctx, selector, updator)
		repo.logProfile(ctx, "UpdateAll", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"updator":        updator,
		})

		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		return 0, handleDatabaseError(ctx, err, "UpdateAll", collectionName, selector, updator)
	}

	return int(info.ModifiedCount), nil
}

func (repo *mongoDBRepository) UpdateAllWithAggregation(ctx context.Context, collectionName string, selector bson.M, pipeline []bson.M) (int, error) {
	ctx = util.DuplicateContext(ctx)
	var info *qmgo.UpdateResult
	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		info, err = collection.UpdateAll(ctx, selector, pipeline)
		repo.logProfile(ctx, "UpdateAllWithAggregation", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"pipeline":       pipeline,
		})

		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		return 0, handleDatabaseError(ctx, err, "UpdateAllWithAggregation", collectionName, selector, pipeline)
	}

	return int(info.ModifiedCount), nil
}

func (repo *mongoDBRepository) Insert(ctx context.Context, collectionName string, docs ...interface{}) ([]bson.ObjectId, error) {
	ctx = util.DuplicateContext(ctx)
	var ids []bson.ObjectId
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		var err error
		var arr []interface{}
		for _, doc := range docs {
			arr = append(arr, doc)
		}
		if len(docs) > 1 {
			var res *qmgo.InsertManyResult
			res, err = collection.InsertMany(ctx, arr)
			defer func() {
				if r := recover(); r != nil {
					stack := make([]byte, log.MaxStackSize)
					stack = stack[:runtime.Stack(stack, false)]
					log.ErrorTrace(ctx, "GetInsertManyIds", log.Fields{
						"error": fmt.Sprintf("%v", r),
					}, stack)
				}
			}()
			ids = append(ids, bson.GetInsertManyIds(res)...)
		} else {
			var res *qmgo.InsertOneResult
			res, err = collection.InsertOne(ctx, docs[0])
			if res != nil {
				defer func() {
					if r := recover(); r != nil {
						stack := make([]byte, log.MaxStackSize)
						stack = stack[:runtime.Stack(stack, false)]
						log.ErrorTrace(ctx, "GetInsertManyIds", log.Fields{
							"error": fmt.Sprintf("%v", r),
						}, stack)
					}
				}()
				ids = append(ids, bson.GetInsertOneId(res))
			}
		}
		repo.logProfile(ctx, "Insert", startTime, map[string]interface{}{
			"collectionName": collectionName,
		})
		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		return ids, handleDatabaseError(ctx, err, "Insert", collectionName, docs)
	}
	return ids, handleDatabaseError(ctx, err, "Insert", collectionName, docs)
}

func (repo *mongoDBRepository) RemoveOne(ctx context.Context, collectionName string, selector bson.M) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		err := collection.Remove(ctx, selector)
		repo.logProfile(ctx, "RemoveOne", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "RemoveOne", collectionName, selector)
}

func (repo *mongoDBRepository) RemoveAll(ctx context.Context, collectionName string, selector bson.M) (int, error) {
	ctx = util.DuplicateContext(ctx)
	var info *qmgo.DeleteResult
	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		info, err = collection.RemoveAll(ctx, selector)
		repo.logProfile(ctx, "RemoveAll", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
		})
		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		return 0, handleDatabaseError(ctx, err, "RemoveAll", collectionName, selector)
	}

	return int(info.DeletedCount), nil
}

func (repo *mongoDBRepository) Count(ctx context.Context, collectionName string, selector bson.M) (int, error) {
	ctx = util.DuplicateContext(ctx)
	var count int
	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		q := collection.Find(ctx, selector)
		qCount, err := q.Count()
		count = int(qCount)
		repo.logProfile(ctx, "Count", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
		})

		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		return 0, handleDatabaseError(ctx, err, "Count", collectionName, selector)
	}

	return count, nil
}

func (repo *mongoDBRepository) Aggregate(ctx context.Context, collectionName string, pipeline interface{}, one bool, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		startTime := time.Now()
		for i := range pipeline.([]bson.M) {
			formatSelectorIn(ctx, (pipeline.([]bson.M))[i])
		}
		pipe := collection.Aggregate(ctx, pipeline)
		if one {
			err = pipe.One(result)
		} else {
			err = pipe.All(result)
		}
		repo.logProfile(ctx, "Aggregate", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"pipeline":       fmt.Sprint(pipeline), // pipeline是mgo library内部定义的结构体，它的字段都是私有的，不能json化，所以此处用Sprint()
			"one":            one,
			"result":         result,
		})
		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "Aggregate", collectionName, pipeline, one)
}

func (repo *mongoDBRepository) Distinct(ctx context.Context, collectionName string, selector bson.M, key string, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		err := collection.Find(ctx, selector).Distinct(key, result)

		repo.logProfile(ctx, "Distinct", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"key":            key,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "Distinct", collectionName, selector, key)
}

func (repo *mongoDBRepository) Upsert(ctx context.Context, collectionName string, selector bson.M, updator bson.M) (interface{}, error) {
	ctx = util.DuplicateContext(ctx)
	var info *mongo.UpdateResult

	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		startTime := time.Now()
		formatSelectorIn(ctx, selector)
		mongoCollection := *(**mongo.Collection)(unsafe.Pointer(collection))
		info, err = mongoCollection.UpdateOne(ctx, selector, updator, options.Update().SetUpsert(true))
		repo.logProfile(ctx, "Upsert", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"updator":        updator,
		})

		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		return nil, handleDatabaseError(ctx, err, "Upsert", collectionName, selector, updator)
	}

	return info.UpsertedID, nil
}

func (repo *mongoDBRepository) InsertUnordered(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException) {
	ctx = util.DuplicateContext(ctx)
	// 由于 qmgo 在 0.76 版本还未封装 mongo driver 的 bulk 相关操作，并且隐藏了 mongo 的 collection，此处先自行调用 driver，后面根据 qmgo 的支持再做修改
	var result *mongo.BulkWriteResult

	runMongodbFunc := func(collection *qmgo.Collection) (err error) {
		// 获取 mongo driver 的 collection 变量
		mongoCollection := *(**mongo.Collection)(unsafe.Pointer(collection))
		startTime := time.Now()
		flag := false
		bulkOption := &options.BulkWriteOptions{
			Ordered: &flag,
		}
		var models []mongo.WriteModel
		for _, item := range docs {
			model := mongo.NewInsertOneModel()
			model.SetDocument(item)
			models = append(models, model)
		}
		result, err = mongoCollection.BulkWrite(ctx, models, bulkOption)
		repo.logProfile(ctx, "InsertUnordered", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"docs":           docs,
		})
		return
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		bErr := handleBulkError(ctx, err, "InsertUnordered", collectionName, docs)
		return nil, bErr
	}

	return result, nil
}

func (repo *mongoDBRepository) BatchUpsert(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException) {
	ctx = util.DuplicateContext(ctx)
	var result *mongo.BulkWriteResult

	runMongodbFunc := getBatchUpsertFunc(ctx, collectionName, repo, true, result, nil, docs...)

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		bErr := handleBulkError(ctx, err, "BatchUpsert", collectionName, docs)
		return nil, bErr
	}

	return result, nil
}

func (repo *mongoDBRepository) BatchUpdate(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException) {
	ctx = util.DuplicateContext(ctx)
	var result *mongo.BulkWriteResult

	runMongodbFunc := getBatchUpsertFunc(ctx, collectionName, repo, false, result, nil, docs...)

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		bErr := handleBulkError(ctx, err, "BatchUpsert", collectionName, docs)
		return nil, bErr
	}

	return result, nil
}

func (repo *mongoDBRepository) BatchUpdateUnordered(ctx context.Context, collectionName string, docs ...interface{}) (*mongo.BulkWriteResult, *mongo.BulkWriteException) {
	ctx = util.DuplicateContext(ctx)
	var result *mongo.BulkWriteResult

	runMongodbFunc := getBatchUpsertFunc(ctx, collectionName, repo, false, result, &options.BulkWriteOptions{
		Ordered: new(bool),
	}, docs...)

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)
	if err != nil {
		bErr := handleBulkError(ctx, err, "BatchUpsert", collectionName, docs)
		return nil, bErr
	}

	return result, nil
}

type TransactionOption struct {
	CanRetry bool
}

func (repo *mongoDBRepository) Transaction(ctx context.Context, transactionFunc func(sessCtx context.Context) (interface{}, error), opts ...*TransactionOption) (interface{}, error) {
	ctx = util.DuplicateContext(ctx)
	client := repo.clientSelector.GetClient(ctx, "")

	needRetry := true

	session, err := client.Session()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)

	transactionMaxCommitTime := 10 * time.Second

	defer timeCost(ctx, "Transaction time cost")()
	return session.StartTransaction(ctx, wrapperTransactionFunc(session, transactionFunc, needRetry), &qmgo_options.TransactionOptions{
		TransactionOptions: &options.TransactionOptions{
			MaxCommitTime: &transactionMaxCommitTime,
		},
	})
}

func timeCost(ctx context.Context, msg string) func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		log.Warn(ctx, msg, log.Fields{
			"time(ms)": tc.Milliseconds(),
		})
	}
}

func wrapperTransactionFunc(session *qmgo.Session, tFunc func(ctx context.Context) (interface{}, error), needRetry bool) func(ctx context.Context) (interface{}, error) {
	return func(sessCtx context.Context) (result interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				if cmdErr, ok := r.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
					// 暂时只重试 concurrent write conflict 错误
					if needRetry && cmdErr.Name == TRANSACTION_ERROR_WRITE_CONFLICT {
						log.Warn(sessCtx, "TransientTransactionError", log.Fields{"error": cmdErr.Error(), "msg": "retrying transaction"})
						err = qmgo.ErrTransactionRetry
					} else {
						err = errors.New(cmdErr.Error())
						session.AbortTransaction(sessCtx)
					}
				} else {
					stack := make([]byte, log.MaxStackSize)
					stack = stack[:runtime.Stack(stack, false)]
					log.ErrorTrace(sessCtx, "Panic in transaction", log.Fields{
						"error":     fmt.Sprintf("%v", r),
						"accountId": util.GetAccountId(sessCtx),
					}, stack)
					err = errors.New(fmt.Sprintf("%v", r))
					session.AbortTransaction(sessCtx)
				}
			}
			return
		}()
		result, err = tFunc(sessCtx)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (repo *mongoDBRepository) CreateCollection(ctx context.Context, name string) error {
	ctx = util.DuplicateContext(ctx)
	client := repo.clientSelector.GetClient(ctx, "")

	err := client.CreateCollection(ctx, name)

	return handleDatabaseError(ctx, err, "createCollection")
}

func getBatchUpsertFunc(ctx context.Context, collectionName string, repo *mongoDBRepository, needUpsert bool, result *mongo.BulkWriteResult, opt *options.BulkWriteOptions, docs ...interface{}) func(collection *qmgo.Collection) (err error) {
	return func(collection *qmgo.Collection) (err error) {
		mongoCollection := *(**mongo.Collection)(unsafe.Pointer(collection))
		startTime := time.Now()
		var models []mongo.WriteModel
		if len(docs)%2 != 0 {
			panic("Bulk update requires an even number of parameters")
		}
		for i := 0; i < len(docs); i += 2 {
			model := mongo.NewUpdateOneModel()
			model.SetUpsert(needUpsert)
			model.SetFilter(docs[i])
			model.SetUpdate(docs[i+1])
			models = append(models, model)
		}
		result, err = mongoCollection.BulkWrite(ctx, models, opt)
		repo.logProfile(ctx, "BatchUpsert", startTime, map[string]interface{}{
			"collectionName": collectionName,
		})
		return
	}
}

type IterWrapper = qmgo.CursorI

func (repo *mongoDBRepository) Iterate(ctx context.Context, collectionName string, selector bson.M, sortor []string) (IterWrapper, error) {
	opt := IterateOption{
		Sortor: sortor,
	}
	return repo.IterateWithOption(ctx, collectionName, selector, opt)
}

func (repo *mongoDBRepository) IterateWithOption(ctx context.Context, collectionName string, selector bson.M, opt IterateOption) (IterWrapper, error) {
	ctx = util.DuplicateContext(ctx)
	startTime := time.Now()
	formatSelectorIn(ctx, selector)
	client := repo.clientSelector.GetClient(ctx, collectionName)
	collection := client.Database.Collection(collectionName)

	q := collection.Find(ctx, selector)
	q = opt.AppendOptionToQueryI(q)
	it := q.Cursor()
	err := it.Err()

	repo.logProfile(ctx, "Iterate", startTime, map[string]interface{}{
		"collectionName": collectionName,
		"selector":       selector,
		"iterateOption":  opt,
	})

	return it, err
}

func (repo *mongoDBRepository) FindAllWithHint(ctx context.Context, collectionName string, selector bson.M, sortor []string, limit int, hint string, result interface{}) error {
	ctx = util.DuplicateContext(ctx)
	runMongodbFunc := func(collection *qmgo.Collection) error {
		formatSelectorIn(ctx, selector)
		q := collection.Find(ctx, selector)
		if len(sortor) > 0 {
			q = q.Sort(sortor...)
		}
		if limit > 0 {
			q = q.Limit(int64(limit))
		}
		if hint != "" {
			// todo
			//hints := strings.Split(hint, ",")
			//q = q.Hint(hints...)
		}
		startTime := time.Now()

		err := q.All(result)

		repo.logProfile(ctx, "FindAll", startTime, map[string]interface{}{
			"collectionName": collectionName,
			"selector":       selector,
			"sortor":         sortor,
			"limit":          limit,
		})

		return err
	}

	err := repo.runMongodb(ctx, collectionName, runMongodbFunc)

	return handleDatabaseError(ctx, err, "FindAll", collectionName, selector, sortor, limit)

}

func (repo *mongoDBRepository) logProfile(ctx context.Context, method string, startTime time.Time, params map[string]interface{}) {
	duration := time.Now().Sub(startTime)
	if duration.Nanoseconds() >= slowQueryThresholdInMS*NANOS_PER_MILLI_SECOND || strings.HasSuffix(util.ExtractRequestIDFromCtx(ctx), "-00") {
		record := log.NewServiceLog()
		record.ReqId = util.ExtractRequestIDFromCtx(ctx)
		record.Level = "warn"
		record.Category = "db_profiling"
		record.Message = "Slow database operations profiling"
		record.Context = map[string]interface{}{
			"method":    method,
			"startTime": startTime,
			"duration":  duration / NANOS_PER_MILLI_SECOND, // millisecond for readable
			"params":    params,
		}

		bytes, _ := json.Marshal(record)
		log.Stdout.Printf("%s", bytes)
	}
}

func (repo *mongoDBRepository) runMongodb(ctx context.Context, collectionName string, op RunMongodbFunc) (err error) {
	return repo.runMongodbRecursively(ctx, collectionName, op, 1, true)
}

func (repo *mongoDBRepository) runMongodbRecursively(ctx context.Context, collectionName string, op RunMongodbFunc, callTimes int, closeSession bool) (err error) {
	// get session and its collection
	client := repo.clientSelector.GetClient(ctx, collectionName)
	//if closeSession {
	//	defer session.Close()
	//}
	collection := client.Database.Collection(collectionName)

	err = op(collection)

	// 1. if run mongodb successfully, return
	// 2. if the error is not server unreachable, return
	// 3. if reached max retry times, return
	if err == nil || callTimes >= SERVER_UNREACHABLE_RETRY_TIMES || !isServerUnreachable(err) {
		return
	}

	// remove the useless session and try once more
	repo.clientRemover.Remove(ctx, collectionName)
	err = repo.runMongodbRecursively(ctx, collectionName, op, callTimes+1, false)

	return
}

func handleDatabaseError(ctx context.Context, err error, funcName string, args ...interface{}) error {
	if err != nil && err != bson.ErrNotFound {
		logUnexpectedError(ctx, err, funcName, args)
	}

	return err
}

func handleBulkError(ctx context.Context, err error, funcName string, args ...interface{}) *mongo.BulkWriteException {
	var (
		bErr mongo.BulkWriteException
		ok   bool
	)

	if bErr, ok = err.(mongo.BulkWriteException); !ok {
		logUnexpectedError(ctx, err, funcName, args)
	}
	return &bErr
}

func logUnexpectedError(ctx context.Context, err error, funcName string, args ...interface{}) {
	stack := make([]byte, log.MaxStackSize)
	stack = stack[:runtime.Stack(stack, false)]
	log.WarnTrace(ctx, "Error happened during accessing database", log.Fields{
		"function":  funcName,
		"arguments": args,
		"error":     err.Error(),
	}, stack)
	panic(err)
}

func isServerUnreachable(err error) bool {
	return strings.Contains(err.Error(), SERVER_UNREACHABLE)
}

// 升级 golang driver 后，$in 未初始化的数组会报错。panic: (BadValue) $in needs an array。
func formatSelectorIn(ctx context.Context, query bson.M) {
	// 理论上不可能 panic，由于太底层，以防万一
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, log.MaxStackSize)
			stack = stack[:runtime.Stack(stack, false)]
			log.ErrorTrace(ctx, "Uncaught exception", log.Fields{
				"error": fmt.Sprintf("%v", r),
			}, stack)
		}
	}()

	for k, v := range query {
		if k == "$in" {
			if v == nil || reflect.ValueOf(v).Len() == 0 {
				sliceType := reflect.SliceOf(reflect.TypeOf(v).Elem())
				emptySlice := reflect.MakeSlice(sliceType, 0, 1).Interface()
				query[k] = emptySlice
			}
			continue
		}
		if reflect.TypeOf(v) == nil {
			continue
		}
		if reflect.TypeOf(v).Kind() == reflect.Map {
			formatSelectorIn(ctx, v.(bson.M))
			continue
		}
		if reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array {
			for i := 0; i < reflect.ValueOf(v).Len(); i++ {
				vv := reflect.ValueOf(v).Index(i).Interface()
				if reflect.TypeOf(vv).Kind() == reflect.Map {
					formatSelectorIn(ctx, vv.(bson.M))
					continue
				}
			}
			continue
		}
	}
}
