package beescaffold

import "github.com/pkg/errors"

var DefErrors map[error]int
var DefaultCode int =-100

type R struct {
	Error int `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type E struct {
	Error error
	Code int
}


//保留session
var Code_Error_Session=-2
var E_Error_Session=Make_E_with_s_code("session效验失败",Code_Error_Session)

var Code_Error_Auth=-3
var E_Error_Auth=Make_E_with_s_code("用户权限不够",Code_Error_Auth)


func Make_Api_R_with_obj(_o interface{})*R{
	return &R{Data:_o}
}

func Make_Api_Error_with_error_code(_err error,_n int)*R{
	return &R{Error:_n,Message:_err.Error()}
}

func Make_Api_Error_with_e(_e *E)*R{
	return Make_Api_Error_with_error_code(_e.Error,_e.Code)
}

func Make_E_defaultcode_with_err(_err error)*E{
	return Make_E_with_error_code(_err,DefaultCode)
}

func Make_E_defaultcode_with_s(_s string)*E{
	return Make_E_with_s_code(_s,DefaultCode)
}

func Make_E_with_s_code(_s string,_code int)*E{
	return Make_E_with_error_code(errors.New(_s),_code)
}



func Make_E_with_error_code(_err error,_code int)*E{
	return &E{Error:_err,Code:_code}
}