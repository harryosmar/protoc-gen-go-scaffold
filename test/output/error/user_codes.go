package error

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	userCodeErrMap = map[CodeErr]CodeErrEntity{
		ErrUserNotFound: {Code: "ERR404User2a12b6d7", Status: http.StatusNotFound, GrpcCode: codes.NotFound, Message: "user not found"},
	}
)

func InitUserCode() {
	for code, entity := range userCodeErrMap {
		codeErrMap[code] = entity
	}
}

// Error code constants
const (
	ErrUserNotFound CodeErr = "ErrUserNotFound"
)
