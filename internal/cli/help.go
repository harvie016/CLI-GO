package cli

import "fmt"

const helpText = `mws — утилита для управления профилями.

Профиль хранится в YAML-файле в текущей директории: профиль test лежит
в файле test.yaml рядом с местом запуска команды.

Использование:
  mws <команда> [аргументы]

Команды:
  profile create --name=<name> --user=<user> --project=<project>
        создать профиль; все три флага обязательны
  profile get --name=<name>
        показать данные профиля
  profile list
        вывести все профили текущей директории
  profile delete --name=<name>
        удалить профиль
  help
        показать эту справку

Пример:
  mws profile create --name=test --user=example --project=new-project
`

// Help печатает справку по всем командам.
func Help() {
	fmt.Print(helpText)
}
