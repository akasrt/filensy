package file

import (
	"log"
	"os"
	"time"

	"github.com/akasrt/filensy/internal/config/env"
	"github.com/akasrt/filensy/internal/cryptox"
	"github.com/akasrt/filensy/internal/filestore"
	"github.com/akasrt/filensy/internal/util/errorx"
	"github.com/rs/xid"
)

type Service interface {
	GetMetaData(code, token string) (RSFileData, error)
	Upload(rq RQFileData) (RSFileData, error)
	Download(code, token string) (*os.File, RSFileData, error)
	Delete(code, token string) error
}

func NewService() Service {
	fs := filestore.NewFileStore(env.GetEnv(env.FileRoot))

	return &service{
		storage:   NewStorage(),
		fileStore: fs,
	}
}

type service struct {
	storage   Storage
	fileStore filestore.FileStore
}

func (s *service) GetMetaData(code, token string) (RSFileData, error) {
	fileData, err := s.storage.Get(code)
	if err != nil {
		return RSFileData{}, err
	}

	isValidToken := cryptox.VerifyToken(token, fileData.Token)
	if !isValidToken {
		return RSFileData{}, errorx.NewUnauthorizedError(nil)
	}

	return fileData.MapToResponse(nil), nil
}

func (s *service) Upload(rq RQFileData) (RSFileData, error) {
	maxRetries := 10
	token, tokenHash, err := cryptox.GenerateFileToken()
	if err != nil {
		return RSFileData{}, err
	}

	storageKey, size, err := s.fileStore.Save(rq.Reader)
	if err != nil {
		return RSFileData{}, err
	}

	now := time.Now()
	var expiresAt time.Time
	if rq.TTL > 0 {
		expiresAt = now.Add(rq.TTL)
	}

	fileData := FileData{
		ID:         xid.New().String(),
		Token:      tokenHash,
		Name:       rq.Name,
		StorageKey: storageKey,
		Size:       size,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
	}

	var savedData FileData
	for i := 1; i <= maxRetries; i++ {
		code, err := cryptox.GenerateFileCode()
		if err != nil {
			s.fileStore.Delete(fileData.StorageKey)
			return RSFileData{}, err
		}
		fileData.Code = code

		savedData, err = s.storage.Create(fileData)
		if err != nil {
			if s.storage.IsDuplicateErr(err) {
				continue
			} else {
				s.fileStore.Delete(storageKey)
				return RSFileData{}, err
			}
		}

		if i == maxRetries {
			log.Println("Max retries reached file insert")
			s.fileStore.Delete(storageKey)
			return RSFileData{}, errorx.NewInternalServerError(nil)
		}
	}

	return savedData.MapToResponse(&token), nil
}

func (s *service) Download(code, token string) (*os.File, RSFileData, error) {
	fileData, err := s.storage.Get(code)
	if err != nil {
		return nil, RSFileData{}, err
	}

	isValidToken := cryptox.VerifyToken(token, fileData.Token)
	if !isValidToken {
		return nil, RSFileData{}, errorx.NewUnauthorizedError(nil)
	}

	file, err := s.fileStore.Get(fileData.StorageKey)
	if err != nil {
		return nil, RSFileData{}, err
	}

	return file, fileData.MapToResponse(nil), nil
}

func (s *service) Delete(code, token string) error {
	fileData, err := s.storage.Get(code)
	if err != nil {
		return err
	}

	isValidToken := cryptox.VerifyToken(token, fileData.Token)
	if !isValidToken {
		return errorx.NewUnauthorizedError(nil)
	}

	err = s.storage.Delete(code)
	if err != nil {
		return err
	}

	err = s.fileStore.Delete(fileData.StorageKey)
	if err != nil {
		return err
	}

	return nil
}
