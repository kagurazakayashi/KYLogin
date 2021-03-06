/*@Author: 0wew0
 * @Date: 2020-03-03 17:44:13
 * @LastEditTime: 2020-03-03 17:45:15
 * @LastEditors: 0wew0
 * @Description: In User Settings Edit
 * @FilePath: /zyz/user/tools/go/cleandb/cleandb.go
 */
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

//Filejson is json的结构.
type Filejson struct {
	Type   string                 `json:"type"`
	Temp   string                 `json:"temp"`
	Todir  string                 `json:"todir"`
	Toname string                 `json:"toname"`
	Info   map[string]interface{} `json:"info"`
	To     [][]interface{}        `json:"to"`
}

// var uploadtmp string
var watermarkimage string
var nickname string
var wmfont string

// var keypath string
var verbose bool

//redis中所有image的key
var arrkeys []string

// redis设置
// 声明一个全局的redisdb变量
var redisdb *redis.Client
var redisaddress string = "127.0.0.1:6379"
var redispassword string = "4XQa7P8kVN1EUgopyU1v43D7tZhYe8jc"
var redisdbid int = 10

func main() {
	flag.BoolVar(&verbose, "v", false, "是否显示输出")
	flag.StringVar(&redisaddress, "r", "127.0.0.1:6379", "redis地址（带端口）")
	flag.StringVar(&redispassword, "rp", "RfryNXkEUG3OnyDI06zK1mqWA7oQslqvc8IEgHh78BpACCaUZIN44zrlUyDIq8xL3unaZJpWd592DrJifvymOdLHCAIN0ycg1TzvatE2tJiu40kr06Aub1vfjYGIWadevBm70UDTClutBxWTjInt3fsZomDXQvYjrRktguqJeGT0RgfJA95XgTDQGqp2Eo7D33EhIU8zSQpjy3e97Bbl5yFvoqERz3wUBvcFd7K95Eas4DZpld3NV7fuk1tdh7Xa", "redis密码")
	flag.IntVar(&redisdbid, "rid", 1, "redis数据库ID")
	// flag.StringVar(&uploadtmp, "path", "/www/wwwroot/zyz/upload_tmp", "需要扫描的路径")
	flag.StringVar(&watermarkimage, "wmimage", "/www/wwwroot/zyz/img/logo.png", "水印logo位置")
	flag.StringVar(&nickname, "nick", "@择择#213", "指定删除多少行")
	flag.StringVar(&wmfont, "wmfont", "/usr/share/fonts/google-noto-cjk/NotoSerifCJK-Regular.ttc", "水印字体")
	// flag.StringVar(&keypath, "kpath", "/www/wwwroot/zyz/user/tools/go/convertfile/file/encrypt.keyinfo", "指定删除多少行")
	flag.Parse()

	fmt.Println("STRART convertimage v1.0")
	err := initRedis()
	check(0, "initRedis", err)

	allkeys()
}

func allkeys() {
	rkeys := redisdb.Keys("i_*")
	arrkeys = rkeys.Val()
	if verbose {
		fmt.Println(arrkeys)
	}
	if len(arrkeys) > 0 {
		runConvert()
	} else {
		rdel := redisdb.Del("ic")
		if verbose {
			fmt.Println("删除处理图片状态key", rdel)
		}
		redisdb.Close()
	}
}

func runConvert() {
	allkeysjson, err := json.Marshal(arrkeys)
	check(0, "创建正在处理图片状态key出错", err)

	err = redisdb.Set("ic", allkeysjson, 0).Err()
	check(0, "创建正在处理图片状态key出错", err)

	for _, flv := range arrkeys {
		fmt.Println("\n 当前key:", flv)
		val2, err := redisdb.Get(flv).Result()
		check(0, "redisdb.Get.Result", err)

		var someOne Filejson
		if err := json.Unmarshal([]byte(val2), &someOne); err == nil {
			if err == redis.Nil {
				fmt.Println(flv, " does not exists")
			} else if err != nil {
				check(0, "图片json错误", err)
			} else {
				// fmt.Println(someOne)
				for _, v := range someOne.To {
					if verbose {
						fmt.Println("正在处理图片：", someOne.Temp, someOne.Todir, someOne.Toname, v)
					}
					scaleImage(someOne.Temp, someOne.Todir, someOne.Toname, v)
				}
				filemd5(readfile(someOne.Temp), someOne.Todir, someOne.Toname)

				rdel := redisdb.Del(flv)
				if verbose {
					fmt.Println("删除处理图片key", rdel)
				}
				err := os.Remove(someOne.Temp)
				check(0, "del file", err)
			}
		} else {
			fmt.Println(err)
		}

	}
	allkeys()
}

