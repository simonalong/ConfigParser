package test

import (
	"context"
	"fmt"
	"github.com/simonalong/gole/config"
	orm2 "github.com/simonalong/gole/extend/orm"
	"testing"
	"xorm.io/xorm/contexts"
)

func TestXorm1(t *testing.T) {
	config.LoadYamlFile("./application-test1.yaml")
	orm2.AddXormHook(&GobaseXormHook{})
	db, _ := orm2.NewXormDb()

	// 删除表
	db.Exec("drop table isc_demo.gobase_demo")

	//新增表
	db.Exec("CREATE TABLE gobase_demo(\n" +
		"  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n" +
		"  `name` char(20) NOT NULL COMMENT '名字',\n" +
		"  `age` INT NOT NULL COMMENT '年龄',\n" +
		"  `address` char(20) NOT NULL COMMENT '名字',\n" +
		"  \n" +
		"  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n" +
		"  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n" +
		"\n" +
		"  PRIMARY KEY (`id`)\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")

	db.Table("gobase_demo").Insert(&GobaseDemo{Name: "zhou", Age: 18, Address: "杭州"})
	// 新增
	db.Table("gobase_demo").Insert(&GobaseDemo{Name: "zhou", Age: 18, Address: "杭州"})

	var demo GobaseDemo
	db.Table("gobase_demo").Where("name=?", "zhou").Get(&demo)

	dd := db.DB()
	dd.Query("select * from gobase_demo")

	// 查询：多行
	fmt.Println(demo)
}

type GobaseXormHook struct {
}

func (*GobaseXormHook) BeforeProcess(c *contexts.ContextHook, driverName string) (context.Context, error) {
	fmt.Println("before-xorm")
	return c.Ctx, nil
}

func (*GobaseXormHook) AfterProcess(c *contexts.ContextHook, driverName string) error {
	fmt.Println("after-xorm")
	return nil
}
