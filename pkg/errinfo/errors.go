package errinfo

import (
	"errors"
	"net/http"
)

type errInf struct {
	Msg    string `json:"message"`
	Status int    `json:"status"`
}

var (
	InternalServerErr   = errors.New("internal server error")
	InvalidExcelData    = errors.New("invalid excel data")
	InvalidExcelEmpty   = errors.New("invalid excel empty")
	InvalidNotFile      = errors.New("invalid, file not found in request")
	InvalidTaskID       = errors.New("invalid task id")
	InvalidSupplierID   = errors.New("invalid supplier id")
	InvalidUserID       = errors.New("invalid supplier id")
	InvalidUserOrSupplierID = errors.New("invalid user or supplier id")
	InvalidLimit        = errors.New("invalid limit")
	InvalidOffset       = errors.New("invalid offset")
	InvalidTecDocParams = errors.New("invalid tecdoc params")

	constErrs = map[error]errInf{
		InternalServerErr: {
			Msg:    "Внутренняя ошибка на сервере, обратитесь к разработчикам",
			Status: http.StatusInternalServerError,
		},
		InvalidExcelData: {
			Msg:    "В excel указаны некорректные данные",
			Status: http.StatusBadRequest,
		},
		InvalidExcelEmpty: {
			Msg:    "в таблице нет данных",
			Status: http.StatusBadRequest,
		},
		InvalidNotFile: {
			Msg:    "нет файла в запросе",
			Status: http.StatusBadRequest,
		},
		InvalidTaskID: {
			Msg:    "Некорректный id задания",
			Status: http.StatusBadRequest,
		},
		InvalidSupplierID: {
			Msg:    "Некорректный id поставщика",
			Status: http.StatusBadRequest,
		},
		InvalidUserOrSupplierID: {
			Msg:    "Некорректный id поставщика или пользователя",
			Status: http.StatusBadRequest,
		},
		InvalidUserID: {
			Msg:    "Некорректный id пользователя",
			Status: http.StatusBadRequest,
		},
		InvalidLimit: {
			Msg:    "Некорректный параметр limit",
			Status: http.StatusBadRequest,
		},
		InvalidOffset: {
			Msg:    "Некорректный параметр offset",
			Status: http.StatusBadRequest,
		},
		InvalidTecDocParams: {
			Msg:    "Некорректные название бренда и номер артикула",
			Status: http.StatusBadRequest,
		},
	}
)

func GetErrorInfo(err error) (int, string) {
	info, found := constErrs[err]
	if found {
		return info.Status, info.Msg
	}
	return 0, ""
}
