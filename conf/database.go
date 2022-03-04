package conf

import (
	"database/sql"

	"github.com/astaxie/beego"
	_ "github.com/denisenkom/go-mssqldb"
)

var Db *sql.DB
var Db2 *sql.DB

func init() {
	Db, _ = sql.Open("mssql", "sqlserver://"+beego.AppConfig.String("mssqluser")+":"+beego.AppConfig.String("mssqlpass")+"@"+beego.AppConfig.String("mssqlurls")+"?database="+beego.AppConfig.String("mssqldb")+"&connection+timeout=0")
	Db2, _ = sql.Open("mssql", "sqlserver://"+beego.AppConfig.String("mssqluser_Detail")+":"+beego.AppConfig.String("mssqlpass_Detail")+"@"+beego.AppConfig.String("mssqlurls_Detail")+"?database="+beego.AppConfig.String("mssqldb_Detail")+"&connection+timeout=0")
}
