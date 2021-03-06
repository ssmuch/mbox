package utils

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func GetHTML(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		logrus.Fatalf("Failed to load: %s", url)
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	return string(html)
}

func GetPageLimit(html string) int {
	reg, _ := regexp.Compile(`尾(\d+)页`)
	dst := []byte("")
	template := "$1"
	subj := reg.FindString(html)
	match := reg.FindStringSubmatchIndex(subj)
	tmp := reg.ExpandString(dst, template, subj, match)

	num, _ := strconv.Atoi(string(tmp))
	return num
}

func RandStr(prefix string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	return prefix + strconv.Itoa(rand.Intn(1000))
}

func ProcessDir(dirPath string) (err error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		var mode os.FileMode = 0777
		err = os.MkdirAll(dirPath, mode)
	}
	return
}

func Basename(s string) string {
	slash := strings.LastIndex(s, "/") // -1 if "/" not found
	s = s[slash+1:]
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}
