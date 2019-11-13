package main

import (
	"context"

	"fmt"
	"math/rand"
	"sparrow/sparrow/db"
	"time"
	//  "sync"
)

const charset = "abcdefghijklmnopqrstuvwxyz1234567890"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset() string {
	length := rand.Intn((15 - 6) + 6)

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type User struct {
	ID      int64
	Name    string
	Balance int
}

type Base struct {
	ID int64 `sql:"id"`
}
type Tran struct {
	//ID     int64
	Base
	UserID int64 `sql:"user_id"`
	Amount int
	Type   uint
	Time   int64
}

func main() {
	cancelCtx, _ := context.WithCancel(context.Background())
	dbClient := db.NewDatabase(cancelCtx)

	dbCtx := dbClient.Open(cancelCtx)

	user := User{
		ID:      time.Now().UnixNano() / 1000 / 1000 / 1000,
		Name:    StringWithCharset(),
		Balance: 1000,
	}
	users := []User{
		/* 		ID:      time.Now().UnixNano() / 1000 / 1000 / 1000,
		   		Name:    StringWithCharset(),
		   		Balance: 1000, */
	}

	insertComPR := dbCtx.NewCommand(cancelCtx)
	insertComPR.For("user").RawSQL("select * from user")
	err := dbCtx.Begin(insertComPR)
	insertComPR.Find(&users)

	dbCtxPR2 := dbClient.Open(cancelCtx)
	insertComPR2 := dbCtxPR2.NewCommand(cancelCtx)
	insertComPR2.For("user").Insert("id", 21).Insert("name", "cc").Insert("balance", 500)
	err = dbCtxPR2.Begin(insertComPR2)
	_, err = insertComPR2.Exec()
	if err != nil {
		fmt.Println(err)
	}
	users = []User{}
	insertComPR.Find(&users)
	fmt.Println(len(users))
	dbCtxPR2.Commit()
	users = []User{}
	insertComPR.Find(&users)
	fmt.Println(len(users))
	// insertComPR3 := dbCtx.NewCommand(cancelCtx)
	insertComPR.For("user").Insert("id", 21).Insert("name", "dd").Insert("balance", 400)
	_, err = insertComPR.Exec()
	if err != nil {
		fmt.Println(err)
	}
	err = dbCtx.Commit()
	if err != nil {
		fmt.Println(err)
	}
	insertCom := dbCtx.NewCommand(cancelCtx)
	insertCom.For("user").RawSQL("select * from user")
	err = dbCtx.Begin(insertCom)
	insertCom.Find(&users)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(users))

	dbCtx2 := dbClient.Open(cancelCtx)
	insertCom2 := dbCtx2.NewCommand(cancelCtx)
	insertCom2.For("user").Insert("id", user.ID).Insert("name", user.Name).Insert("balance", user.Balance)
	err = dbCtx2.Begin(insertCom2)
	_, err = insertCom2.Exec()
	if err != nil {
		fmt.Println(err)
	}

	users2 := []User{}
	insertCom.Find(&users2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(users2))
	dbCtx2.Commit()

	_, err = insertCom.For("user").Insert("id", user.ID).Insert("name", user.Name).Insert("balance", user.Balance).Exec()

	/* 	if err != nil {
	   		fmt.Println(err)
	   	}
	*/
	var tranID int64
	tranID, err = dbClient.NewID(cancelCtx, "tran", "deposit")
	if err != nil {
		fmt.Println(err)
	}
	insertCom = dbCtx.NewCommand(cancelCtx)

	tran := Tran{
		Base: Base{
			ID: tranID,
		},
		//ID:     tranID,
		UserID: user.ID,
		Amount: int(rand.Intn((100 - 6) + 6)),
		Type:   1,
		Time:   time.Now().UnixNano() / 1000 / 1000 / 1000,
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
	if queryUser.Balance-550 >= 0 {
		tranID, err = dbClient.NewID(cancelCtx, "tran", "deposit")
		if err != nil {
			fmt.Println(err)
		}
		tran = Tran{
			Base: Base{
				ID: tranID,
			},
			// ID:     tranID,
			UserID: user.ID,
			Amount: 550,
			Type:   1,
			Time:   time.Now().UnixNano() / 1000 / 1000 / 1000,
		}
		insertCom.For("tran").Insert("id", tran.ID).Insert("user_id", tran.UserID).Insert("amount", tran.Amount).Insert("type", tran.Type).Insert("time", tran.Time)
		_, err = insertCom.Exec()
		if err != nil {
			fmt.Println(err)
			dbCtx.Rollback()
		} else {
			updateCom := dbCtx.NewCommand(cancelCtx)
			updateCom.For("user").Update("balance", queryUser.Balance-550).Where("id={id}", "id")
			vars := map[string]interface{}{"id": 1234512}
			updateCom.Vars(vars)
			updateResult, err := updateCom.Exec()
			if err != nil {
				fmt.Println(err)
				dbCtx.Rollback()
			} else {
				rowAffected, err := updateResult.RowsAffected()
				if err != nil {
					fmt.Println(err)
					dbCtx.Rollback()
				}
				if rowAffected == 0 {
					dbCtx.Rollback()
				} else {
					dbCtx.Commit()
				}

			}
		}

	}

}
