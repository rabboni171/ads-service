package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CreateAdRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_ad_request_count",
		Help: "Количество запросов на создание объявления",
	})
	CreateAdOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_ad_ok_count",
		Help: "Количество успешно созданных объявлений",
	})
	CreateAdError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_ad_error_count",
		Help: "Количестов ошибок при создании объявления",
	})

	GetAdRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_ad_request_count",
		Help: "Количество запросов на получение единичного объявления",
	})
	GetAdOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_ad_ok_count",
		Help: "Количество успешно полученных единичных объявлений",
	})
	GetAdError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_ad_error_count",
		Help: "Количестов ошибок при получении единичного объявления",
	})

	GetAdsRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_ads_request_count",
		Help: "Количество запросов на получение страниц объявлений",
	})
	GetAdsOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_ads_ok_count",
		Help: "Количество успешно полученных страниц объявлений",
	})
	GetAdsError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_ads_error_count",
		Help: "Количестов ошибок при получении страниц объявлений",
	})
)

