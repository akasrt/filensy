package file

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/akasrt/filensy/internal/util/errorx"
	"github.com/akasrt/filensy/internal/util/httputil"
	"github.com/labstack/echo/v5"
)

type Handler interface {
	GetMetaData(c *echo.Context) error
	Upload(c *echo.Context) error
	Download(c *echo.Context) error
	Delete(c *echo.Context) error
}

func NewHandler() Handler {
	return &handler{
		service: NewService(),
	}
}

type handler struct {
	service Service
}

func (h *handler) GetMetaData(c *echo.Context) error {
	code, token, err := parseRequest(c)
	if err != nil {
		return err
	}

	data, err := h.service.GetMetaData(code, token)
	if err != nil {
		return err
	}

	resp := httputil.NewResponse("", data)
	return c.JSON(http.StatusOK, resp)
}

func (h *handler) Upload(c *echo.Context) error {
	contentType := strings.ToLower(c.Request().Header.Get("Content-Type"))
	if !strings.HasPrefix(contentType, "application/octet-stream") {
		return errorx.New(http.StatusUnsupportedMediaType, "Content-Type must be application/octet-stream")
	}

	var fileData RQFileData
	ttl, err := parseTTL(c)
	if err != nil {
		return err
	}

	fileData.Name = c.QueryParam("name")
	encStr, err := strconv.ParseBool(c.QueryParam("enc"))
	if err != nil {
		return errorx.Wrap(err, 400, "invalid query param")
	}
	fileData.Is_Encrypted = encStr
	fileData.TTL = ttl
	visibility := c.QueryParam("visibility")
	if visibility == "" {
		fileData.Visibility = visibilityPrivate
	} else {
		fileData.Visibility = visibility
	}

	fileData.Reader = c.Request().Body
	c.Validate(fileData)

	fd, err := h.service.Upload(fileData)
	if err != nil {
		return err
	}

	resp := httputil.NewResponse("upload success", fd)
	return c.JSON(http.StatusCreated, resp)
}

func (h *handler) Download(c *echo.Context) error {
	code, token, err := parseRequest(c)
	if err != nil {
		return err
	}

	file, metaData, err := h.service.Download(code, token)
	if err != nil {
		return err
	}
	defer file.Close()

	setFileHeaders(c, metaData)

	return c.Stream(
		http.StatusOK,
		"application/octet-stream",
		file,
	)
}

func (h *handler) Delete(c *echo.Context) error {
	code, token, err := parseRequest(c)
	if err != nil {
		return err
	}

	err = h.service.Delete(code, token)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func parseRequest(c *echo.Context) (string, string, error) {
	code := c.Param("code")
	if code == "" {
		return "", "", errorx.New(http.StatusBadRequest, "missing code in param")
	}

	token := c.Request().Header.Get("X-File-Token")

	return code, token, nil
}

func parseTTL(c *echo.Context) (time.Duration, error) {
	defaultTTL := time.Hour * 24 * 30
	maxTTL := time.Hour * 24 * 365

	ttlStr := c.QueryParam("ttl")
	if ttlStr == "" {
		return defaultTTL, nil
	}

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return 0, errorx.Wrap(err, http.StatusBadRequest, "invalid ttl")
	}

	if ttl > maxTTL {
		return maxTTL, nil
	}

	if ttl <= 0 {
		return 0, errorx.New(http.StatusBadRequest, "invalid ttl")
	}

	return ttl, nil
}

func setFileHeaders(c *echo.Context, metaData RSFileData) {
	c.Response().Header().Set("X-File-Name", metaData.Name)
	c.Response().Header().Set("X-Is-Encrypted", fmt.Sprint(metaData.Is_Encrypted))

	var fileName string
	if metaData.Is_Encrypted {
		fileName = "enc_" + metaData.Name + ".bin"
	} else {
		fileName = metaData.Name
	}
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))

	c.Response().Header().Set("Content-Length", strconv.FormatUint(metaData.Size, 10))
}
