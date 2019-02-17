package beescaffold

import (
	"github.com/astaxie/beego"
	"errors"
	"time"
	)

type BeeController struct {
	beego.Controller
	Role int //用户权限
	R interface{}  //返回的obj
	E *E         //返回的error
}

func (this *BeeController)Prepare(){
	if this.Role>=2{ //如果需要权限认证
		var err *E
	    defer func() {
	    	if err!=nil{
	    		this.Data["json"]=Make_Api_Error_with_e(err)
	    		this.ServeJSON(false)
			}
		}()

		uid,err1:=this.GetInt64("uid")
		if err1!=nil{
			err=Make_E_with_error_code(errors.New("效验uid错误"),Code_Error_Auth)
			return
		}
		if uid<1{
			err=Make_E_with_error_code(errors.New("效验uid错误"),Code_Error_Auth)
			return
		}
		session:=this.GetString("session")
		if len(session)<1{
			err=E_Error_Session
			return
		}
		t,err1:=QueryToken_with_uid_session_expire(uid,session,time.Now().Unix())
		if err1!=nil{
			err=E_Error_Session
			return
		}
		if t.Role<this.Role{
			err=E_Error_Auth
			return
		}
		//beego.Debug("查询到的token:",t)

	}

}

func (this *BeeController)Finish(){
	defer func() {
		var r *R
		if this.E!=nil{ //如果被标记错误
			r=Make_Api_Error_with_error_code(this.E.Error,this.E.Code)
		}else{
			if this.R==nil{  //如果未实现，也就是未返回
				//r=Make_Api_Error_with_e(Make_E_defaultcode_with_s("接口未实现"))
				r=Make_Api_R_with_obj(nil)
			}else{   //处理R返回，正常结果
				r=Make_Api_R_with_obj(this.R)
			}
		}
		this.Data["json"]=r
		this.ServeJSON(false)
	}()
}