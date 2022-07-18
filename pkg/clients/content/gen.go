package content

import (
	"context"
)

//go:generate tg client -go --services . --outPath ../generated_client

// @tg http-prefix=content
// @tg jsonRPC-server content
type ContentClient interface {
	// @tg summary=`Работа с карточками товаров`
	//tg desc=`описание метода в документации swagger`
	//tg key.type=string
	//tg key.format=uuid
	// @tg http-path=/content/push
	//tg http-headers=userID|X-User-Id
	PushContent(ctx context.Context, attributes ...map[string]interface{}) (err error)
}
