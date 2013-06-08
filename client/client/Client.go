package client

import (
	"../../server/server"
	"../libs"
	//"../models"
	"../utils"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	outtimes = "Error:"
)

type UploaderHandler struct {
	libs.BaseHandler
}

func (self *UploaderHandler) Get() {
	/*
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.WriteString(outtimes + "SDFS连接超时！")
		} else {
			self.TplNames = "index.html"
			self.Render()
		}
	*/

	self.TplNames = "index.html"
	self.Render()

}

func (self *UploaderHandler) Post() {
	/*
		//TODO: Validate the file type

	*/
	/*
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			_, handler, _ := self.GetFile("uploadfile")

			if handler != nil {
				self.Ctx.WriteString(outtimes + "上传“ " + handler.Filename + " ”失败，请你重新登录，现已超时操作！")
			} else {
				self.Ctx.WriteString(outtimes + "请你重新登录，现已超时操作！")
			}
		} else {

		}
	*/

	//设置临时目录路径
	targetFolder := "/static/temp/"
	//客戶端接收文件 判斷是否接收文件成功
	if file, handler, e := self.GetFile("uploadfile"); e != nil {
		self.Data["MsgErr"] = e
	} else {
		//如果接收文件成功
		if handler != nil {
			//文件後綴
			ext := "." + strings.Split(handler.Filename, ".")[1]
			//產生臨時文件名
			filename := utils.MD5(time.Now().String()) + ext
			//生成路徑目錄
			os.MkdirAll("."+targetFolder, 0644)
			path := targetFolder + filename

			if f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
				self.Data["MsgErr"] = err
			} else {

				//把內存中的file通過io copy到臨時文件f
				io.Copy(f, file)

				//縮略圖處理
				input_file := "." + path
				output_file := input_file //跟input_file路徑一樣是爲了把文件複寫到同一個文件路徑
				output_size := "950"
				output_align := "center"
				background := "black"
				utils.Thumbnail(input_file, output_file, output_size, output_align, background)

				//讀取縮略處理后的文件之filehash值
				hash, err := utils.Filehash(output_file)
				if err != nil {
					fmt.Println("filehash error:", err)
				}

				//正態隨機產生目標節點
				serverid_tmp, nodename := server.RandNode(-1)
				serverid := strconv.Itoa(int(serverid_tmp))

				//構造發送命令url
				actionurl := "http://" + nodename + "/setfile/" + hash + "_" + serverid + ext

				//todo:檢查目標節點是否可用以及要保存的文件是否已經存在

				//往目標節點發送文件
				if resp, err := utils.PostFile(output_file, actionurl, "file"); err != nil {

					var fsize int64 = 0
					if fileInfo, err := os.Stat(output_file); err == nil {
						fsize = fileInfo.Size() / 1024
					}

					fmt.Println(resp, err)
					fmt.Println(">>>文件" + path + "保存错误，filehash：" + string(hash) + "filesize:" + string(fsize))
				} else {

					if body, err := ioutil.ReadAll(resp.Body); err != nil {
						fmt.Println("resp.Body err", err)
					} else {
						fmt.Println("body", string(body))
						self.Data["MsgErr"] = "<img src=\"" + string(body) + "\" alt=\"" + string(hash) + "\" />"

						self.Ctx.WriteString(self.Data["MsgErr"].(string))
						//models.SetFile(0, pid, 0, handler.Filename, "", string(hash), string(body) , "", fsize)
					}

				}

			}

		} else {
			self.Data["MsgErr"] = "error:文件没上传！"
		}
	}

}
