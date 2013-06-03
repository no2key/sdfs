package libs

import (
	"../../utils"
	"github.com/insionng/torgo"
	//"github.com/astaxie/beego"
	//"../torgo"
)

var (
	bc *torgo.BeeCache
)

type BaseHandler struct {
	torgo.Handler
	//beego.Controller
}

func init() {
	bc = torgo.NewBeeCache()
	bc.Every = 259200 //該單位為秒，0為不過期，259200 三天,604800 即一個星期清空一次緩存
	bc.Start()
}

//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
func (self *BaseHandler) Prepare() {

}

func (self *BaseHandler) Render() (err error) {

	var ivalue []byte
	ck, _ := self.Ctx.Request.Cookie("lang")
	lang := ""

	if ck != nil {
		lang = ck.Value
	} else {
		lang = "normal"
	}

	if self.GetString("lang") != "" {

		if self.GetString("lang") == "normal" {
			lang = "normal"
		}

		if self.GetString("lang") == "cn" {
			lang = "zh-cn"
		}

		if self.GetString("lang") == "hk" {
			lang = "zh-hk"
		}

	}

	self.Ctx.SetCookie("lang", lang, "", "", 0)
	self.Data["lang"] = lang

	rb, e := self.RenderBytes()
	rs := string(rb)
	ikey := utils.MD5(rs + lang)
	if bc.IsExist(ikey) {
		ivalue = bc.Get(ikey).([]byte)
	} else {

		if lang == "normal" {
			ivalue = rb
		} else {
			ivalue = utils.Convzh(rs, lang)
		}

		bc.Put(ikey, ivalue, 259200)

	}

	return self.RenderCore(ivalue, e)

}
