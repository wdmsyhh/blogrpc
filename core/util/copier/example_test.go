package copier

import (
	"blogrpc/core/extension/bson"
	"blogrpc/core/util"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestCopyModelToModelWithOption(t *testing.T) {
	type OriginLocation struct {
		City      string
		Latitude  float64
		Longitude float64
	}
	type originModel struct {
		Id        bson.ObjectId
		Name      string
		Age       int
		Money     float64
		BirthDay  time.Time
		IsDeleted bool
		Location  OriginLocation
		Detail    *OriginLocation
	}

	var target originModel
	origin := originModel{
		Id:        bson.NewObjectId(),
		Name:      "MockModel",
		Money:     7.7,
		BirthDay:  time.Now(),
		IsDeleted: false,
		Location: OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
		Detail: &OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
	}

	target = originModel{
		Age:   48,
		Money: 100,
	}
	assert.Nil(t, Instance(NewOption().SetIgnoreEmpty(true).SetOverwrite(false)).From(origin).CopyTo(&target))
	origin.Age = 48
	origin.Money = 100
	assert.Equal(t, origin, target)
}

func TestCopyModelToProtoWithArrayField(t *testing.T) {
	type OriginLocation struct {
		City      string
		Latitude  float64
		Longitude float64
	}
	type originModel struct {
		Id        bson.ObjectId
		Name      string
		Age       int
		Money     float64
		BirthDay  time.Time
		IsDeleted bool
		Locations []OriginLocation
	}
	type TargetLocation struct {
		City      string
		Latitude  float32
		Longitude float32
	}
	type targetModel struct {
		Id        string
		Name      string
		Age       int
		Money     float64
		BirthDay  string
		Locations []*TargetLocation
	}

	var target targetModel
	originId := bson.NewObjectId()
	originTime := time.Now()
	origin := originModel{
		Id:        originId,
		Name:      "MockModel",
		Age:       21,
		Money:     7.7,
		BirthDay:  originTime,
		IsDeleted: false,
		Locations: []OriginLocation{
			{
				City:      "ShangHai",
				Latitude:  234.123412,
				Longitude: 3423.43265,
			},
			{
				City:      "Beijing",
				Latitude:  232.123443,
				Longitude: 12323.141,
			},
		},
	}

	assert.Nil(t, Instance(nil).From(origin).CopyTo(&target))
	assert.Equal(t, target, targetModel{
		Id:       originId.Hex(),
		Name:     "MockModel",
		Age:      21,
		Money:    7.7,
		BirthDay: originTime.Format(util.RFC3339),
		Locations: []*TargetLocation{
			{
				City:      "ShangHai",
				Latitude:  234.123412,
				Longitude: 3423.43265,
			},
			{
				City:      "Beijing",
				Latitude:  232.123443,
				Longitude: 12323.141,
			},
		},
	})
}

func TestCopyModelToProto(t *testing.T) {
	type OriginLocation struct {
		City      string
		Latitude  float64
		Longitude float64
	}
	type originModel struct {
		Id        bson.ObjectId
		Name      string
		Age       int
		Money     float64
		BirthDay  time.Time
		IsDeleted bool
		Location  OriginLocation
	}
	type TargetLocation struct {
		City      string
		Latitude  float32
		Longitude float32
	}
	type targetModel struct {
		Id       string
		Name     string
		Age      int
		Money    float64
		BirthDay string
		Location TargetLocation
	}

	var target targetModel
	originId := bson.NewObjectId()
	originTime := time.Now()
	origin := originModel{
		Id:        originId,
		Name:      "MockModel",
		Age:       21,
		Money:     7.7,
		BirthDay:  originTime,
		IsDeleted: false,
		Location: OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
	}

	assert.Nil(t, Instance(nil).From(origin).CopyTo(&target))
	assert.Equal(t, target, targetModel{
		Id:       originId.Hex(),
		Name:     "MockModel",
		Age:      21,
		Money:    7.7,
		BirthDay: originTime.Format(util.RFC3339),
		Location: TargetLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
	})
}

func TestCopySlice(t *testing.T) {
	type OriginLocation struct {
		City      string
		Latitude  float64
		Longitude float64
	}
	type originModel struct {
		Id        bson.ObjectId
		Name      string
		Age       int
		Money     float64
		BirthDay  time.Time
		IsDeleted bool
		Location  OriginLocation
		Detail    *OriginLocation
	}

	var targets []*originModel
	origin := originModel{
		Id:        bson.NewObjectId(),
		Name:      "MockModel",
		Age:       21,
		Money:     7.7,
		BirthDay:  time.Now(),
		IsDeleted: false,
		Location: OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
		Detail: &OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
	}
	origins := []originModel{origin, origin}

	assert.Nil(t, Instance(nil).From(origins).CopyTo(&targets))
	assert.Equal(t, 2, len(targets))
	assert.Equal(t, origin, *targets[0])
	assert.Equal(t, origin, *targets[1])
}

func TestCopyByConvertor(t *testing.T) {
	type ProtoModel struct {
		Name  string
		Age   int
		Money float64
	}
	type DbModel struct {
		Name  string
		Age   int
		Money interface{}
	}
	type originModel struct {
		Obj *ProtoModel
	}
	type targetModel struct {
		Obj DbModel
	}

	target := targetModel{}
	origin := &originModel{
		&ProtoModel{
			Name:  "MockModel",
			Age:   21,
			Money: 7.7,
		},
	}

	assert.Nil(t, Instance(NewOption().SetSkipUnsuited(false)).RegisterConverter(Target{
		From: reflect.TypeOf(&ProtoModel{}),
		To:   reflect.TypeOf(DbModel{}),
	}, func(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
		return reflect.ValueOf(from), nil
	}).From(origin).CopyTo(&target))

	target = targetModel{
		DbModel{
			Name:  "MockModel",
			Age:   21,
			Money: 7.7,
		},
	}
	origin = &originModel{}
	function := func(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
		val := from.Interface().(targetModel)
		result := originModel{
			&ProtoModel{
				Name:  val.Obj.Name,
				Age:   val.Obj.Age,
				Money: cast.ToFloat64(val.Obj.Money),
			},
		}
		return reflect.ValueOf(result), nil
	}
	assert.Nil(t, Instance(NewOption().SetSkipUnsuited(false)).RegisterConverter(Target{
		From: reflect.TypeOf(targetModel{}),
		To:   reflect.TypeOf(originModel{}),
	}, function).From(target).CopyTo(&origin))
	assert.Equal(t, origin.Obj.Money, float64(7.7))

	assert.NotNil(t, Instance(NewOption().SetSkipUnsuited(false)).RegisterConverter(Target{
		From: reflect.TypeOf(targetModel{}),
		To:   reflect.TypeOf(&originModel{}),
	}, function).From(target).CopyTo(&origin))
}

func TestCopySliceWithDiffField(t *testing.T) {
	type OriginLocation struct {
		City      string
		Latitude  float64
		Longitude float64
	}
	type originModel struct {
		Id        bson.ObjectId
		Name      string
		Age       int
		Money     float64
		BirthDay  time.Time
		IsDeleted bool
		Location  OriginLocation
	}
	type TargetLocation struct {
		City      string
		Latitude  float32
		Longitude float32
	}
	type targetModel struct {
		Id       string
		Name     string
		Age      int
		Money    float64
		BirthDay string
		Loc      *TargetLocation
	}

	originId := bson.NewObjectId()
	originTime := time.Now()
	origin := originModel{
		Id:        originId,
		Name:      "MockModel",
		Age:       21,
		Money:     7.7,
		BirthDay:  originTime,
		IsDeleted: false,
		Location: OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
	}

	var targets []*targetModel

	origins := []originModel{origin}

	assert.Nil(t, Instance(nil).RegisterResetDiffField([]DiffFieldPair{{
		Origin:  "Location",
		Targets: []string{"Loc"},
	}}).From(origins).CopyTo(&targets))
	assert.Equal(t, 1, len(targets))
	assert.Equal(t, origin.Id.Hex(), (*targets[0]).Id)
	assert.Equal(t, origin.Age, (*targets[0]).Age)
	assert.Equal(t, origin.Location.City, (*targets[0]).Loc.City)
}

func TestSkipExists(t *testing.T) {
	type OriginLocation struct {
		City      string
		Latitude  float64
		Longitude float64
	}
	type originModel struct {
		Id        bson.ObjectId
		Name      string
		Age       int
		Money     float64
		BirthDay  time.Time
		IsDeleted bool
		Location  OriginLocation
		Detail    *OriginLocation
	}

	target := originModel{}
	origin := originModel{
		Id:        bson.NewObjectId(),
		Name:      "MockModel",
		Age:       21,
		Money:     7.7,
		BirthDay:  time.Now(),
		IsDeleted: false,
		Location: OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
		Detail: &OriginLocation{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
	}

	target.Name = "I already have a name"
	assert.Nil(t, Instance(NewOption().SetOverwrite(false)).From(origin).CopyTo(&target))
	origin.Name = target.Name
	assert.Equal(t, origin, target)
}

func TestTransformerModelToProto(t *testing.T) {
	type Location struct {
		City      string
		Latitude  float64
		Longitude float64
	}
	type originModel struct {
		Id       bson.ObjectId
		Name     string
		BirthDay time.Time
		StoreId  string
	}

	type targetModel struct {
		Id         string
		TargetName string
		Name       string
		CreatedAt  string
		Location   *Location
	}

	var targets []targetModel
	locationMapper := map[string]*Location{
		"12306": &Location{
			City:      "ShangHai",
			Latitude:  234.123412,
			Longitude: 3423.43265,
		},
	}

	origins := []originModel{
		{
			Id:       bson.NewObjectId(),
			Name:     "MockModel1",
			BirthDay: time.Now(),
			StoreId:  "12306",
		},
		{
			Id:       bson.NewObjectId(),
			Name:     "MockModel2",
			BirthDay: time.Now(),
			StoreId:  "12345",
		},
	}
	assert.Nil(t, Instance(nil).
		RegisterTransformer(map[string]interface{}{
			"Location": func(storeId string) *Location {
				if location, ok := locationMapper[storeId]; ok {
					return location
				}
				return nil
			},
		}).
		RegisterResetDiffField([]DiffFieldPair{
			{Origin: "Name", Targets: []string{"TargetName", "Name"}},
			{Origin: "BirthDay", Targets: []string{"CreatedAt"}},
			{Origin: "StoreId", Targets: []string{"Location"}}},
		).From(origins).CopyTo(&targets))

	assert.Equal(t, targetModel{
		Id:         origins[0].Id.Hex(),
		TargetName: "MockModel1",
		Name:       "MockModel1",
		CreatedAt:  origins[0].BirthDay.Format(util.RFC3339),
		Location:   locationMapper["12306"],
	}, targets[0])
	assert.Equal(t, targetModel{
		Id:         origins[1].Id.Hex(),
		TargetName: "MockModel2",
		Name:       "MockModel2",
		CreatedAt:  origins[1].BirthDay.Format(util.RFC3339),
		Location:   nil,
	}, targets[1])
}

func TestCopyModelToProtoWithMultiLevelAndTransformer(t *testing.T) {
	type Age struct {
		Value int
	}
	type OriginCityInfo struct {
		Age  Age
		Area float64
	}
	type TargetCityInfo struct {
		Age  Age
		Name string
	}
	type OriginLocation struct {
		City     string
		CityInfo OriginCityInfo
	}
	type originModel struct {
		Name     string
		Location OriginLocation
	}
	type TargetLocation struct {
		City         string
		CityName     string
		CityNickName string
		CityInfo     TargetCityInfo
	}
	type targetModel struct {
		Name string
		Loc  *TargetLocation
	}

	var targets []targetModel
	origins := []originModel{
		{
			Name: "MockModel",
			Location: OriginLocation{
				City: "ShangHai",
				CityInfo: OriginCityInfo{
					Age:  Age{Value: 1},
					Area: 1,
				},
			},
		},
	}
	assert.Nil(t, Instance(nil).RegisterTransformer(map[string]interface{}{
		"Loc.CityNickName": func(city string) string {
			return "Transformer city nick name"
		},
		"Loc.CityInfo.Age": func(city Age) Age {
			city.Value++
			return city
		},
	}).RegisterResetDiffField([]DiffFieldPair{
		{Origin: "Location", Targets: []string{"Loc"}},
		{Origin: "Location.City", Targets: []string{"Loc.CityName", "Loc.CityNickName", "Name"}},
		{Origin: "Location.CityInfo.Age", Targets: []string{"Loc.CityInfo.Age"}},
	}).From(origins).CopyTo(&targets))

	assert.Equal(t, targetModel{
		Name: "ShangHai",
		Loc: &TargetLocation{
			City:         "ShangHai",
			CityName:     "ShangHai",
			CityNickName: "Transformer city nick name",
			CityInfo: TargetCityInfo{
				Age: Age{Value: 2},
			},
		},
	}, targets[0])
}

func TestCopyModelToProtoFromDeepLevel(t *testing.T) {
	type rChannel struct {
		Id     string
		Origin string
	}
	type modelChannel struct {
		Id     string
		Origin string
	}
	type ScoreInfo struct {
		Channel modelChannel
	}
	type response struct {
		Channel *rChannel
	}
	type model struct {
		Info ScoreInfo
	}

	resp := []response{}
	origin := []model{
		{
			Info: ScoreInfo{
				Channel: modelChannel{
					Id:     "123456",
					Origin: "retail",
				},
			},
		},
	}
	assert.Nil(t, Instance(nil).RegisterTransformer(Transformer{
		"Channel": func(s ScoreInfo) *rChannel {
			return &rChannel{
				Id:     s.Channel.Id,
				Origin: s.Channel.Origin,
			}
		},
	}).RegisterResetDiffField([]DiffFieldPair{
		{Origin: "Info", Targets: []string{"Channel"}},
	}).From(origin).CopyTo(&resp))

	assert.Equal(t, response{
		Channel: &rChannel{
			Id:     "123456",
			Origin: "retail",
		},
	}, resp[0])
}

func TestCopyModelToProtoWithTransformerMultiField(t *testing.T) {
	type originModel struct {
		Name     string
		Age      int
		BirthDay time.Time
	}
	type targetModel struct {
		Name    string
		NameArr []string
		Content *string
	}

	var targets []targetModel
	origins := []originModel{
		{
			Name:     "MockModel",
			BirthDay: time.Now(),
			Age:      18,
		},
	}
	assert.Nil(t, Instance(nil).RegisterResetDiffField([]DiffFieldPair{
		{Origin: "BirthDay", Targets: []string{"Name", "NameArr"}},
		{Origin: "Name", Targets: []string{"Name", "NameArr"}},
		{Origin: "Age", Targets: []string{"Content"}},
	}).RegisterTransformer(map[string]interface{}{
		"NameArr": func(value interface{}, originFieldKey string, target []string) []string {
			switch originFieldKey {
			case "BirthDay":
				target = append(target, value.(time.Time).Format(time.Kitchen))
			case "Name":
				target = append(target, value.(string))
			}
			return target
		},
		"Content": func(age int, originFieldKey string, target *string) *string {
			switch originFieldKey {
			case "Age":
				result := strconv.FormatInt(int64(age), 10)
				return &result
			}
			return nil
		},
		"Name": func(value interface{}, originFieldKey string, target string, targetKey string) string {
			switch originFieldKey {
			case "BirthDay":
				target = target + ", I was born on " + value.(time.Time).Format(time.Kitchen) + "."
			case "Name":
				target = "My name is " + value.(string)
			}
			return target
		},
	}).From(origins).CopyTo(&targets))

	age := "18"
	t.Log()
	assert.Equal(t, targetModel{
		Name:    "My name is MockModel, I was born on " + time.Now().Format(time.Kitchen) + ".",
		NameArr: []string{"MockModel", time.Now().Format(time.Kitchen)},
		Content: &age,
	}, targets[0])
}

func TestCopyModelToModelWithIgnoreInvalidOption(t *testing.T) {
	type originModel struct {
		Id        bson.ObjectId
		BirthDay  time.Time
		IsDeleted bool
	}
	type targetModel struct {
		Id        int
		BirthDay  int
		IsDeleted string
	}

	target := targetModel{}
	origin := originModel{
		Id:        bson.NewObjectId(),
		BirthDay:  time.Now(),
		IsDeleted: false,
	}

	assert.Nil(t, Instance(NewOption().SetIgnoreEmpty(true).SetOverwrite(false)).RegisterIgnoreTargetFields([]FieldKey{"Id", "BirthDay"}).From(origin).CopyTo(&target))
}

func TestCopyModelToProtoModelWithOverwriteOriginalCopyFieldOption(t *testing.T) {
	type originModel struct {
		Id        bson.ObjectId
		ProductId string
	}
	type targetModel struct {
		Id        string
		ProductId string
	}

	target := targetModel{}
	origin := originModel{
		Id:        bson.NewObjectId(),
		ProductId: "TestProductId",
	}
	copier := Instance(NewOption().SetOverwriteOriginalCopyField(true))
	assert.Nil(t, copier.RegisterResetDiffField([]DiffFieldPair{{Origin: "Id", Targets: []string{"ProductId"}}}).From(origin).CopyTo(&target))
	assert.Equal(t, targetModel{
		Id:        "",
		ProductId: origin.Id.Hex(),
	}, target)

	targets := []targetModel{}
	origins := []originModel{
		{Id: bson.NewObjectId(), ProductId: "TestProductId1"},
		{Id: bson.NewObjectId(), ProductId: "TestProductId2"},
	}
	assert.Nil(t, copier.RegisterResetDiffField([]DiffFieldPair{{Origin: "Id", Targets: []string{"Id", "ProductId"}}}).From(origins).CopyTo(&targets))
	assert.Equal(t, targetModel{
		Id:        origins[1].Id.Hex(),
		ProductId: origins[1].Id.Hex(),
	}, targets[1])
}

func TestDoubleModelIntoOneProtoModel(t *testing.T) {
	type Template struct {
		Id        string `bson:"id"`
		AppSecret bool   `bson:"appSecret"`
		Title     string `bson:"title"`
	}
	type Task struct {
		Id       bson.ObjectId `bson:"_id"`
		Name     string        `bson:"name"`
		Type     string        `bson:"type"`
		Template Template      `bson:"template"`
	}
	type TaskDetail struct {
		Id        string
		Name      string
		Type      string
		AppSecret bool
		Title     string
	}

	task := Task{
		Id:   bson.NewObjectId(),
		Name: "task",
		Type: "task",
		Template: Template{
			Id:        "templateId",
			AppSecret: true,
			Title:     "templateTitle",
		},
	}
	taskDetail := TaskDetail{}
	err := Instance(nil).RegisterResetDiffField([]DiffFieldPair{
		{
			Origin:  "Template.Id",
			Targets: []string{"Id"},
		},
		{
			Origin:  "Template.Title",
			Targets: []string{"Title"},
		},
		{
			Origin:  "Template.AppSecret",
			Targets: []string{"AppSecret"},
		},
	}).From(task).CopyTo(&taskDetail)
	assert.Nil(t, err)
	assert.Equal(t, task.Template.Id, taskDetail.Id)
	assert.Equal(t, task.Template.Title, taskDetail.Title)
	assert.Equal(t, task.Template.AppSecret, taskDetail.AppSecret)
}

func TestCopierUnexportedField(t *testing.T) {
	type incScoreOption struct {
		Name string
	}
	type ScoreInfo struct {
		Score          int
		Description    *string
		businessId     string
		reason         *string
		incScoreOption *incScoreOption
	}

	scoreDescription := "inc by register"
	reason := "task"
	scoreInfo := &ScoreInfo{
		Score:       12,
		Description: &scoreDescription,
		businessId:  "2022-03-09",
		reason:      &reason,
		incScoreOption: &incScoreOption{
			Name: "origin",
		},
	}

	result := ScoreInfo{}
	err := Instance(NewOption().SetCopyUnexported(true)).From(scoreInfo).CopyTo(&result)

	assert.Nil(t, err)
	assert.Equal(t, result.incScoreOption.Name, scoreInfo.incScoreOption.Name)
	assert.Equal(t, result.Score, scoreInfo.Score)
	assert.Equal(t, *result.Description, *scoreInfo.Description)
	assert.Equal(t, result.Description, scoreInfo.Description)
	assert.Equal(t, result.reason, scoreInfo.reason)
	assert.NotSame(t, result.Description, scoreInfo.Description)
	assert.NotSame(t, result.reason, scoreInfo.reason)
	assert.Equal(t, result.businessId, scoreInfo.businessId)
}

func TestCopierWithIgnoreDeepEmpty(t *testing.T) {
	type Info struct {
		Sku string
	}
	type Tag struct {
		TagName string
	}
	type EcProduct struct {
		Info *Info
		Tags *[]Tag
	}

	from := &EcProduct{
		Info: &Info{},
		Tags: &[]Tag{},
	}
	to := &EcProduct{}
	err := Instance(NewOption().SetIgnoreDeepEmpty(true)).From(from).CopyTo(to)

	assert.Nil(t, err)
	assert.Nil(t, to.Info)
	assert.Nil(t, to.Tags)

	err = Instance(NewOption().SetIgnoreDeepEmpty(false)).From(from).CopyTo(to)

	assert.Nil(t, err)
	assert.Equal(t, to.Info.Sku, "")
	assert.NotNil(t, to.Tags)
}
