package actioninfo

import (
	"fmt"
	"log"
)

// DataParser описывает поведение парсеров данных об активностях.
type DataParser interface {
	// Parse разбирает строку с данными об активности.
	Parse(data string) error
	// ActionInfo формирует строку с информацией об активности.
	ActionInfo() (string, error)
}

// Info перебирает строки датасета, парсит каждую из них и выводит
// информацию об активности. Ошибки логируются, обработка продолжается.
func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Println(err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(info)
	}
}
