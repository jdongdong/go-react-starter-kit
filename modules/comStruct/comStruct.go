package comStruct

import (
	"fmt"
	"net/http"

	"github.com/go-xorm/xorm"
	"github.com/jdongdong/go-lib/slog"
	"github.com/jdongdong/go-react-starter-kit/modules/apiCode"
	"github.com/labstack/echo"
	"github.com/pquerna/ffjson/ffjson"
)

type ResModel struct {
	R       interface{} `json:"r"`
	Success string      `json:"success"`
	Code    int         `json:"code"`
	Info    string      `json:"info"`
}

type PageResModel struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

type CusHandlerFunc func(*CustomContext) error

type CustomContext struct {
	UserID string
	echo.Context
}

func (this *CustomContext) toResModel(err_code int, info string, r interface{}) ResModel {
	if r == nil {
		r = ""
	}
	success := "T"
	if err_code != apiCode.Success {
		success = "F"
	}

	res := ResModel{Success: success, Code: err_code, Info: info, R: r}
	defer func() {
		vlu, err := ffjson.Marshal(r)
		if err == nil {
			slog.Debug(fmt.Sprintf("url:%s success:%s code:%d info:%s r:%s", this.Request().RequestURI, success, err_code, info, string(vlu)))
		} else {
			slog.Debug(fmt.Sprintf("url:%s success:%s code:%d info:%s r:%+v", this.Request().RequestURI, success, err_code, info, r))
		}
	}()
	return res
}

func (this *CustomContext) Success(r interface{}) error {
	vlu, err := ffjson.Marshal(this.toResModel(apiCode.Success, "", r))
	if err != nil {
		return err
	}
	return this.String(http.StatusOK, string(vlu))
}

func (this *CustomContext) Fail(err error) error {
	vlu, err := ffjson.Marshal(this.toResModel(apiCode.FormatApiCode(err), err.Error(), nil))
	if err != nil {
		return err
	}
	return this.String(http.StatusOK, string(vlu))
}

func (this *CustomContext) AutoPageDataRs(i interface{}, total int64, err error) error {
	if err != nil {
		return err
	}
	return this.Success(PageResModel{Data: i, Total: total})
}

func (this *CustomContext) AutoDataRs(i interface{}, err error) error {
	if err != nil {
		return err
	}
	return this.Success(i)
}

//func (this *CustomContext) ToJson(i interface{}) error {
//	//body, err := ioutil.ReadAll(this.Request().Body)
//	//if err != nil {
//	//	return errCode.CheckErrorInvalidJson(err)
//	//}
//	//
//	//return errCode.CheckErrorInvalidJson(ffjson.Unmarshal(body, &i))
//	return errCode.CheckErrorInvalidJson(json.NewDecoder(this.Request().Body).Decode(i))
//}

type SeaModel struct {
	PageIndex    int `json:"pageIndex"`
	PageSize     int `json:"pageSize"`
	seaInterface SeaInterface
}

type SeaDtlModel struct {
	seaDtlInterface SeaDtlInterface
}

type TreeModel struct {
	Key      string      `json:"key"`
	Title    string      `json:"title"`
	Type     string      `json:"type"`
	Checked  bool        `json:"checked"`
	Children []TreeModel `json:"children"`
}

type LeftMenuModel struct {
	Id     int64  `json:"id"`
	Pid    int64  `json:"pid"`
	MPid   int64  `json:"breadPid"`
	Sort   int    `json:"sort"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Router string `json:"router"`
}

type SeaInterface interface {
	where(session *xorm.Session)
}

type SeaDtlInterface interface {
	whereDtl(session *xorm.Session)
}

type PagingInterface interface {
	GetPaging() (interface{}, int64, error)
}

type PagingDtlInterface interface {
	GetDtlPaging() (interface{}, int64, error)
}

type InsertInterface interface {
	Insert() error
}

type UpdateByIdInterface interface {
	UpdateById() error
}

type DeleteByIdInterface interface {
	DeleteById() error
}

type InsertTransInterface interface {
	InsertTrans() error
}

type UpdateTransInterface interface {
	UpdateTrans() error
}

type DeleteTransInterface interface {
	DeleteTrans() error
}
