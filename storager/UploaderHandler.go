package storager

import (
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
			targetFolder := "/archives/upload/"

			file, handler, e := self.GetFile("uploadfile")
			pid, _ := self.GetInt("pid")

			if e != nil {
				self.Data["MsgErr"] = "0"
			} else {

				if handler != nil {
					ext := "." + strings.Split(handler.Filename, ".")[1]
					filename := utils.MD5(time.Now().String()) + ext

					path := targetFolder + time.Now().Format("2006/01/02/")

					os.MkdirAll("."+path, 0644)
					path = path + filename
					f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644)
					defer f.Close()

					if err != nil {
						self.Data["MsgErr"] = "0"
					} else {
						io.Copy(f, file)
						defer file.Close()
						input_file := "." + path
						output_file := "." + path
						output_size := "534"
						output_align := "center"
						background := "black"
						utils.Thumbnail(input_file, output_file, output_size, output_align, background)
						hash := utils.Filehash(output_file)
						fileInfo, err := os.Stat(output_file)
						var fsize int64 = 0
						if err == nil {
							fsize = fileInfo.Size() / 1024
						}

						self.Data["MsgErr"] = "<img src=\"" + path + "\" alt=\"" + hash + "\" />"

						models.SetFile(0, pid, 0, handler.Filename, "", hash, path, "", fsize)
					}

				} else {
					self.Data["MsgErr"] = "0"
				}
			}

			self.Ctx.WriteString(self.Data["MsgErr"].(string))
		}
	*/

	targetFolder := "/static/temp/"

	file, handler, e := self.GetFile("uploadfile")

	if e != nil {
		self.Data["MsgErr"] = e
	} else {

		if handler != nil {

			ext := "." + strings.Split(handler.Filename, ".")[1]
			filename := utils.MD5(time.Now().String()) + ext

			path := targetFolder + time.Now().Format("2006/01/02/15/04/")

			os.MkdirAll("."+path, 0644)
			path = path + filename
			f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644)
			defer f.Close()

			if err != nil {
				self.Data["MsgErr"] = err
			} else {
				fmt.Println("处理上传文件")
				io.Copy(f, file)
				defer file.Close()
				input_file := "." + path
				output_file := "." + path
				output_size := "950"
				output_align := "center"
				background := "black"
				//fmt.Println("file block hash:", utils.Filehash_block(input_file, 2048))
				utils.Thumbnail(input_file, output_file, output_size, output_align, background)
				hash, err := utils.Filehash(output_file)
				if err != nil {
					fmt.Println("filehash error:", err)
				}

				serverid_tmp, nodename := RandNode(-1)
				serverid := strconv.Itoa(int(serverid_tmp))
				fileInfo, err := os.Stat(output_file)
				var fsize int64 = 0
				if err == nil {
					fsize = fileInfo.Size() / 1024
				}
				actionurl := "http://" + nodename + "/setfile/" + hash + "_" + serverid + ext

				if resp, err := utils.PostFile(output_file, actionurl, "file"); err != nil {
					fmt.Println(resp, err)
					fmt.Println("文件" + path + "保存错误，filehash：" + string(hash) + "filesize:" + string(fsize))
				} else {

					if body, err := ioutil.ReadAll(resp.Body); err != nil {
						fmt.Println("resp.Body err", err)
					} else {
						fmt.Println("body", string(body))
						self.Data["MsgErr"] = "<img src=\"" + string(body) + "\" alt=\"" + string(hash) + "\" />"

						//models.SetFile(0, pid, 0, handler.Filename, "", string(hash), path, "", fsize)
					}

				}

			}

		} else {
			self.Data["MsgErr"] = "error:文件没上传！"
		}
	}

	self.Ctx.WriteString(self.Data["MsgErr"].(string))
}
