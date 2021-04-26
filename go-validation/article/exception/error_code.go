package exception

import "net/http"

var (
	/*
	   system  : 10000
	   article : 20000
	*/

	ErrInvalidRequest = NewCustomError(http.StatusBadRequest, 10100, "Request is invalid")
	ErrBinding        = NewCustomError(http.StatusUnprocessableEntity, 10200, "Request binding error")

	//todo : 실제 request에서 요청하는 ID도 메시지에 포함시키면 좋을 듯하다
	ErrArticleNotFound = NewCustomError(http.StatusNotFound, 20100, "Article is not found")
)
