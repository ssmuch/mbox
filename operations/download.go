package operations

import (
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"pandora/conf"
	"pandora/constants"
	"pandora/models"

	"github.com/sirupsen/logrus"
)

// DownloadSubject download images by subjectID
func DownloadSubject(sID uint64) {
	images := GetNotDownloadedImagesBySubjectID(sID)
	for _, i := range images {
		//go download(i)
		x := rand.Intn(10)
		time.Sleep(time.Duration(x) * time.Second)
		go func(i models.Image) {
			err := download(i)
			if err != nil {
				logrus.Printf("%v", err)
			}
		}(i)

	}
}

func download(img models.Image) error {
	resp, err := http.Get(img.URL)
	if err != nil {
		logrus.Printf("%v", err)
		return err
	}
	defer resp.Body.Close()

	// Build path
	imgByte, err := ioutil.ReadAll(resp.Body)

	var fh *os.File
	cTitle := GetCategoryTitleByID(img.CategoryID)
	file := conf.Setup.Section("download").Key("image_path").String() + cTitle + "/" + img.Title + "/" + img.Name + ".jpg"
	fh, err = os.Create(file)
	if err != nil {
		logrus.Fatalf("Failed to create img file: %s", file)
	} else {
		logrus.Printf("Creating: %s", file)
	}

	defer fh.Close()
	fh.Write(imgByte)

	// Save base64 to db
	db := conf.GlobalDb.Get()

	img.Base64 = base64.StdEncoding.EncodeToString(imgByte)
	img.DownloadStatus = constants.DOWNLOAD_STATUS__DONE
	db.Save(&img)

	thumbID := FetchThumbImageBySubjectID(img.SubjectID)
	if img.ID == thumbID {
		s := models.Subject{
			PandoraObj: models.PandoraObj{
				ID: img.SubjectID,
			},
		}
		s.DownloadStatus = constants.DOWNLOAD_STATUS__DONE
		db.Save(&s)
	}
	return err
}
