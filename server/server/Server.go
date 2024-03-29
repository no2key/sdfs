package server

import (
	"../../utils"
	"../libs"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	filehash, ext, local_serverid, remote_serverid, nodename string

	publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`)

	privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)
)

type RStorageHandler struct {
	libs.BaseHandler
}

type WStorageHandler struct {
	libs.BaseHandler
}

func LocalNode() (int64, string) {
	/*
		sdfs每台机器设置一个 数字id 和对应的 二级域名
		node name  应该使用域名，不要使用ip，譬如file1.veryhour.com file2.veryhour.com file3.veryhour.com

	*/
	// number id & node name (sub domain)
	return 0, "file0.veryhour.com:8800"
}

func RandNode(m int64) (int64, string) {
	nodeskv := make([]string, 3, 10)
	nodeskv[0] = "file0.veryhour.com:8800"
	nodeskv[1] = "file1.veryhour.com:8801"
	nodeskv[2] = "file2.veryhour.com:8802"
	if m < 0 {
		//根据正态分布随机整数选择节点服务器
		n := int64(utils.Nrand(int64(len(nodeskv))))
		return n, nodeskv[n]
	} else {
		//根据m值查找对应编号的节点名称
		return m, nodeskv[m]
	}
}

func (self *RStorageHandler) Get() {
	fmt.Println("读取模式")
	//根据URL去指定服务器读文件

	//URL HASH组成：
	//http://file0.veryhour.com/getfile/FileHash_ServerId.jpg
	// http://file0.veryhour.com/getfile/0fe7bfba7443fa894f97e544085ca6c7_0.jpg
	/*
	   get状态：
	   根据请求的来路URL，分拆URL hash字符串，读出文件hash值
	   以文件hash值为对象，找出应该去哪个节点读取文件
	   判断是否为本节点，如果是则读取本机上的文件并返回给用户，
	   如果不是则跳转到指定节点去处理。
	*/

	filename := self.GetString(":filename")
	for i, v := range strings.Split(filename, "_") {

		if i == 0 {
			filehash = v
		}
		if i == 1 {
			for ii, vv := range strings.Split(v, ".") {

				if ii == 0 {
					local_serverid = vv
				}
				if ii == 1 {
					ext = vv
				}
			}
		}
	}

	nn, _ := strconv.Atoi(local_serverid)
	_, nodename := RandNode(int64(nn))

	fileurl := "http://" + nodename + "/getfile/" + filehash + "_" + local_serverid + "." + ext

	//判断是否为本节点，如果是则接着处理，如果不是则跳转到指定节点。
	if nid, _ := LocalNode(); nid == int64(nn) {

		filedir := "./data/" + utils.FixedpathByString(filehash, 3)

		filename = filehash + "." + ext
		filepath := filedir + filename

		if file, err := os.Open(filepath); err != nil {
			fmt.Println("找不到指定文件", err)
		} else {

			if _, e := io.Copy(self.Ctx.ResponseWriter, file); e != nil {
				fmt.Println(e)
			} else {
				fmt.Println("Server fileurl:", fileurl)
			}
		}

	} else {
		fmt.Println("Redirect fileurl:", fileurl)
		self.Redirect(fileurl, 302)
	}

}

func (self *WStorageHandler) Post() {
	fmt.Println("写入模式")
	/*
		获取到文件后，我们用正态分布随机算法选择出存储节点，
		然后判断该节点是不是本机，如果是则直接把文件保存到本节点上即可，
		如果不是当前节点，则跳到目标节点去保存。
	*/
	//随机写文件，如果遇到节点服务器无法写入（已满）则再次随机，
	//直到随机次数大于节点数依然无法写入的时候则终止写入，并返回写入失败的状态（表示整个存储网络都不可用）

	//URL HASH组成：
	//http://localhost/setfile/FileHash_senderServerId.png

	/*
		post状态：
		接受文件
		计算文件hash值
		以ServerId+FileHash+文件后缀作为文件名，
		以正态随机算法算出一个随机数值R，然后判断R是否为当前节点，如果是则保存到当前机器。如果不是则跳转到节点R
		如果为当前节点，则在当前机器存储文件，
	*/

	filename := self.GetString(":filename")
	for i, v := range strings.Split(filename, "_") {

		if i == 0 {
			filehash = v
		}
		if i == 1 {
			for ii, vv := range strings.Split(v, ".") {

				if ii == 0 {
					remote_serverid = vv
				}
				if ii == 1 {
					ext = vv
				}
			}
		}
	}
	nn, _ := strconv.Atoi(remote_serverid)
	_, nodename := RandNode(int64(nn))

	//判断是否为本节点，如果是则接着处理，如果不是则跳转到指定节点。
	if nid, _ := LocalNode(); nid == int64(nn) {
		//接收POST过来的文件
		if file, handler, e := self.GetFile("file"); handler == nil || e != nil {
			fmt.Println("GET THE POST FILE %s ERROR!", self.Ctx.Request.RequestURI)
		} else {

			tmpfiledir := "./static/temp/"
			tmpfilepath := tmpfiledir + filename
			filedir, filepath := "", ""
			//打开 临时文件的句柄tf 以便保存 io.Copy 过来的 文件句柄file
			if tf, err := os.OpenFile(tmpfilepath, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
				fmt.Println("os.OpenFile tf error:", err)
			} else {
				//利用io.Copy转换文件句柄类型
				if _, e := io.Copy(tf, file); e != nil {
					fmt.Println("io.Copy tf error:", e)
				} else {

					filehash, _ = utils.Filehash(tmpfilepath)
					filedir = "./data/" + utils.FixedpathByString(filehash, 3)

					filename = filehash + "." + ext
					filepath = filedir + filename
					//预设目录
					os.MkdirAll(filedir, 0644)
					//保存文件
					if e := utils.MoveFile(tmpfilepath, filepath); e != nil {
						fmt.Println("SaveFile error:", e)
						self.Ctx.WriteString("error")
					} else {
						fmt.Println("SaveFile Okay!")
						fileurl := "http://" + nodename + "/getfile/" + filehash + "_" + remote_serverid + "." + ext

						if data, err := utils.RsaEncrypt([]byte(fileurl), publicKey); err != nil {
							fmt.Println(err)
						} else {
							if origData, err := utils.RsaDecrypt(data, privateKey); err != nil {
								fmt.Println(err)
							} else {
								fmt.Println("origData:", string(origData))
							}
						}
						//self.Redirect(fileurl, 302)
						self.Ctx.WriteString(fileurl)
					}
				}
			}
		}

	} else {
		fileurl := "http://" + nodename + "/setfile/" + filehash + "_" + remote_serverid + "." + ext
		fmt.Println("Redirect SaveFile fileurl:", fileurl)
		self.Redirect(fileurl, 307)
	}

	//判断是否为本节点，如果是则接着处理，如果不是则跳转到指定节点。
	/*
		if nid, _ := LocalNode(); nid == int64(nn) {
			//接收POST过来的文件
			if file, handler, e := self.GetFile("file"); handler == nil || e != nil {
				fmt.Println("SetFile fileurl %s error!", fileurl)
			} else {
				//预设目录
				os.MkdirAll(filedir, 0644)

				//打开 临时文件的句柄tf 以便保存 io.Copy 过来的 文件句柄file
				if tf, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
					fmt.Println("os.OpenFile tf error:", err)
				} else {
					//利用io.Copy转换文件句柄类型
					if _, e := io.Copy(f, file); e != nil {
						fmt.Println("io.Copy error:", e)
					} else {
						//保存文件
						if e := utils.SaveFile(filedir, filename, f); e != nil {
							fmt.Println("SaveFile error:", e)
							self.Ctx.WriteString("error")
						} else {
							fmt.Println("SaveFile Okay!")
							fileurl = "http://" + nodename + "/getfile/" + filehash + "_" + remote_serverid + "." + ext
							//self.Redirect(fileurl, 302)
							self.Ctx.WriteString(fileurl)
						}
					}
				}
			}
		} else {
			fmt.Println("Redirect SaveFile fileurl:", fileurl)
			self.Redirect(fileurl, 307)
		}
	*/

}
