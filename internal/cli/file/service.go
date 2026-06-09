package file

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/akasrt/filensy/internal/cli/api"
	"github.com/akasrt/filensy/internal/cli/localstore"
	"github.com/akasrt/filensy/internal/config/userconfig"
	"github.com/akasrt/filensy/internal/cryptox"
	"github.com/akasrt/filensy/internal/util/errorx"
)

const (
	visibilityPublic  = "public"
	visibilityPrivate = "private"
)

type FileService interface {
	UploadFile(path string, opts FileOptions) (RSFileData, error)
	GetFile(dir, code, token, password string) error
	FindFile(code, token string) (RSFileData, error)
	DeleteFile(code, token string) error
}

type fileService struct {
	config     userconfig.Config
	localStore localstore.Store
}

func NewFileService() (FileService, error) {
	store, err := localstore.New()
	if err != nil {
		return nil, err
	}

	return &fileService{
		config:     userconfig.GetConfig(),
		localStore: store,
	}, nil
}

func (f *fileService) UploadFile(path string, opts FileOptions) (RSFileData, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return RSFileData{}, errorx.ErrFileNotExists
		}
		return RSFileData{}, err
	}

	if info.IsDir() {
		return RSFileData{}, errorx.ErrFileNotExists
	}

	file, err := os.Open(path)
	if err != nil {
		return RSFileData{}, err
	}
	defer file.Close()

	fileName := filepath.Base(path)

	var payloadReader io.Reader = file
	if opts.IsEncrypted {
		pr, pw := io.Pipe()

		go func() {
			err = cryptox.Encrypt(payloadReader, pw, opts.Password)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			pw.Close()
		}()

		payloadReader = pr
	}

	var visibility string
	if opts.IsPublic {
		visibility = visibilityPublic
	} else {
		visibility = visibilityPrivate
	}
	status, resp, err := api.PostFile(fileName, opts.TTL, visibility, opts.IsEncrypted, payloadReader)
	if err != nil {
		return RSFileData{}, err
	}

	if status != http.StatusCreated {
		return RSFileData{}, errorx.WrapServerError(status, resp)
	}

	fileData, err := parseResponse(resp.Data)
	if err != nil {
		return RSFileData{}, err
	}

	local := localstore.FileData{
		Code:     fileData.Code,
		Token:    *fileData.Token,
		Password: opts.Password,
		Name:     fileData.Name,
	}
	err = f.localStore.Create(local)
	if err != nil {
		return RSFileData{}, err
	}

	return fileData, nil
}

func (f *fileService) GetFile(dir, code, token, password string) error {
	if token == "" || password == "" {
		localData, exists := f.localStore.Get(code)
		if exists {
			if localData.Token != "" && token == "" {
				token = localData.Token
			}
			if localData.Password != "" && password == "" {
				password = localData.Password
			}
		}
	}

	status, resp, reader, headers, err := api.GetFile(code, token)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errorx.WrapServerError(status, resp)
	}

	defer reader.Close()

	fileData := parseFileHeader(headers)

	if fileData.Is_Encrypted {
		if password == "" {
			return errorx.ErrPasswordMissing
		}

		pr, pw := io.Pipe()
		go func() {
			err := cryptox.Decrypt(reader, pw, password)
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			pw.Close()
		}()
		reader = pr
	}

	if dir == "" {
		defaultDir := userconfig.GetConfig().Directory
		if defaultDir != "" {
			dir = defaultDir
		} else {
			dir = "."
		}
	}
	err = SaveFile(reader, dir, fileData.Name)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileService) FindFile(code string, token string) (RSFileData, error) {
	if token == "" {
		localData, exists := f.localStore.Get(code)
		if exists {
			token = localData.Token
		}
	}
	status, resp, err := api.GetFileMetadata(code, token)
	if err != nil {
		return RSFileData{}, err
	}

	if status != http.StatusOK {
		return RSFileData{}, errorx.WrapServerError(status, resp)
	}

	fileData, err := parseResponse(resp.Data)
	if err != nil {
		return RSFileData{}, err
	}
	return fileData, nil
}

func (f *fileService) DeleteFile(code, token string) error {
	if token == "" {
		localData, exists := f.localStore.Get(code)
		if exists {
			token = localData.Token
		}
	}
	status, resp, err := api.DeleteFile(code, token)
	if err != nil {
		return err
	}

	if status != http.StatusOK && status != http.StatusNoContent {
		return errorx.WrapServerError(status, resp)
	}

	err = f.localStore.Delete(code)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(data io.Reader, dir, fileName string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		return err
	}
	return nil
}

func parseResponse(data any) (RSFileData, error) {
	var result RSFileData
	b, err := json.Marshal(data)
	if err != nil {
		return RSFileData{}, err
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return RSFileData{}, err
	}

	return result, nil
}

func parseFileHeader(headers map[string]string) RSFileData {
	fileName := headers["X-File-Name"]
	if fileName == "" {
		fileName = headers["Content-Disposition"]
	}
	isEncrypted, _ := strconv.ParseBool(headers["X-Is-Encrypted"])

	return RSFileData{
		Name:         fileName,
		Is_Encrypted: isEncrypted,
	}
}
