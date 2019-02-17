自动完成beego权限认证以及快捷开发api的工具

#简单快速完成你的高性能api#

## install ##

    go get github.com/zm4728/beescaffold



use

# 最简单的入门例子 #


1.创建gofile

    import "github.com/zm4728/beescaffold"
    
    type YouController struct {
    	beescaffold.BeeController
    }
    
    func (this *ImportExcelController)Get(){
        this.R=map[string]string{
            "result":"ok",
        }
    }
    
2.刷新你的chrome

    {
      "error": 0,
      "message": "",
      "data": {
          "result": "ok"
      }
    }
    
看到吗一切都是自动的,

    this.R is interface{}
    
    this.R="o" 
    
这种写法也是正确的

简单分析
  继承 beescaffold.BeeController
  然后设置一个this.R的值他会自动处理json，包括错误



# 现在学习下怎么处理api错误结果 #

1.创建gofile 

    import "github.com/zm4728/beescaffold"
    
    type YouController struct{
    beescaffold.BeeController
    }
    
    func (this *YouController)Get(){
       	this.E=beescaffold.Make_E_defaultcode_with_s("this is error")
    }

2.刷新你的chrome
    
    {
      "error": -1,
      "message": "this is error",
      "data": null
    }

一样是自动的，你只需要设置一个this.E

注意:this.R或者this.E基本任何api都是这种模式


# **设置自定义错误代码** #

替换 beescaffold.Make_E_defaultcode_with_s 为 Make_E_with_s_code

    func (this *YouController)Get(){
       	this.E=beescaffold.Make_E_with_s_code("this is error",500)
    }

这时候会返回自定义错误代码    


# 添加你的自定义代码，如果没有好的位置，通常在controller的init函数操作 #

    var E_NotLogin=beescaffold.Make_E_with_s_code("not login",-2)
    
    func (this *YouController)Get(){
       	this.E=E_NotLogin
    }
    


设置 默认的错误代码 

    
    func main(){
    beescaffold.DefaultCode=-7
    }
    
    


预定义保留错误代码：

    //保留session
    var Code_Error_Session=-2
    var E_Error_Session=Make_E_with_s_code("session效验失败",Code_Error_Session)
    var Code_Error_Auth=-3
    var E_Error_Auth=Make_E_with_s_code("用户权限不够",Code_Error_Auth)

我建议你试用-100以下的代码，因为功能还会拓展


# 自动认证session和权限功能 #



### Create_Mysql_with_host_port_u_p_db(_host string,_port int,_u string,_p string,_db string)(*xorm.Engine,error) ###
创建xorm的mysql连接，并且设置线程池最大闲置50，最大并发500,并且转换所有表和字段为小写格式的api




### Register_tokenmodel_with_xorm ###
注册auth的表到数据库中


#### beescaffold.Register_router_Allow("*") ####

对路径加入AllowOrigins，以解决跨域问题
参数是正则表达式可以随意匹配path


# 为beego添加自动认证能力 #

## 1.继承beescaffold的BeeController ##

    type PackageController struct {
	    beescaffold.BeeController
    }

## 2.注册路由到beego ##

    func init() {
	    beescaffold.Register_router_Allow("*")
	    beego.Router("/app",&controllers.PackageController{beescaffold.BeeController{Role:2}})
    }

#### 好了，现在你的beego已经具备了权限认证，并且是多级权限认证，下边验证一下

#### 打开数据库的 token 表 添加一条token记录 uid:1 session:1 expire:当前时间+1天(注意这里要写时间戳unix time) role写1 ####

#### 现在打开http://localhost:8080/app 这时候你会发现权限被阻拦了，因为role的权限小于路由权限 ####

#### 修改role为2，这时候权限匹配，允许访问 ####


####注意:不要使用权限小于2的role，因为他是被保留的####



## 实例 自动认证 ##


    package models
    
    import (
    	"github.com/go-xorm/xorm"
    	"fmt"
    	"beescaffold"
    )
    
    var db * xorm.Engine
    func Run(){
    var err error
	db,err=beescaffold.Create_Mysql_with_host_port_u_p_db("127.0.0.1",3306,"root","Aa1231231","N")
    	if err != nil {
    		panic(err)
    	}
    	beescaffold.Register_tokenmodel_with_xorm(db)    
    	err=db.Sync(Shenbo{},)
    	if err!=nil{
    		fmt.Println(err)
    	}
    	fmt.Println("已经成功构造数据库")
    }




    package controllers

    import (
	        "github.com/zm4728/beescaffold"
		    "papi/models"
	    )

    type PackageController struct {
	    beescaffold.BeeController
    }

    func (this* PackageController) Get(){
	    alias:=this.GetString(":alias")
	    p,err:=models.Get_package_with_alias(alias)
	    if err!=nil{
		    this.E=beescaffold.Make_E_defaultcode_with_err(err)
		    return
	    }
	    this.R=p
    }



    package routers

    import (
	    "github.com/astaxie/beego"
	    "papi/controllers"
	    "github.com/zm4728/beescaffold"
    )


    func init() {
	    beescaffold.Register_router_Allow("*")
	    beego.Router("/app/:alias",&controllers.PackageController     {beescaffold.BeeController{Role:2}})
	    beego.Router("/app1",&controllers.PackageController{})
    }


这个项目有对应的vue客户端