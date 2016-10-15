package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"syscall"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// BingRoot Bing Web Site Root.
	BingRoot = "http://cn.bing.com/"
	// GetTimeGap How Long to Get a New Images.
	GetTimeGap = 24 * 60 * 60
	// Version  for application.
	Version = "1.0.0"
	// LogFile log every picture url.
	LogFile = "./img_url.txt"
	// PicRegexp regexp pictures root.
	PicRegexp = `[^"]+.jpg`
	// O_APPEND append string when write file.
	O_APPEND = syscall.O_APPEND
)

// regexp example:
// g_img={url: "http://s.cn.bing.net/az/hprichbg/rb/LacsdesCheserys_ZH-CN10032851647_1920x1080.jpg",id:'bgDiv'

// SaveRoot The Root You Save Images.
var SaveRoot = kingpin.Flag("root", "the root you want to save these pictures.").Short('r').Required().String()

func init() {
	kingpin.Parse()
	kingpin.Version(Version)
}

func main() {
	// Start right now.
	HandleTime := time.NewTimer(time.Millisecond)

	for {
		select {
		case <-HandleTime.C:
			GetAndSave()
			HandleTime.Reset(GetTimeGap * time.Second)
		}
	}
}

// GetAndSave bing pictures
func GetAndSave() {
	content, status := getBingContent(BingRoot)
	if status != http.StatusOK {
		log.Println(errors.New("无法访问必应主站，请检查你的网络"))
		os.Exit(-1)
	}

	PicURL := findPicURL(content)
	fmt.Println("PicURL =>", PicURL)
	err := getPic(PicURL)
	if err != nil {
		// TODO: Send E-mail to Me.
		panic(err)
	}
	err = logFile(PicURL)
	if err != nil {
		panic(err)
	}
	log.Println("Get And Save Today Picture Done.")
}

func getBingContent(url string) (content string, status int) {
	res, err := http.Get(url)
	if err != nil {
		// TODO: Send Email to Me.
		log.Println("Get Bing Connect Error:", err)
		return "", http.StatusGatewayTimeout
	}
	defer res.Body.Close()
	bys, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Get Bing Content Error:", err.Error())
		return "", http.StatusBadRequest
	}
	status = http.StatusOK
	content = string(bys)
	return
}

func getPic(url string) (err error) {
	picRes, err := http.Get(url)
	if err != nil {
		return err
	}

	defer picRes.Body.Close()

	err = saveFile(picRes.Body)
	if err != nil {
		log.Println("save file error:", err)
		return err
	}
	return nil
}

func findPicURL(content string) (url string) {
	return regexp.MustCompile(PicRegexp).FindString(content)
}

func saveFile(rc io.ReadCloser) error {
	date := time.Now().Format("2006-01-02")
	fileName := verifyRoot(*SaveRoot) + date + ".jpg"

	file, err := os.Create(fileName)
	if err != nil {
		if err == os.ErrExist {
			file, _ = os.Open(fileName)
		} else {
			return err
		}
	}
	_, err = io.Copy(file, rc)
	defer file.Close()
	return err
}

func logFile(urlContent string) error {
	var file *os.File
	var err error
	if isExist(LogFile) {
		file, err = os.OpenFile(LogFile, O_APPEND|os.O_WRONLY, os.ModeAppend)
	} else {
		log.Println("log file not exist, creating...")
		file, err = os.Create(LogFile)
	}
	if err != nil {
		return err
	}

	_, err = file.WriteString("\n" + time.Now().Format("2006-01-02") + " : \n" + urlContent)
	if err != nil {
		log.Println("write picture URL to file error:", err)
	}
	defer file.Close()
	return err
}

func verifyRoot(root string) string {
	if root == "" {
		return "./"
	}

	if root[len(root)-1] != '/' {
		return root + "/"
	}
	return root
}

func isExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
