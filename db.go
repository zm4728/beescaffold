package beescaffold

import(
	_ "github.com/go-sql-driver/mysql"
	_"github.com/mattn/go-sqlite3"
	"github.com/go-xorm/xorm"
	"strings"
	"fmt"
	"errors"
	"time"
)

type Token struct {
	Id      int64  `json:"id"`                //id主键
	Uid     int64  `xorm:"unique" json:"uid"` //用户uid
	Session string `json:"session"`             //session
	Expire  int64  `json:"expire"`            //过期时间
	Role    int    `json:"role"`              //用户权限
	Created int64  `xorm:"created" json:"created"`
	Updated int64  `xorm:"updated" json:"updated"`
}

var _eng * xorm.Engine

type ISameMapper struct {
}
func (m ISameMapper) Obj2Table(o string) string {
	return strings.ToLower(o)
}
func (m ISameMapper) Table2Obj(t string) string {
	return strings.ToLower(t)
}

func Ping_mysql_with_db_date(_db* xorm.Engine,_d time.Duration){
	go func() {
		for{
			_db.Ping()
			time.Sleep(_d)
		}
	}()
}



func Create_sqlite3_with_db(_db string)(*xorm.Engine,error){
	eng,err:=xorm.NewEngine("sqlite3",_db)
	if err!=nil{
		return nil,err
	}
	//初始化所有字段和标明为小写
	eng.SetMapper(ISameMapper{})

	//线程池处理
	eng.DB().SetMaxIdleConns(50)
	eng.DB().SetMaxOpenConns(500)
	return eng,nil
}

func Create_Mysql_with_host_port_u_p_db(_host string,_port int,_u string,_p string,_db string)(*xorm.Engine,error){
	eng,err:=xorm.NewEngine("mysql",fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",_u,_p,_host,_port,_db))
	if err!=nil{
		return nil,err
	}
	//初始化所有字段和标明为小写
	eng.SetMapper(ISameMapper{})

	//线程池处理
	eng.DB().SetMaxIdleConns(50)
	eng.DB().SetMaxOpenConns(500)
	return eng,nil
}
//注册token模块
func Register_tokenmodel_with_xorm(_db* xorm.Engine)error{
	_eng=_db
	return _db.Sync(Token{})
}


func Update_token_with_uid_session_expire_role(_uid int64,_session string,_expire int64,_role int)error{
	n,err:=_eng.Table(Token{}).Where("uid=?",_uid).Update(map[string]interface{}{
		"uid":_uid,
		"session":_session,
		"expire":_expire,
		"role":_role,
	})
	if err!=nil{
		fmt.Println(err)
		return err
	}
	fmt.Println("当前没有发生错误",n)
	if n==1{
		return nil
	}
	fmt.Println(n)
	err=Insert_token_with_uid_session_expire_role(_uid,_session,_expire,_role)
	return err
}

func Insert_token_with_uid_session_expire_role(_uid int64,_session string,_expire int64,_role int)error{
	tk:=Token{
		Uid:     _uid,
		Session: _session,
		Expire:  _expire,
		Role:    _role,
	}
	n,err:=_eng.InsertOne(&tk)
	if err!=nil{
		return err
	}
	if n!=1{
		return errors.New("添加token失败")
	}
	return nil
}

func QueryToken_with_uid_session_expire(_uid int64,_session string,_expire int64)(*Token,error){
	t:=Token{}
	b,err:=_eng.Table(Token{}).Where("uid=? and session=? and expire>=?",_uid,_session,_expire).Get(&t)
	if err!=nil{
		return nil,err
	}
	if b==false{
		return nil,errors.New("session无效")
	}
	return &t,nil
}