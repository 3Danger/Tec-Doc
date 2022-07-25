package services

import (
	"context"
)

//go:generate tg client -go --services . --outPath ../content

// @tg http-prefix=content
// @tg jsonRPC-server content
type Source interface {
	// @tg summary=`Работа с карточками товаров`
	//tg desc=`описание метода в документации swagger`
	// @tg http-path=http://source.content-card.svc.k8s.stage-dp/source/migration
	Migration(ctx context.Context, attributes ...map[string]interface{}) (err error)
}
