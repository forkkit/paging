package paging

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// ValidateLimitMarker returns true if limit and offset/cursor values are valid
func ValidateLimitMarker(limit int64, marker int64) bool {
	values := []int64{limit, marker}

	for _, v := range values {
		if v < 0 {
			return false
		}
	}

	return true
}

// GetLimitFromRequest returns current limit.
func GetLimitFromRequest(request *http.Request, options *Options) int64 {
	var (
		limit int64
		err   error
	)

	requestLimit := request.URL.Query().Get(options.LimitKeyName)

	if requestLimit != "" {
		limit, err = strconv.ParseInt(requestLimit, 10, 64)
		if err != nil {
			limit = options.DefaultLimit

		}
		if options.MaxLimit > 0 && limit > options.MaxLimit {
			limit = options.MaxLimit
		}
	} else {
		limit = options.DefaultLimit
	}

	return limit
}

// GetOffsetFromRequest returns current offset.
func GetOffsetFromRequest(request *http.Request, options *Options) int64 {
	var (
		offset int64
		err    error
	)

	requestOffset := request.URL.Query().Get(options.OffsetKeyName)

	if requestOffset != "" {
		offset, err = strconv.ParseInt(requestOffset, 10, 64)
		if err != nil {
			offset = 0
		}
	} else {
		offset = 0
	}

	return offset
}

// GetCursorFromRequest returns current cursor.
func GetCursorFromRequest(request *http.Request, options *Options) int64 {
	var (
		cursor int64
		err    error
	)

	requestCursor := request.URL.Query().Get(options.CursorKeyName)

	if requestCursor != "" {
		cursor, err = strconv.ParseInt(requestCursor, 10, 64)
		if err != nil {
			cursor = 0
		}
	} else {
		cursor = 0
	}

	return cursor
}

// GenerateOffsetURI generates the pagination URI.
func GenerateOffsetURI(limit int64, offset int64, options *Options) string {
	return fmt.Sprintf(
		"?%s=%d&%s=%d",
		options.LimitKeyName,
		limit,
		options.OffsetKeyName,
		offset)
}

// GenerateCursorURI generates the pagination URI for cursor system.
func GenerateCursorURI(limit int64, cursor interface{}, options *Options) string {
	return fmt.Sprintf(
		"?%s=%d&%s=%d",
		options.LimitKeyName,
		limit,
		options.CursorKeyName,
		cursor)
}

// GetPaginationType returns the pagination type "offeset|cursor"
// (use constant CursorType or OffsetType)
// return OffsetType by default
func GetPaginationType(request *http.Request, options *Options) string {
	if options == nil {
		options = NewOptions()
	}

	if cursor := GetCursorFromRequest(request, options); cursor > 0 {
		return CursorType
	}

	return OffsetType
}

// Last gets the last element ID value of array.
func Last(arr interface{}, field string) (interface{}, error) {
	value := reflect.ValueOf(arr)
	valueType := value.Type()

	kind := value.Kind()
	if kind == reflect.Ptr {
		value = value.Elem()
		valueType = value.Type()
		kind = value.Kind()
	}

	if kind == reflect.Array || kind == reflect.Slice {
		if value.Len() == 0 {
			return 0, nil
		}
		item := value.Index(value.Len() - 1)
		cursor := item.FieldByName(field)
		if cursor.Kind() == reflect.Struct {
			return cursor.Interface().(time.Time), nil
		}

		return cursor.Int(), nil
	}

	return 0, fmt.Errorf("Type %s is not supported by Last", valueType.String())
}