// 初始化连接
func initRedis() (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     redisaddress,
		Password: redispassword, // no password set
		DB:       redisdbid,     // use default DB
	})

	_, err = redisdb.Ping().Result()
	check(0, "redisdb.Ping.Result", err)
	return nil
}

// func getFilesList(path string) []string {
// 	var fl []string
// 	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
// 		if f == nil {
// 			return err
// 		}
// 		if f.IsDir() {
// 			return nil
// 		}
// 		isjson := strings.Split(path, ".")
// 		isjsonlast := len(isjson) - 1
// 		if isjson[isjsonlast] == "json" {
// 			fl = append(fl, path)
// 		}
// 		return nil
// 	})
// 	check(0,"filepath.Walk", err)
// 	return fl
// }

func readfile(fp string) []byte {
	configfile, err := os.OpenFile(fp, os.O_RDONLY, 0755)
	defer configfile.Close()
	check(0, "os.OpenFile", err)
	fi, _ := configfile.Stat()
	data := make([]byte, fi.Size())
	n, err := configfile.Read(data)
	check(0, "configfile.Read", err)
	// fmt.Println(string(data[:n]))
	return data[:n]
}

// func readJSON(data []byte) *Filejson {
// 	var filejson Filejson
// 	data = []byte(os.ExpandEnv(string(data)))
// 	err := json.Unmarshal(data, &filejson)
// 	check(0,"json.Unmarshal", err)
// 	return &filejson
// }

// func mkdir(todir string) {
// 	cmd := exec.Command("mkdir", "-p", todir)
// 	stdout, err := cmd.StdoutPipe()
// 	check(0,"cmd.StdoutPipe", err)
// 	defer stdout.Close()
// 	err = cmd.Start()
// 	check(0,"cmd.Start", err)
// 	opBytes, err := ioutil.ReadAll(stdout)
// 	check(0,"ioutil.ReadAll", err)
// 	fmt.Println(string(opBytes))
// }

func filemd5(data []byte, todir string, toname string) {
	// cmd := exec.Command("mkdir", "-p", todir)
	// stdout, err := cmd.StdoutPipe()
	// check(err)
	// defer stdout.Close()
	// err = cmd.Start()
	// check(err)
	// opBytes, err := ioutil.ReadAll(stdout)
	// check(err)
	// fmt.Println(string(opBytes))

	ret := md5.Sum(data)
	MD5Str := hex.EncodeToString(ret[:])
	// fmt.Printf("\n-----%s-----", MD5Str)
	// 将保存的字符串转换为字节流
	str := []byte(MD5Str)
	// 保存到文件
	tomd5 := fmt.Sprintf("%s/%s.md5", todir, toname)
	// fmt.Printf("\n-----%s-----\n", tomd5)
	err := ioutil.WriteFile(tomd5, str, 0666)
	check(0, "ioutil.WriteFile", err)
}

