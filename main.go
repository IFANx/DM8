//package main
//
////引用包
//import (
//	"fmt"
//	_ "github.com/alexbrainman/odbc" // go官方驱动
//	"github.com/go-xorm/xorm"
//)
//
////定义引擎
//var (
//	engine *xorm.Engine
//)
//
//func main() {
//	//数据库连接参数
//	params := fmt.Sprintf("driver={DM8 ODBC DRIVER};server=110.40.129.180:5236;database=DAMENG;uid=SYSDBA;pwd=SYSDBA;charset=utf8")
//
//	/*driver 对应 DM8-ODBC 数据源名称；
//	  erver =数据库 IP
//	  database = 数据库名称
//	  uid = 数据库用户
//	  pwd = 数据库用户密码
//	  charset = 字符集
//	*/
//
//	//创建 Engine 实例连接数据库
//	Engine, err := xorm.NewEngine("odbc", params)
//	if err != nil {
//		fmt.Println("new engine got error:", err)
//		return
//	}
//
//	//通过 Engine.Ping() 验证数据库是否连接成功
//	if err := Engine.Ping(); err != nil {
//		fmt.Println("ping got error:", err)
//		return
//	}
//
//	Engine.ShowSQL(true)
//	//var sql1 = `insert into DMHR.XUKANG VALUES('CD','成都',1);`
//	//rows, err := Engine.Query(sql1)
//	var sql1 = `select * FROM TEST."city";`
//	rows, err := Engine.Query(sql1)
//	//var sql1 = `delete from DMHR.XUKANG where DMHR.XUKANG.c3=8;`
//	//rows, err := Engine.Exec(sql1)
//	if err != nil {
//		return
//	}
//	fmt.Println("queryTable results:")
//	fmt.Println(rows)
//
//	return
//}
//
////-------------------------------------------------------------
///*该例程实现插入数据，修改数据，删除数据，数据查询等基本操作。*/
/*该例程实现插入数据，修改数据，删除数据，数据查询等基本操作。*/
package main

// 引入相关包
import (
	"database/sql"
	"fmt"
	"gitee.com/chunanyong/dm"
	"io/ioutil"
	"time"
)

var db *sql.DB
var err error

func main() {
	driverName := "dm"
	dataSourceName := "dm://SYSDBA:SYSDBA@localhost:5236"

	if db, err = connect(driverName, dataSourceName); err != nil {
		fmt.Println(err)
		return
	}
	if err = insertTable(); err != nil {
		fmt.Println(err)
		return
	}
	if err = updateTable(); err != nil {
		fmt.Println(err)
		return
	}
	if err = queryTable(); err != nil {
		fmt.Println(err)
		return
	}
	if err = deleteTable(); err != nil {
		fmt.Println(err)
		return
	}
	if err = disconnect(); err != nil {
		fmt.Println(err)
		return
	}
}

/* 创建数据库连接 */
func connect(driverName string, dataSourceName string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	if db, err = sql.Open(driverName, dataSourceName); err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Printf("connect to \"%s\" succeed.\n", dataSourceName)
	return db, nil
}

/* 往产品信息表插入数据 */
func insertTable() error {
	var inFileName = "C:\\Users\\Administrator\\Pictures\\xukang.jpg"
	var sql = `INSERT INTO production.product(name,author,publisher,publishtime, 
				product_subcategoryid,productno,satetystocklevel,originalprice,nowprice,discount, 
                description,photo,type,papertotal,wordtotal,sellstarttime,sellendtime) 
                VALUES(:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17);`
	data, err := ioutil.ReadFile(inFileName)
	if err != nil {
		return err
	}
	t1, _ := time.Parse("2006-Jan-02", "2005-Apr-01")
	t2, _ := time.Parse("2006-Jan-02", "2006-Mar-20")
	t3, _ := time.Parse("2006-Jan-02", "1900-Jan-01")
	_, err = db.Exec(sql, "三国演义", "罗贯中", "中华书局", t1, 4, "9787101046121", 10, 19.0000, 15.2000, 8.0,
		"《三国演义》是中国第一部长篇章回体小说，中国小说由短篇发展至长篇的原因与说书有关。",
		data, "25", 943, 93000, t2, t3)
	if err != nil {
		return err
	}
	fmt.Println("insertTable succeed")
	return nil
}

/* 修改产品信息表数据 */
func updateTable() error {
	var sql = "UPDATE production.product SET name = :name WHERE productid = 11;"
	if _, err := db.Exec(sql, "三国演义（上）"); err != nil {
		return err
	}
	fmt.Println("updateTable succeed")
	return nil
}

/* 查询产品信息表 */
func queryTable() error {
	var productid int
	var name string
	var author string
	var description dm.DmClob
	var photo dm.DmBlob
	var sql = "SELECT productid,name,author,description,photo FROM production.product WHERE productid<11"
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("queryTable results:")
	for rows.Next() {
		if err = rows.Scan(&productid, &name, &author, &description, &photo); err != nil {
			return err
		}
		blobLen, _ := photo.GetLength()
		fmt.Printf("%v	%v	%v	%v	%v\n", productid, name, author, description, blobLen)
	}
	return nil
}

/* 删除产品信息表数据 */
func deleteTable() error {
	var sql = "DELETE FROM production.product WHERE productid = 11;"
	if _, err := db.Exec(sql); err != nil {
		return err
	}
	fmt.Println("deleteTable succeed")
	return nil
}

/* 关闭数据库连接 */
func disconnect() error {
	if err := db.Close(); err != nil {
		fmt.Printf("db close failed: %s.\n", err)
		return err
	}
	fmt.Println("disconnect succeed")
	return nil
}
