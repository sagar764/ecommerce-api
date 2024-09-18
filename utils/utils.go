package utils

import (
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

func PrepareVersionName(version string) string {
	version = strings.ReplaceAll(version, ".", "_")
	return version
}

func ValidateName(data string) bool {
	for _, r := range data {
		if !(unicode.IsLetter(r) || r == '\'' || r == '-' || r == ' ' || r == '_') {
			return false
		}
	}
	return len(data) > 0 && len(data) <= 50
}

// ValidateUUID checks if the provided string is a valid UUID
func ValidateUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// ValidateURL checks if a given string is a valid URL.
func ValidateURL(urlStr string) error {
	re := regexp.MustCompile(consts.UrlRegex)
	if !re.MatchString(urlStr) {
		return errors.New("invalid URL format")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("error parsing URL")
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("URL must have a valid scheme and host")
	}

	return nil
}

func Paginate(page, limit, defaultLimit int) (int, int) {
	if page <= 0 {
		page = consts.DefaultPage
	}
	if (limit <= 0 || limit > consts.MaxLimit) && defaultLimit > 0 {
		limit = defaultLimit
	} else if limit <= 0 {
		limit = consts.DefaultLimit
	}
	return page, limit
}

// MetaDataInfo calculates values for pagination.
func MetaDataInfo(metaData *entities.MetaData) *entities.MetaData {
	if metaData.Total < 1 {
		return nil
	}
	if metaData.CurrentPage*metaData.PerPage < metaData.Total {
		metaData.Next = metaData.CurrentPage + 1
	}
	if metaData.CurrentPage > 1 {
		metaData.Prev = metaData.CurrentPage - 1
	}
	return metaData
}
