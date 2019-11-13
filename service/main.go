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

type Sport struct {
	ID   int64
	Name string
}

type Match struct {
	ID       int64
	SportID  int64  `sql:"sport_id"`
	HomeName string `sql:"home_name"`
	AwayName string `sql:"away_name"`
	Time     int64
}

type Odds struct {
	ID          int64
	MatchID     int64  `sql:"match_id"`
	BetTypeName string `sql:"bet_type_name"`
	Selection1  float64
	Selection2  float64
	Time        int64
}

func main() {
	cancelCtx, _ := context.WithCancel(context.Background())
	dbClient := db.NewDatabase(cancelCtx)

	dbCtx := dbClient.Open(cancelCtx)

	// 插入兩個同業務不同產品的match table, match有參與分表
	soccerMatchID, _ := dbClient.NewID(cancelCtx, "match", "soccer")
	fmt.Println("soccer match id :  ", soccerMatchID)
	match01 := Match{
		ID:       soccerMatchID,
		SportID:  1,
		HomeName: "Chelsea",
		AwayName: "Arsenal",
		Time:     time.Now().UnixNano() / 1000 / 1000 / 1000,
	}
	insertSoccerMatchCom := dbCtx.NewCommand(cancelCtx)
	insertSoccerMatchCom.For("match").Insert("id", match01.ID).Insert("sport_id", match01.SportID).Insert("home_name", match01.HomeName).Insert("away_name", match01.AwayName).Insert("time", match01.Time)
	_, err := insertSoccerMatchCom.Exec()
	if err != nil {
		fmt.Println(err)
	}

	basketballMatchID, _ := dbClient.NewID(cancelCtx, "match", "basketball")
	fmt.Println("basketball match id : ", basketballMatchID)
	match02 := Match{
		SportID:  basketballMatchID,
		HomeName: "Laker",
		AwayName: "King",
		Time:     time.Now().UnixNano() / 1000 / 1000 / 1000,
	}
	insertBasketballMatchCom := dbCtx.NewCommand(cancelCtx)
	insertBasketballMatchCom.For("match").Insert("id", match02.ID).Insert("sport_id", match02.SportID).Insert("home_name", match02.HomeName).Insert("away_name", match02.AwayName).Insert("time", match02.Time)
	_, err = insertBasketballMatchCom.Exec()
	if err != nil {
		fmt.Println(err)
	}

	// select 單場 match
	singleMatch := Match{}
	singleMatchSelCom := dbCtx.NewQuery(cancelCtx)
	singleMatchSelCom.For("match").SelectModel(&singleMatch).Where("id={id}", "id")
	singleMatchSelCom.Vars(map[string]interface{}{"id": match01.ID})
	err = singleMatchSelCom.Find(&singleMatch)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("singleMatch: %+v\n", singleMatch)

	// select 運動列表
	sportList := []Sport{}
	sportListSelCom := dbCtx.NewQuery(cancelCtx)
	err = sportListSelCom.For("sport").RawSQL("SELECT * FROM sport").Find(&sportList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("sport list: %+v\n", sportList)

	//  透過parent 的分片ID, 來生成同庫的其他業務分片ID
	oddsiD, _ := dbClient.NewSubID(cancelCtx, "odds", "soccer", match01.ID)
	odds := Odds{
		ID:          oddsiD,
		MatchID:     match01.ID,
		BetTypeName: "hadicap",
		Selection1:  1.95,
		Selection2:  2.02,
		Time:        time.Now().UnixNano() / 1000 / 1000 / 1000,
	}
	fmt.Printf("oddsid: %+v\n", odds)

	// 啟動事物
	soccerMatchID, _ = dbClient.NewID(cancelCtx, "match", "soccer")
	fmt.Println("soccer match id :  ", soccerMatchID)

	match01 = Match{
		ID:       soccerMatchID,
		SportID:  1,
		HomeName: "Crystal",
		AwayName: "Manchester City",
		Time:     time.Now().UnixNano() / 1000 / 1000 / 1000,
	}
	insertSoccerMatchCom = dbCtx.NewCommand(cancelCtx)
	insertSoccerMatchCom.For("match").Insert("id", match01.ID).Insert("sport_id", match01.SportID).Insert("home_name", match01.HomeName).Insert("away_name", match01.AwayName).Insert("time", match01.Time)
	err = dbCtx.Begin(insertSoccerMatchCom)
	if err != nil {
		fmt.Println(err)
	}
	_, err = insertSoccerMatchCom.Exec()
	if err != nil {
		dbCtx.Rollback()
		fmt.Println(err)
	}

	oddsiD, _ = dbClient.NewSubID(cancelCtx, "odds", "soccer", match01.ID)
	odds = Odds{
		ID:          oddsiD,
		MatchID:     match01.ID,
		BetTypeName: "over under",
		Selection1:  1.88,
		Selection2:  2.12,
		Time:        time.Now().UnixNano() / 1000 / 1000 / 1000,
	}

	insertSoccerOddsCom := dbCtx.NewCommand(cancelCtx)
	insertSoccerOddsCom.For("odds").Insert("id", odds.ID).Insert("match_id", odds.MatchID).Insert("bet_type_name", odds.BetTypeName).Insert("selection1", odds.Selection1).Insert("selection2", odds.Selection2).Insert("time", odds.Time)
	_, err = insertSoccerOddsCom.Exec()
	if err != nil {
		// 因為odds表格不存在, 所以整個事物被回滾, 剛剛還沒commit的match insert操作也被回滾了
		dbCtx.Rollback()
		fmt.Println(err)
	} else {
		dbCtx.Commit()
	}

	// end
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
	err = dbCtx.Begin(insertComPR)
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