func scaleImage(tempfile string, todir string, toname string, tofile []interface{}) {

	// time.Sleep(500000000)
	resolutionratio := fmt.Sprintf("%sx%s", strconv.FormatFloat(tofile[1].(float64), 'f', -1, 64), strconv.FormatFloat(tofile[2].(float64), 'f', -1, 64))
	precision := fmt.Sprintf("%s%s", strconv.FormatFloat(tofile[3].(float64), 'f', -1, 64), "%")
	// toext := fmt.Sprintf("%s", tofile[0])
	topath := fmt.Sprintf("%s%s.%s", todir, toname, tofile[0])
	// fmt.Printf("\n0:%s\n1:%s\n2:%s\n3:%s\n4:%s\n5:%s\n", tempfile, resolutionratio, precision, toname, toext, topath)

	tempfileimagepath := strings.Split(todir, ".")
	shfile := fmt.Sprintf("%s%s.%s.sh", tempfileimagepath[0], toname, tofile[0])

	// println("==tempfileimagepath:", tempfileimagepath[0])
	// println("shfile:", shfile)
	sh := ""
	if precision == "0%" {
		sh = fmt.Sprintf("cp %s %s", tempfile, topath)
	} else {
		temparr := strings.Split(tempfile, ".")
		// fmt.Println("====", temparr)
		if temparr[1] == "gif" {
			// sh = fmt.Sprintf("/usr/local/bin/convert %s -gravity southeast -geometry +0+0 -fill white -font %s -pointsize 16 -draw \"image over 5,5 24,24 '%s' %s\" -coalesce -resize 256x -deconstruct %s", tempfile, wmfont, watermarkimage, wmnickname, topath)
			fmt.Println("====", resolutionratio)
			sh = fmt.Sprintf("/usr/local/bin/convert -coalesce %s -coalesce -set dispose previous -gravity southeast -stroke '#000C' -font %s -pointsize 16 -strokewidth 4 -annotate +29+5  \"%s\" -stroke none -fill white -font %s -pointsize 16 -annotate +29+5 \"%s\" -geometry +3+5 null: %s -layers composite -resize %s -layers optimize %s", tempfile, wmfont, nickname, wmfont, nickname, watermarkimage, resolutionratio, topath)
		} else {
			// wmnickname := fmt.Sprintf("text 29,5 '%s'", nickname)
			sh = fmt.Sprintf("/usr/local/bin/convert %s -resize %s -gravity southeast -stroke '#000C' -font %s -pointsize 16 -strokewidth 4 -annotate +29+5  \"%s\" -stroke none -fill white -font %s -pointsize 16 -annotate +29+5 \"%s\" \\( %s -geometry +3+5 \\) -quality %s -composite %s >>convertimage.log 2>&1", tempfile, resolutionratio, wmfont, nickname, wmfont, nickname, watermarkimage, precision, topath)
		}
		// cmd = exec.Command("convert", url, "-resize", resolutionratio, watermarkimage, "-gravity", "southeast", "-geometry", "+0+0", "-gravity", "southeast", "-fill", "white", "-font", wmfont, "-pointsize", "16", "-draw", wmnickname, "-quality", precision, "-composite", tourl)
		//GIF
		// convert i.gif -coalesce -set dispose previous -gravity southeast -geometry +0+0 -fill white -font /usr/share/fonts/google-noto-cjk/NotoSerifCJK-Regular.ttc -pointsize 16 -draw "image over 5,5 24,24 '/www/wwwroot/zyz/img/logo.png' text 29,5 '@择择#213'" -coalesce -resize 256x -layers Optimize 2.gif
	}
	sh = fmt.Sprintf("cat %s >>convertimage.log && echo >>convertimage.log && date >>convertimage.log && %s && echo ===== >>convertimage.log && rm -f %s", shfile, sh, shfile)
	println("===sh:", sh)
	content := []byte(sh)
	err := ioutil.WriteFile(shfile, content, 0777)
	check(1, "ioutil.WriteFile", err)
	cmd := exec.Command("sh", shfile)
	stdout, err := cmd.StdoutPipe()
	check(1, "cmd.StdoutPipe", err)
	defer stdout.Close()
	err = cmd.Start()
	check(1, "cmd.Start", err)
	opBytes, err := ioutil.ReadAll(stdout)
	check(1, "ioutil.ReadAll", err)
	fmt.Println(string(opBytes))
}

//对错误检查的封装
func check(code int, msg string, err error) {
	if err != nil {
		fmt.Printf("msg:%s\nerror:%s\n", msg, err)
		if code == 1 {
			os.Exit(code)
		}
	}
}
