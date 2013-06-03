package models

import (
	"../../utils"
	"fmt"
	"github.com/coocood/qbs"
	_ "github.com/mattn/go-sqlite3"
	"time"
	//_ "github.com/lib/pq"
	"os"
)

const (
	DbName         = "./data/sqlite.db"
	DbUser         = "root"
	mysqlDriver    = "mymysql"
	mysqlDrvformat = "%v/%v/"
	pgDriver       = "postgres"
	pgDrvFormat    = "user=%v dbname=%v sslmode=disable"
	sqlite3Driver  = "sqlite3"
	dbtypeset      = "sqlite"
)

type File struct {
	Id              int64
	Cid             int64 `qbs:"index"`
	Nid             int64 `qbs:"index"`
	Uid             int64 `qbs:"index"`
	Pid             int64 `qbs:"index"`
	Ctype           int64
	Filename        string
	Content         string
	Hash            string
	Location        string
	Url             string
	Size            int64
	Created         time.Time `qbs:"index"`
	Updated         time.Time `qbs:"index"`
	Hotness         float64   `qbs:"index"`
	Hotup           int64     `qbs:"index"`
	Hotdown         int64     `qbs:"index"`
	Views           int64     `qbs:"index"`
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
}

// k/v infomation
type Kvs struct {
	Id int64
	/*
		Cid int64
		Nid int64
		Tid int64
		Rid int64
	*/
	K string
	V string
}

func RegisterDb() {

	switch {
	case dbtypeset == "sqlite":
		qbs.Register("sqlite3", "./data/sqlite.db", "", qbs.NewSqlite3())

	case dbtypeset == "mysql":
		qbs.Register("mysql", "qbs_test@/qbs_test?charset=utf8&parseTime=true&loc=Local", "dbname", qbs.NewMysql())

	case dbtypeset == "pgsql":
		qbs.Register("postgres", "qbs_test@/qbs_test?charset=utf8&parseTime=true&loc=Local", "dbname", qbs.NewPostgres())
	}

}

func ConnDb() (q *qbs.Qbs, err error) {
	RegisterDb()
	q, err = qbs.GetQbs()
	return q, err
}

func SetMg() (mg *qbs.Migration, err error) {
	RegisterDb()
	mg, err = qbs.GetMigration()
	return mg, err
}

func CreateDb() bool {
	q, err := ConnDb()
	defer q.Close()
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		mg, _ := SetMg()
		defer mg.Close()

		mg.CreateTableIfNotExists(new(Kvs))
		mg.CreateTableIfNotExists(new(File))

		return true
	}

	return false

}

func AddFile(ctype int64, location string, url string) error {
	q, _ := ConnDb()
	defer q.Close()
	_, err := q.Save(&File{Ctype: ctype, Location: location, Url: url})
	return err
}

func DelFile(id int64) error {
	q, _ := ConnDb()
	defer q.Close()
	f := GetFile(id)

	if utils.Exist("." + f.Location) {
		if err := os.Remove("." + f.Location); err != nil {
			return err
			fmt.Println(err)
		}
	}

	//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
	_, err := q.Delete(&f)
	fmt.Println(err)
	return err
}

func GetFile(id int64) (f File) {
	q, _ := ConnDb()
	defer q.Close()
	q.Where("id=?", id).Find(&f)
	return f
}

func GetAllFile() (f []*File) {
	q, _ := ConnDb()
	defer q.Close()
	q.OrderByDesc("id").FindAll(&f)
	return f
}

func GetAllFileByCtype(ctype int64) (f []*File) {
	q, _ := ConnDb()
	defer q.Close()
	q.WhereEqual("ctype", ctype).OrderByDesc("id").FindAll(&f)
	return f
}

func SaveFile(f File) error {
	q, _ := ConnDb()
	defer q.Close()
	_, e := q.Save(&f)
	return e
}

func SetFile(id int64, pid int64, ctype int64, filename string, content string, hash string, location string, url string, size int64) error {
	q, _ := ConnDb()
	defer q.Close()
	var f File
	if q.WhereEqual("id", id).Find(&f); f.Id == 0 {
		_, err := q.Save(&File{Id: id, Pid: pid, Ctype: ctype, Filename: filename, Content: content, Hash: hash, Location: location, Url: url, Size: size})
		return err
	} else {
		type File struct {
			Pid      int64
			Ctype    int64
			Filename string
			Content  string
			Hash     string
			Location string
			Url      string
			Size     int64
		}
		_, err := q.WhereEqual("id", id).Update(&File{Pid: pid, Ctype: ctype, Filename: filename, Content: content, Hash: hash, Location: location, Url: url, Size: size})

		return err
	}
	return nil
}

func AddKV(k string, v string) error {
	q, _ := ConnDb()
	defer q.Close()
	_, err := q.Save(&Kvs{K: k, V: v})
	return err
}

func SetKV(k string, v string) error {
	q, _ := ConnDb()
	defer q.Close()
	var kvs Kvs
	if q.Where("k=?", k).Find(&kvs); kvs.Id == 0 {
		_, err := q.Save(&Kvs{K: k, V: v})
		return err
	} else {
		type Kvs struct {
			K string
			V string
		}

		_, err := q.WhereEqual("k", k).Update(&Kvs{K: k, V: v})

		return err
	}
	return nil
}

func GetKV(k string) (v string) {
	q, _ := ConnDb()
	defer q.Close()
	var kvs Kvs
	q.Where("k=?", k).Find(&kvs)
	return kvs.V
}
