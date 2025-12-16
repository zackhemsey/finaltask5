package actioninfo

import (
	"fmt"
)

type DataParser interface {
	Parse(data string) error
	ActionInfo() (string, error)
}

func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			fmt.Printf("Ошибка парсинга: %v\n", err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			fmt.Printf("Ошибка формирования информации: %v\n", err)
			continue
		}

		fmt.Print(info)

	}
}
