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
			ItmName:   "Уведомления",
			ItmLink:   "/notify",
			ItmIsMenu: false,
			ItmMenu:   nil,
		},
		{
			ItmName:   "Конфигурация",
			ItmLink:   "/config",
			ItmIsMenu: true,
			ItmMenu: []options.NaviDropMenu{
				{
					ItmName: "Wildberries",
					ItmLink: "/parser_config_wb",
				},
				{
					ItmName: "Ozon",
					ItmLink: "/parser_config_ozon",
				},
				{
					ItmName: "Yandex Market",
					ItmLink: "/parser_config_yandex_market",
				},
				{
					ItmName: "Avito",
					ItmLink: "/parser_config_avito",
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
