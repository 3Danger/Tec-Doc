package errinfo

import (
	"errors"
	"net/http"
)

type errInf struct {
	Msg  string `json:"message"` // Error message
	Code int    `json:"status"`  // Http status code
}

var (
	InternalServerErr   = errors.New("internal server error")
	InvalidExcelData    = errors.New("invalid excel data")
	InvalidTaskID       = errors.New("invalid task id")
	InvalidSupplierID   = errors.New("invalid supplier id")
	InvalidUserID       = errors.New("invalid supplier id")
	InvalidLimit        = errors.New("invalid limit")
	InvalidOffset       = errors.New("invalid offset")
	InvalidTecDocParams = errors.New("invalid tecdoc params")

	constErrs = map[error]errInf{
		InternalServerErr: {
			Msg:  "Внутренняя ошибка на сервере, обратитесь к разработчикам",
			Code: http.StatusInternalServerError,
		},
		InvalidExcelData: {
			Msg:  "В excel указаны некорректные данные",
			Code: http.StatusBadRequest,
		},
		InvalidTaskID: {
			Msg:  "Некорректный id задания",
			Code: http.StatusBadRequest,
		},
		InvalidSupplierID: {
			Msg:  "Некорректный id поставщика",
			Code: http.StatusBadRequest,
		},
		InvalidUserID: {
			Msg:  "Некорректный id пользователя",
			Code: http.StatusBadRequest,
		},
		InvalidLimit: {
			Msg:  "Некорректный параметр limit",
			Code: http.StatusBadRequest,
		},
		InvalidOffset: {
			Msg:  "Некорректный параметр offset",
			Code: http.StatusBadRequest,
		},
		InvalidTecDocParams: {
			Msg:  "Некорректные название бренда и номер артикула",
			Code: http.StatusBadRequest,
		},
	}
)

func GetErrorInfo(err error) (string, int) {
	info, found := constErrs[err]
	if found {
		return info.Msg, info.Code
	}
	return "", 0
}
