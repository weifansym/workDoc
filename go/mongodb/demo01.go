package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

type MongoConfig struct {
	Host      string
	PoolLimit int
	DbName    string
	DbCol     string
}

type Config struct {
	Addr    string
	Mongodb MongoConfig
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("read filename error: ", err)
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println("json Unmarshal error: ", err)
		return
	}
}

func init() {
	JsonParse := NewJsonStruct()
	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
	JsonParse.Load("./config.json", &config)
}

var config = Config{}

var session *mgo.Session
var database *mgo.Database


type dataNews struct {
	nick    string
	uid     int
	title   string
	content string
	star    int
}

type Person struct {
	ID       bson.ObjectId `bson:"_id"`
	Name  string
	Phone string `bson:"phone"`
	CreateTime string
}


func main() {
	fmt.Println("BBBBB")
	fmt.Println(config.Mongodb)
	var err error

	dialInfo := &mgo.DialInfo{
		Addrs:     []string{config.Mongodb.Host},
		Direct:    false,
		Timeout:   time.Second * 1,
		PoolLimit: config.Mongodb.PoolLimit, // Session.SetPoolLimit
	}
	//创建一个维护套接字池的session
	session, err = mgo.DialWithInfo(dialInfo)
	defer session.Close()

	if err != nil {
		log.Println(err.Error())
	}
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(config.Mongodb.DbName).C("people")
	/*
	//  增（插入）
	err = c.Insert(&Person{bson.NewObjectId(),"Ale", "+55 53 8116 9639"},
		&Person{bson.NewObjectId(),"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}*/



	/*
	// 更新
	idStr := "5b87d053705715430ea84ed5"
	err = c.UpdateId(bson.ObjectIdHex(idStr), bson.M{"$set": bson.M{"CreateTime": time.Now().Add(8 * time.Hour), "phone": "1111111"}})
	//err = c.Update(bson.M{"_id": bson.ObjectIdHex(idStr)}, bson.M{"$set": bson.M{"phone": "18015582925"}})
	if err != nil {
		log.Fatal("update person error: ", err)
	}*/

	//  查询
	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("person : ", result)
	fmt.Println("person id: ", result.ID.Hex())


	/*
	//使用指定数据库
	database = session.DB(config.Mongodb.DbName)
	coll := database.C(config.Mongodb.DbCol)
	err = coll.Insert(&dataNews{
		"duhuo",
		12345,
		"测试demo01",
		"龙哥被砍死了",
		11,
	})
	if err != nil {
		log.Fatal(err)
	}

	result := dataNews{}
	err = coll.Find(bson.M{"nick": "duhuo"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("news data :", result)*/
}
