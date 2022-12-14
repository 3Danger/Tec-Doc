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
	InternalServerErr        = errors.New("internal server error")
	InvalidExcelData         = errors.New("invalid excel data")
	InvalidExcelEmpty        = errors.New("invalid excel empty")
	InvalidExcelLimit        = errors.New("invalid excel limit")
	InvalidNotFile           = errors.New("invalid, file not found in request")
	InvalidTaskID            = errors.New("invalid task id")
	InvalidSupplierID        = errors.New("invalid supplier id")
	InvalidUserID            = errors.New("invalid supplier id")
	InvalidUserOrSupplierID  = errors.New("invalid user or supplier id")
	InvalidLimit             = errors.New("invalid limit")
	InvalidOffset            = errors.New("invalid offset")
	InvalidTecDocParams      = errors.New("invalid tecdoc params")
	SupplierIsNotUUID        = errors.New("supplier is not uuid")
	FailOldSupplierID        = errors.New("can't get old supplier ID")
	CheckAcessError          = errors.New("check access error")
	NoTecDocArticlesFound    = errors.New("no articles found")
	MoreThanOneArticlesFound = errors.New("found more than 1")
	NoTecDocBrandFound       = errors.New("no brand found")

	constErrs = map[error]errInf{
		MoreThanOneArticlesFound: {
			Msg:    "Количество найденных артикулов по запросу больше 1",
			Status: http.StatusBadRequest,
		},
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
		InvalidExcelLimit: {
			Msg:    "Превышен лимит объектов в таблице",
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
		SupplierIsNotUUID: {
			Msg:    "ID поставщика не является UUID",
			Status: http.StatusUnauthorized,
		},
		FailOldSupplierID: {
			Msg:    "нельзя получить старый ID поставщика",
			Status: http.StatusUnauthorized,
		},
		CheckAcessError: {
			Msg:    "ошибка доступа, проверьте имеются ли необходимые права",
			Status: http.StatusUnauthorized,
		},
		NoTecDocArticlesFound: {
			Msg:    "товары не найдены, проверьте корректно ли переданы бренд и артикул",
			Status: http.StatusUnauthorized,
		},
		NoTecDocBrandFound: {
			Msg:    "бренд не найден, проверьте корректно ли переданы параметры",
			Status: http.StatusUnauthorized,
		},
	}
)

func GetErrorInfo(err error) (int, string) {
	info, found := constErrs[err]
	if !found {
		info, _ = constErrs[InternalServerErr]
	}
	return info.Status, info.Msg
}
