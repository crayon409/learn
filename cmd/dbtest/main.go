package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var conf string
var h bool

func init() {
	flag.StringVar(&conf, "conf", "./db.yml", "test conf file")
	flag.BoolVar(&h, "h", false, "")
}

type DbNode struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}
type DB struct {
	List []DbNode `yaml:"db"`
}

var (
	db     DB
	data   []byte
	err    error
	client *gorm.DB
)

func main() {
	flag.Parse()
	if h {
		help()
		return
	}
	fmt.Printf("conf:%s\n", conf)
	fd, _ := os.Open(conf)
	data, _ = io.ReadAll(fd)
	err = yaml.Unmarshal(data, &db)
	if err != nil {
		panic(err)
	}

	for _, node := range db.List {
		dsn := parse(&node)
		if client, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		}); err != nil {
			log.Fatal(err)
		}
		sql := "select 1"
		err = client.Exec(sql).Error
		if err != nil {
			log.Fatal(fmt.Errorf("%w", err))
		}
	}
	fmt.Errorf("")
}

func parse(node *DbNode) string {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true&loc=Local",
		node.User,
		node.Pass,
		node.Host,
		node.Port,
		node.Name)
}

func help() {
	fmt.Printf("help start\n")
	tmp := `db:
  - host: ''
    name: ''
    user: ''
    pass: ''
    port: ''`
	fd, err := os.OpenFile(conf, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer fd.Close()
	if err != nil {
		fmt.Printf("open err: %v\n", err)
		return
	}
	_, err = fd.WriteString(tmp)
	if err != nil {
		fmt.Printf("write err: %v\n", err)
		return
	}
	err = fd.Sync()
	if err != nil {
		return
	}

}
