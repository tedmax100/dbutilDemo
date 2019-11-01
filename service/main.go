package main 

import (
	"context"

	"sparrow/sparrow/db"
	"math/rand"
	"time"
	"fmt"
	//  "sync"
)
const charset = "abcdefghijklmnopqrstuvwxyz1234567890"
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset() string {
	length := rand.Intn((15-6) + 6)

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type User struct {
	ID int64
	Name string
	Balance int
}

type Tran struct {
	ID int64 `sql:"id"`
	UserID int64`sql:"user_id"`
	Amount int 
	Type uint
	Time  int64
}

func main() {
	cancelCtx, _ := context.WithCancel(context.Background())
	dbClient := db.NewDatabase(cancelCtx)

	dbCtx := dbClient.Open(cancelCtx)

	user := User{
		ID : time.Now().UnixNano() / 1000 / 1000 / 1000,
		Name : StringWithCharset(),
		Balance : 1000,
	}
	insertCom := dbCtx.NewCommand(cancelCtx)
	_, err := insertCom.For("user").Insert("id", user.ID).Insert("name", user.Name).Insert("balance", user.Balance).Exec()
	if err != nil {
		fmt.Println(err)
	}

/* 	var wg sync.WaitGroup 
	for k:= 1 ; k < 100 ; k++ {
		wg.Add(5)
		for i := 0 ; i< 5; i ++ {
			 go func() {
				defer wg.Done()
				tranId, err := dbClient.NewID(cancelCtx, "tran", "deposit")
				if err != nil {
					fmt.Println(err)
				}
				insertCom := dbCtx.NewCommand(cancelCtx)
			
				tran := Tran{
					ID : tranId,
					UserID: user.ID,
					Amount: int(rand.Intn((100-6) + 6)),
					Type : 1,
					Time : time.Now().UnixNano() / 1000 / 1000/ 1000,
				}
		
				_, err = insertCom.For("tran").Insert("id", tran.ID).Insert("user_id", tran.UserID).Insert("amount", tran.Amount).Insert("type", tran.Type).Insert("time", tran.Time).Exec()
				if err != nil {
					fmt.Println(err)
				}
			 }()
			}
		wg.Wait()
		 time.Sleep(10 * time.Millisecond)
	} */
	
	var tranID int64
	tranID, err = dbClient.NewID(cancelCtx, "tran", "deposit")
	if err != nil {
		fmt.Println(err)
	}
	insertCom = dbCtx.NewCommand(cancelCtx)
			


	tran := Tran{
		ID : tranID,
		UserID: user.ID,
		Amount: int(rand.Intn((100-6) + 6)),
		Type : 1,
		Time : time.Now().UnixNano() / 1000 / 1000/ 1000,
	}

	insertCom.For("tran").Insert("id", tran.ID).Insert("user_id", tran.UserID).Insert("amount", tran.Amount).Insert("type", tran.Type).Insert("time", tran.Time)
	err = dbCtx.Begin(insertCom)
	if err != nil {
		fmt.Println(err)
	} 
	_, err = insertCom.Exec()
	if err != nil {
		fmt.Println(err)
		dbCtx.Rollback()
	} 

	queryUser := User{}
	queryCom := dbCtx.NewQuery(cancelCtx)
	queryCom.For("user").RawSQL("SELECT id, name, balance FROM user").Where("id={id}", "id")
	vars := map[string]interface{}{"id": user.ID}
	queryCom.Vars(vars)
	err = queryCom.Find(&queryUser)
	if err != nil {
		dbCtx.Rollback()
	}
	dbCtx.Commit()
}