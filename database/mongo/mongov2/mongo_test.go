package mongov2

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Trainer 도메인 모델
type Trainer struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Age  int    `bson:"age" json:"age"`
	City string `bson:"city" json:"city"`
}

type mongoTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
	server     *memongo.Server
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(mongoTestSuite))
}

func (s *mongoTestSuite) SetupSuite() {
	s.ctx = context.Background()

	// memongo: 인메모리 MongoDB 시작
	server, err := memongo.Start("6.0.0")
	if err != nil {
		// memongo가 실패하면 (mongod 바이너리 없을 수 있음) 로컬 MongoDB 시도
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			uri = "mongodb://localhost:27017"
		}
		client, err := mongo.Connect(options.Client().ApplyURI(uri))
		s.Require().NoError(err)
		s.client = client
		s.db = client.Database("test_mongov2")
	} else {
		s.server = server
		client, err := mongo.Connect(options.Client().ApplyURI(server.URI()))
		s.Require().NoError(err)
		s.client = client
		s.db = client.Database(memongo.RandomDatabase())
	}

	s.collection = s.db.Collection("trainers")
}

func (s *mongoTestSuite) TearDownTest() {
	s.collection.Drop(s.ctx)
}

func (s *mongoTestSuite) TearDownSuite() {
	if s.client != nil {
		s.client.Disconnect(s.ctx)
	}
	if s.server != nil {
		s.server.Stop()
	}
}

// --- CRUD 테스트 ---

func (s *mongoTestSuite) TestInsertOne_FindOne() {
	trainer := Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"}

	// InsertOne
	result, err := s.collection.InsertOne(s.ctx, trainer)
	s.NoError(err)
	s.Equal("t1", result.InsertedID)

	// FindOne
	var found Trainer
	err = s.collection.FindOne(s.ctx, bson.D{{Key: "_id", Value: "t1"}}).Decode(&found)
	s.NoError(err)
	s.Equal("Ash", found.Name)
	s.Equal(10, found.Age)
}

func (s *mongoTestSuite) TestInsertMany_FindAll() {
	trainers := []interface{}{
		Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"},
		Trainer{ID: "t2", Name: "Misty", Age: 12, City: "Cerulean City"},
		Trainer{ID: "t3", Name: "Brock", Age: 15, City: "Pewter City"},
	}

	// InsertMany
	result, err := s.collection.InsertMany(s.ctx, trainers)
	s.NoError(err)
	s.Len(result.InsertedIDs, 3)

	// Find (전체 조회)
	cursor, err := s.collection.Find(s.ctx, bson.D{})
	s.NoError(err)

	var results []Trainer
	err = cursor.All(s.ctx, &results)
	s.NoError(err)
	s.Len(results, 3)
}

func (s *mongoTestSuite) TestFind_WithOptions() {
	// 데이터 삽입
	for i := 0; i < 5; i++ {
		s.collection.InsertOne(s.ctx, Trainer{
			ID:   "t" + string(rune('1'+i)),
			Name: "Trainer" + string(rune('A'+i)),
			Age:  10 + i,
		})
	}

	// Limit + Sort 옵션
	opts := options.Find().SetLimit(3).SetSort(bson.D{{Key: "age", Value: -1}})
	cursor, err := s.collection.Find(s.ctx, bson.D{}, opts)
	s.NoError(err)

	var results []Trainer
	err = cursor.All(s.ctx, &results)
	s.NoError(err)
	s.Len(results, 3)
	s.Equal(14, results[0].Age) // 가장 나이 많은 순
}

func (s *mongoTestSuite) TestUpdateOne() {
	s.collection.InsertOne(s.ctx, Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"})

	// $set 연산자로 업데이트
	_, err := s.collection.UpdateOne(s.ctx,
		bson.D{{Key: "_id", Value: "t1"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 11}}}},
	)
	s.NoError(err)

	var updated Trainer
	s.collection.FindOne(s.ctx, bson.D{{Key: "_id", Value: "t1"}}).Decode(&updated)
	s.Equal(11, updated.Age)

	// $inc 연산자로 증가
	_, err = s.collection.UpdateOne(s.ctx,
		bson.D{{Key: "_id", Value: "t1"}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "age", Value: 5}}}},
	)
	s.NoError(err)

	s.collection.FindOne(s.ctx, bson.D{{Key: "_id", Value: "t1"}}).Decode(&updated)
	s.Equal(16, updated.Age)
}

