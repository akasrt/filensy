package file

import (
	"os"
	"time"

	"github.com/akasrt/filensy/internal/cryptox"
	"github.com/akasrt/filensy/internal/filestore"
	"github.com/akasrt/filensy/internal/util/errorx"
	"github.com/rs/xid"
)

const (
	visibilityPublic  = "public"
	visibilityPrivate = "private"
)

type Service interface {
	GetMetaData(code, token string) (RSFileData, error)
	Upload(rq RQFileData) (RSFileData, error)
	Download(code, token string) (*os.File, RSFileData, error)
	Delete(code, token string) error
}

func NewService() Service {
	fs := filestore.Get()

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

	if fileData.Visibility == visibilityPrivate {
		isValidToken := cryptox.VerifyToken(token, fileData.Token)
		if !isValidToken {
			return RSFileData{}, errorx.WrapUnauthorizedError(nil)
		}
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
		ID:           xid.New().String(),
		Token:        tokenHash,
		Name:         rq.Name,
		StorageKey:   storageKey,
		Size:         size,
		Visibility:   rq.Visibility,
		Is_Encrypted: rq.Is_Encrypted,
		CreatedAt:    now,
		ExpiresAt:    expiresAt,
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
		if err == nil {
			break
		}

		if !s.storage.IsDuplicateErr(err) {
			s.fileStore.Delete(storageKey)
			return RSFileData{}, err
		}

		if i == maxRetries {
			s.fileStore.Delete(storageKey)
			return RSFileData{}, errorx.WrapInternalServerError(nil)
		}
	}
	return savedData.MapToResponse(&token), nil
}

func (s *service) Download(code, token string) (*os.File, RSFileData, error) {
	fileData, err := s.storage.Get(code)
	if err != nil {
		return nil, RSFileData{}, err
	}

	if fileData.Visibility == visibilityPrivate {
		isValidToken := cryptox.VerifyToken(token, fileData.Token)
		if !isValidToken {
			return nil, RSFileData{}, errorx.WrapUnauthorizedError(nil)
		}
	}

	file, err := s.fileStore.Get(fileData.StorageKey)
	if err != nil {
		return nil, RSFileData{}, errorx.WrapInternalServerError(err)
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
		return errorx.WrapUnauthorizedError(nil)
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
