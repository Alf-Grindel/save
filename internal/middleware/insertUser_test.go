package middleware

import (
	"github.com/Alf_Grindel/save/conf"
	"github.com/Alf_Grindel/save/internal/dal"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/middleware/redis"
	"github.com/Alf_Grindel/save/internal/model"
	"github.com/Alf_Grindel/save/pkg/utils"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Init() {
	conf.LoadConfig()
	dal.Init()
	redis.Init()
}

const (
	insertNum = 100000

	worker    = 5
	benchSize = 5000
)

var (
	snow = utils.NewSnowflake(0)

	users []*model.User
	wg    sync.WaitGroup
)

func TestDoInsertUser(t *testing.T) {
	Init()
	start := time.Now()

	for i := 0; i < insertNum; i++ {
		id := snow.GenerateID()
		account := "fakeSave" + strconv.Itoa(i)
		user := &model.User{
			Id:       id,
			Account:  account,
			Password: "12345678",
			UserName: "fakeSave",
			Avatar:   "https://636f-codenav-8grj8px727565176-1256524210.tcb.qcloud.la/img/logo.png",
			Profile:  "test - insert user",
			Tags:     "[]",
		}
		users = append(users, user)
		if len(users) == benchSize {
			db.DB.Select("id", "account", "password", "user_name", "avatar", "profile", "tags").CreateInBatches(&users, benchSize)
			users = nil
		}
		//db.DB.Select("id", "account", "password", "user_name", "avatar", "profile", "tags").Create(&user)
	}
	if len(users) > 0 {
		db.DB.Select("id", "account", "password", "user_name", "avatar", "profile", "tags").CreateInBatches(&users, benchSize)
	}

	duration := time.Since(start)
	t.Log(duration.Milliseconds())
	/*
		10万条数据 benchsize = 1000  time: 2992 ms
		10万条数据 benchsize = 10000  time: 1422 ms
		10万条数据 benchsize = 50000  time: 1147ms
	*/
}

func TestDoInsertUserByConcurrent(t *testing.T) {
	Init()

	start := time.Now()

	batchChan := make(chan []*model.User, worker)

	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for batch := range batchChan {
				if len(batch) == 0 {
					continue
				}
				res := db.DB.Select("id", "account", "password", "user_name", "avatar", "profile", "tags").CreateInBatches(&batch, len(batch))
				if err := res.Error; err != nil {
					t.Fatal(err)
				}
			}
		}(i)
	}

	batch := make([]*model.User, 0, benchSize)

	j := 9

	for i := 0; i < insertNum; i++ {
		user := &model.User{
			Id:       snow.GenerateID(),
			Account:  "fakeSave" + strconv.Itoa(i+100000*j),
			Password: "12345678",
			UserName: "fakeSave",
			Avatar:   "https://636f-codenav-8grj8px727565176-1256524210.tcb.qcloud.la/img/logo.png",
			Profile:  "test - insert user",
			Tags:     "[]",
		}
		batch = append(batch, user)
		if len(batch) == benchSize {
			batchChan <- batch
			batch = make([]*model.User, 0, benchSize)
		}
	}

	if len(batch) > 0 {
		batchChan <- batch
	}
	close(batchChan)

	wg.Wait()

	duration := time.Since(start)
	t.Log(duration.Milliseconds())
	/*
		10万条数据 benchsize = 5000  time: 890 ms
	*/
}