func (s *mongoTestSuite) TestDeleteOne() {
	s.collection.InsertOne(s.ctx, Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"})

	result, err := s.collection.DeleteOne(s.ctx, bson.D{{Key: "_id", Value: "t1"}})
	s.NoError(err)
	s.Equal(int64(1), result.DeletedCount)

	// 삭제 확인
	err = s.collection.FindOne(s.ctx, bson.D{{Key: "_id", Value: "t1"}}).Decode(&Trainer{})
	s.ErrorIs(err, mongo.ErrNoDocuments)
}

func (s *mongoTestSuite) TestBsonM_Filter() {
	s.collection.InsertOne(s.ctx, Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"})
	s.collection.InsertOne(s.ctx, Trainer{ID: "t2", Name: "Misty", Age: 12, City: "Cerulean City"})

	// bson.M (map 형태 필터)
	var found Trainer
	err := s.collection.FindOne(s.ctx, bson.M{"name": "Misty"}).Decode(&found)
	s.NoError(err)
	s.Equal("Cerulean City", found.City)
}

// --- 부분 업데이트 테스트 ---

func (s *mongoTestSuite) TestPartialUpdate_WithFlatbson() {
	s.collection.InsertOne(s.ctx, Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"})

	// 이름만 업데이트 (다른 필드는 유지)
	_, err := s.collection.UpdateOne(s.ctx,
		bson.D{{Key: "_id", Value: "t1"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: "Ash Ketchum"}}}},
	)
	s.NoError(err)

	var updated Trainer
	s.collection.FindOne(s.ctx, bson.D{{Key: "_id", Value: "t1"}}).Decode(&updated)
	s.Equal("Ash Ketchum", updated.Name)
	s.Equal(10, updated.Age)           // 유지
	s.Equal("Pallet Town", updated.City) // 유지
}

// --- 인덱스 관리 테스트 ---

func (s *mongoTestSuite) TestCreateIndex_Single() {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}}, // 오름차순
	}

	indexName, err := s.collection.Indexes().CreateOne(s.ctx, indexModel)
	s.NoError(err)
	assert.Contains(s.T(), indexName, "name")
}

func (s *mongoTestSuite) TestCreateIndex_Compound() {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "age", Value: -1},
		},
	}

	indexName, err := s.collection.Indexes().CreateOne(s.ctx, indexModel)
	s.NoError(err)
	s.NotEmpty(indexName)
}

func (s *mongoTestSuite) TestCreateIndex_Unique() {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := s.collection.Indexes().CreateOne(s.ctx, indexModel)
	s.NoError(err)

	// 같은 name으로 두 번 삽입 시 에러
	s.collection.InsertOne(s.ctx, Trainer{ID: "t1", Name: "Ash", Age: 10})
	_, err = s.collection.InsertOne(s.ctx, Trainer{ID: "t2", Name: "Ash", Age: 15})
	s.Error(err) // duplicate key error
}

func (s *mongoTestSuite) TestListIndexes() {
	// 인덱스 생성
	s.collection.Indexes().CreateOne(s.ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}},
	})

	// 인덱스 목록 조회
	cursor, err := s.collection.Indexes().List(s.ctx)
	s.NoError(err)

	var indexes []bson.M
	err = cursor.All(s.ctx, &indexes)
	s.NoError(err)
	s.GreaterOrEqual(len(indexes), 2) // _id 기본 인덱스 + name 인덱스
}

func (s *mongoTestSuite) TestDropIndex() {
	// 인덱스 생성
	s.collection.Indexes().CreateOne(s.ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "city", Value: 1}},
	})

	// 인덱스 삭제
	err := s.collection.Indexes().DropOne(s.ctx, "city_1")
	s.NoError(err)
}

// --- Aggregation Pipeline 테스트 ---

func (s *mongoTestSuite) TestAggregate_GroupSum() {
	// 데이터 삽입
	s.collection.InsertMany(s.ctx, []interface{}{
		Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"},
		Trainer{ID: "t2", Name: "Gary", Age: 10, City: "Pallet Town"},
		Trainer{ID: "t3", Name: "Misty", Age: 12, City: "Cerulean City"},
		Trainer{ID: "t4", Name: "Brock", Age: 15, City: "Pewter City"},
		Trainer{ID: "t5", Name: "Erika", Age: 18, City: "Cerulean City"},
	})

	// $group + $sum: 도시별 트레이너 수
	pipeline := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$city"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	cursor, err := s.collection.Aggregate(s.ctx, pipeline)
	s.NoError(err)

	var results []bson.M
	err = cursor.All(s.ctx, &results)
	s.NoError(err)
	s.Len(results, 3) // 3개 도시

	// Pallet Town: 2명, Cerulean City: 2명, Pewter City: 1명
	s.Equal(int32(2), results[0]["count"])
}

func (s *mongoTestSuite) TestAggregate_MatchSort() {
	s.collection.InsertMany(s.ctx, []interface{}{
		Trainer{ID: "t1", Name: "Ash", Age: 10, City: "Pallet Town"},
		Trainer{ID: "t2", Name: "Misty", Age: 12, City: "Cerulean City"},
		Trainer{ID: "t3", Name: "Brock", Age: 15, City: "Pewter City"},
		Trainer{ID: "t4", Name: "Erika", Age: 18, City: "Celadon City"},
	})

	// $match + $sort: 나이 12 이상, 나이 내림차순
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 12}}}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "age", Value: -1}}}},
	}

	cursor, err := s.collection.Aggregate(s.ctx, pipeline)
	s.NoError(err)

	var results []Trainer
	err = cursor.All(s.ctx, &results)
	s.NoError(err)
	s.Len(results, 3)
	s.Equal("Erika", results[0].Name) // 18세
	s.Equal("Brock", results[1].Name) // 15세
	s.Equal("Misty", results[2].Name) // 12세
}
