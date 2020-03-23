package repositories

import (
	"fmt"
	"sync"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"flower/common"
	"flower/config"
)

type DB struct {
	Engine *xorm.Engine
}

func (session *DB) Builder() *Builder {
	b := builder.Dialect(session.Engine.DriverName())
	return &Builder{b}
}
func (DB *DB) GetNewSession() *Session {
	session := DB.Engine.NewSession()
	return &Session{
		Session: session,
		engine:  DB.Engine,
	}
}
func (DB *DB) GetNewTransactionSession() *Session {
	session := DB.Engine.NewSession()
	return &Session{
		Session: session,
		engine:  DB.Engine,
	}
}
func (DB *DB) Transaction(f func(session *Session) (interface{}, error)) (interface{}, error) {
	return DB.Engine.Transaction(func(s *xorm.Session) (interface{}, error) {
		return f(&Session{
			Session: s,
			engine:  DB.Engine,
		})
	})
}

func WhereNoDelete(session *xorm.Session) *xorm.Session {
	session.Where("status!=?", common.TableStatus_Delete)
	return session
}

var dBSession *DB
var biSession *DB
var once sync.Once
var bionce sync.Once

func initEngine(config *config.Config, dbconfig *config.DBConfig) (*xorm.Engine, error) {
	//postgres://postgres:123456@192.168.2.163/deepface_v5
	connstr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbconfig.DBUser, dbconfig.DBPasswd, dbconfig.DBAddr, dbconfig.DBPort, dbconfig.Database)
	fmt.Println(connstr)
	engine, err := xorm.NewEngine("postgres", connstr)
	if err != nil {
		glog.Fatalf("init db error - %s", err.Error())
		return engine, err
	}
	//engine.SetLogger()
	engine.SetMaxOpenConns(int(config.MaxConnSize))
	engine.SetMaxIdleConns(int(config.MinIdleConnSize))
	//cacher := utils.NewRedisCacher(1*time.Minute, nil)
	//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	//engine.SetDefaultCacher(cacher)
	//engine.MapCacher("",cacher)
	// engine.MapCacher(&model.CivilAttrs{}, nil)
	// engine.MapCacher(&model.CivilImages{}, nil)
	engine.ShowExecTime(false)
	engine.ShowSQL(false)
	if config.Debug {
		engine.ShowExecTime(true)
		engine.ShowSQL(true)
	}
	return engine, nil
}

func GetInstance() *DB {
	once.Do(func() {
		dBSession = new(DB)
		var err error
		if dBSession.Engine, err = initEngine(config.GetConfig(), config.GetConfig().DbData); err != nil {
			panic(err.Error())
		}
	})
	return dBSession
}

func GetBiDB() *DB {
	bionce.Do(func() {
		biSession = new(DB)
		var err error
		if biSession.Engine, err = initEngine(config.GetConfig(), config.GetConfig().DbBi); err != nil {
			panic(err.Error())
		}
	})
	return biSession
}

var dbDeepSession *DB
var deepOnce sync.Once

func GetDeepDB() *DB {
	deepOnce.Do(func() {
		dbDeepSession = new(DB)
		var err error
		if dbDeepSession.Engine, err = initEngine(config.GetConfig(), config.GetConfig().DbDeep); err != nil {
			panic(err.Error())
		}
	})
	return dbDeepSession
}
