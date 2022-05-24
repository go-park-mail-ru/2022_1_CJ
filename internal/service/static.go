package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type StaticService interface {
	UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
	UploadFile(fileHeader *multipart.FileHeader) (string, error)
	UploadFileChat(uuid string, extension string, binaryData string) (string, error)
}

type staticServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *staticServiceImpl) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	uuid, err := core.GenUUID()
	if err != nil {
		return "", err
	}

	filename := uuid + filepath.Ext(fileHeader.Filename)

	dst, err := os.Create("/opt/pics/" + filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	_, err = svc.db.LikeRepo.CreateLike(ctx, &core.Like{Subject: filename})
	if err != nil {
		svc.log.Errorf("CreateLike error: %s", err)
		return "", err
	}

	url := "/" + filename
	return url, nil
}

func (svc *staticServiceImpl) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	uuid, err := core.GenUUID()
	if err != nil {
		return "", err
	}

	filename := uuid + filepath.Ext(fileHeader.Filename)

	dst, err := os.Create("/opt/files/" + filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	url := "/" + filename
	return url, nil
}

func (svc *staticServiceImpl) UploadFileChat(uuid string, extension string, binaryData string) (string, error) {
	filename := uuid + extension

	dst, err := os.Create("/opt/files/" + filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = dst.Write([]byte(binaryData)); err != nil {
		return "", err
	}

	url := "/" + filename
	return url, nil
}

func NewStaticService(log *logrus.Entry, db *db.Repository) StaticService {
	return &staticServiceImpl{log: log, db: db}
}
