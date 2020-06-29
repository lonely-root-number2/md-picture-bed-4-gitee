package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

//Create Date:2020 六月 28 16:16:55 周日
//Auther:DY
//Desc:
//@Email:3309647521@qq.com
func getIMGFile(s string) ([]string, int) {
	// 返回包含根路径，数量因此减一
	ret := make([]string, 0, 10)
	i := 0
	pth := path.Dir(filepath.ToSlash(s))
	filepath.Walk(pth, func(path string, info os.FileInfo, err error) error {
		i++
		ret = append(ret, path)
		return nil
	})
	return ret[1:i], i - 1
}
func getGiteeName(s string) string {
	pth := path.Dir(filepath.ToSlash(s))
	return filepath.Base(pth)
}

type gitee struct {
	accessToken      string
	userName         string
	repositoriesName string
	message          string
}

func (g *gitee) upload(s string) {
	now := getTime()
	now = now + "_"
	imgSlice, _ := getIMGFile(s)
	name := now + getGiteeName(s)
	if g.message == "" {
		// 自定义commit message
		g.message = "upload by DY"
	}
	url := "https://gitee.com/api/v5/repos/"    //hades0/imgStore/contents/c.png"
	client := &http.Client{}
	for k, v := range imgSlice {
		if path.Ext(v) == "md" {
			continue
		}
		fmt.Printf("is uploading %dth\n", k+1)
		f, _ := os.Open(filepath.ToSlash(v))
		buf := make([]byte, 1024000)
		n, _ := f.Read(buf)
		rets := base64.StdEncoding.EncodeToString(buf[:n])
		fmt.Println(url + g.userName + g.repositoriesName + `/contents/` + name + `/` + filepath.Base(v))
		req, _ := http.NewRequest("POST", url+g.userName+`/`+g.repositoriesName+`/contents/`+name+`/`+filepath.Base(v), strings.NewReader(`{"access_token":"`+g.accessToken+`","content":"`+rets+`","message":"`+g.message+`"}`))
		req.Header.Set("Content-Type", `application/json;charset=UTF-8`)
		resp, _ := client.Do(req)
		result, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(result))
		resp.Body.Close()
		f.Close()
	}
	// 处理替换逻辑
	fmd, _ := os.OpenFile(s, os.O_RDWR, 0644)
	bufmd, _ := ioutil.ReadAll(fmd)
	fmd.Close()
	sbuf := string(bufmd)
	for _, v := range imgSlice {
		sbuf = strings.Replace(sbuf, `./`+filepath.Base(v), `https://`+g.userName+`.gitee.io/`+g.repositoriesName+`/`+name+`/`+filepath.Base(v), 1)
	}
	os.Rename(s, filepath.Dir(s)+`\`+filepath.Base(s)+`.bak`)
	fmd, _ = os.OpenFile(s, os.O_CREATE, 0644)
	fmd.Write([]byte(sbuf))
	fmd.Close()
}
func getTime() string {
	// 你可以自定义路径
	return time.Now().Format("2006-01-02-15-04")
}
func main() {

	g := &gitee{
		accessToken:      `你的token`,
		userName:         "用户名",
		repositoriesName: "仓库名",
	}
	g.upload("markdown文件路径") //`C:\Users\Administrator\Desktop\微服务架构\微服务.md`

}
