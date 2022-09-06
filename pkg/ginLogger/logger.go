package ginLogger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func Logger(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		nowTime := time.Now()
		start := nowTime.Format("2006-01-02 15:04:05")
		c.Next()
		latency := fmt.Sprint(time.Now().Sub(nowTime))
		Method := methodColor(c.Request.Method) + reset
		path := magenta + c.Request.URL.Path + reset
		raw := magenta + c.Request.URL.RawQuery + reset
		code := statusColor(c.Writer.Status())
		ip := blue + "ip: " + c.ClientIP() + reset
		errMsg := red + c.Errors.ByType(1).String() + reset
		bodyLen := blue + "body-length: " + strconv.Itoa(c.Writer.Size()) + reset
		contentType := cyan + "Content-Type: " + c.GetHeader("Content-Type") + reset
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + raw
		}
		_, err := fmt.Fprintf(out,
			"%s %s %s %s %s %s %s %s %s \n",
			start, latency, Method, path,
			code, ip,
			errMsg, bodyLen, contentType)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func statusColor(status int) (result string) {
	defer func() { result = result + reset }()
	if status == 200 {
		return green + "code: " + "200"
	} else if status < 300 {
		return yellow + "code: " + strconv.Itoa(status)
	} else if status < 300 {
		return magenta + "code: " + strconv.Itoa(status)
	} else if status < 400 {
		return cyan + "code: " + strconv.Itoa(status)
	}
	return red + "code: " + strconv.Itoa(status)
}

func methodColor(method string) (result string) {
	defer func() { result = result + reset }()
	switch method {
	case http.MethodGet:
		return blue + method
	case http.MethodPost:
		return cyan + method
	case http.MethodPut:
		return yellow + method
	case http.MethodDelete:
		return red + method
	case http.MethodPatch:
		return green + method
	case http.MethodHead:
		return magenta + method
	case http.MethodOptions:
		return white + method
	default:
		return method
	}
}
