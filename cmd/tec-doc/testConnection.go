package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"tec-doc/internal/config"
	"tec-doc/internal/web/internalserver"
	"time"
)

func testConnection() {
	//Init
	conf := new(config.Config)
	if err := envconfig.Process("TEC_DOC", conf); err != nil {
		log.Err(err).Send()
	}

	//Start server
	serv := internalserver.NewInternalServer(conf.InternalServAddress)
	go func() {
		err := serv.Start()
		if err != nil {
			log.Error().Err(err).Send()
		}
	}()

	// Stop server
	ctx, closer := context.WithTimeout(context.Background(), time.Second*2500)
	defer closer()
	go func(ctxgo context.Context) {
		<-ctxgo.Done()
		err := serv.Stop()
		if err != nil {
			log.Error().Err(err).Send()
		}
	}(ctx)
	<-ctx.Done()

	// When Done
	time.Sleep(time.Second)
	fmt.Println("Done!")
}

/*
	1. сделать internalserver пакет
		*находится в internal/web/internalserver
		?*это сервер построенный на gin фреймворке https://github.com/gin-gonic/gin
		*выносим порт запуска в env (базово 8000)
		создаем структуру и функции запуска и остановки internalserver
		делаем роутинг с 3 ручками (/helth, /readiness, /metrics) - первые 2 просто возвращают 200, последняя должна возвращать метрики контейнера
		инициализируем в сервисе - запускаем сервис - убеждаемся что ручки работают
	2. Сделать external server пакет:
		по аналогии с интерналом все, но работает на другом порту и не имеет пока ручек
		/internal/web/externalserver
	3. Пример discounts-prices/dp-api/web/
*/
