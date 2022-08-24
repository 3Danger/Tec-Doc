package services

import (
	"context"
)

//go:generate tg client -go --services . --outPath ../abac

// @tg http-prefix=abac
// @tg jsonRPC-server log metrics
type ABAC interface {
	// @tg summary=`Проверка наличия права в контесте scope`
	// @tg desc=`Проверка наличия у пользователя права **feature**. **userID** автоматически заполняется из http заголовка X-User-Id`
	// @tg key.type=string
	// @tg key.format=uuid
	// @tg http-path=/access/check
	// @tg http-headers=userID|X-User-Id
	CheckAccess(ctx context.Context, scope, featureKey string, userID *uint64, key [16]byte, attributes ...map[string]interface{}) (decision bool, err error)
}
