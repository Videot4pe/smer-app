package uploader

import (
	"backend/pkg/logging"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/vincent-petithory/dataurl"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Uploader struct {
	Logger *logging.Logger
}

const MaxUploadSize = 0

func GetFiler(logger *logging.Logger) *Uploader {
	return &Uploader{
		Logger: logger,
	}
}

var instance *Uploader
var once sync.Once

func GetUploader(logger *logging.Logger) *Uploader {
	once.Do(func() {
		instance = &Uploader{
			Logger: logger,
		}
	})
	return instance
}

func (u *Uploader) Base64Upload(data string) (string, string, error) {
	dataURL, err := dataurl.DecodeString(data)
	if err != nil {
		u.Logger.Println(err)
		return "", "", err
	}
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", "", err
	}

	name := uniuri.NewLen(64)
	var fileType string
	if dataURL.ContentType() == "image/png" {
		fileType = "png"
	} else if dataURL.ContentType() == "image/jpg" {
		fileType = "jpg"
	}
	path := fmt.Sprintf("./uploads/%d%s.%s", time.Now().UnixNano(), name, fileType)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dataURL.Data); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	return path, name, nil
}

func (u *Uploader) MultipleUpload(files []*multipart.FileHeader) error {
	for _, fileHeader := range files {
		if fileHeader.Size > MaxUploadSize {
			u.Logger.Error("Error")
			return errors.New("Error")
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			u.Logger.Error(err)
			return err
		}

		defer file.Close()

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			u.Logger.Error(err)
			return err
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			u.Logger.Error(err)
			return err
		}

		fl, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		if err != nil {
			u.Logger.Error(err)
			return err
		}

		defer fl.Close()

		_, err = io.Copy(fl, file)
		if err != nil {
			u.Logger.Error(err)
			return err
		}

		pr := &Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(fl, io.TeeReader(file, pr))
		if err != nil {
			u.Logger.Error(err)
			return err
		}
	}
	return nil
}

func (u *Uploader) Upload(path string, file io.Reader) error {
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		u.Logger.Error(err)
		return err
	}

	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(path)))
	if err != nil {
		u.Logger.Error(err)
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		u.Logger.Error(err)
		return err
	}

	u.Logger.Infoln("Upload success")
	return nil
}

func (u *Uploader) GetFileType(file io.Reader) (string, error) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		u.Logger.Error(err)
		return "", err
	}

	filetype := http.DetectContentType(buff)

	return filetype, nil
}
