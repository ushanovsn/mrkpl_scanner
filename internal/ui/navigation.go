package ui

import (
	"mrkpl_scanner/internal/options"
)

// Main navigation menu structure with data returns
func GetNavigationMenu() []options.NaviMenu {
	return []options.NaviMenu{
		{
			ItmName:   "Главная",
			ItmLink:   "/",
			ItmIsMenu: false,
			ItmMenu:   nil,
		},
		{
			ItmName:   "Статус",
			ItmLink:   "/status",
			ItmIsMenu: false,
			ItmMenu:   nil,
		},
		{
			ItmName:   "Конфигурация",
			ItmLink:   "/config",
			ItmIsMenu: true,
			ItmMenu: []options.NaviDropMenu{
				{
					ItmName: "Основные параметры",
					ItmLink: "/config_base",
				},
				{
					ItmName: "Wildberries",
					ItmLink: "/config_wb",
				},
				{
					ItmName: "Ozon",
					ItmLink: "/config_ozon",
				},
				{
					ItmName: "Yandex Market",
					ItmLink: "/config_yandex_market",
				},
				{
					ItmName: "Avito",
					ItmLink: "/config_avito",
				},
				{
					ItmName: "Уведомления",
					ItmLink: "/config_notify",
				},
			},
		},
		{
			ItmName:   "Параметры задач",
			ItmLink:   "/task_param",
			ItmIsMenu: true,
			ItmMenu: []options.NaviDropMenu{
				{
					ItmName: "Сканирование",
					ItmLink: "/task_param_scan",
				},
				{
					ItmName: "Мониторинг",
					ItmLink: "/task_param_monitor",
				},
				{
					ItmName: "Поиск",
					ItmLink: "/task_param_search",
				},
			},
		},
		{
			ItmName:   "Лог",
			ItmLink:   "/log",
			ItmIsMenu: false,
			ItmMenu:   nil,
		},
	}
}
