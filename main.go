package packages_metadata

import (
	"fmt"
	"io"
	"strings"
)

// Интерфейс библиотеку получения метаданных используемых пакетов в приложении
type PackagesMetadata interface {
	// Получение версии для конкретного пакета
	GetVersion(name string) string
	// Получение списка сервисов и их версий
	GetMap() map[string]string
}

// Таблица метаданных пакетов
type packages struct {
	// Таблица метаданных пакетов, где ключ - алиас, значение - структура с метаданными {версия и название библиотеки}
	m map[string]packageMeta
}

// Метаданные пакета
type packageMeta struct {
	// Название пакета
	name string
	// Версия пакета
	version string
}

func newPackage(name, version string) packageMeta {
	return packageMeta{
		name:    name,
		version: version,
	}
}

func (p packages) GetVersion(name string) string {
	if pac, ok := p.m[name]; ok {
		return pac.version
	}
	return ""
}

func (p packages) GetMap() map[string]string {
	result := make(map[string]string)
	for alias, pack := range p.m {
		result[alias] = pack.version
	}
	return result
}

/*
	Получение списка интересующих пакетов с их версиями
*/
func GetPackagesMetadata(packagesNameVersion, packagesAliases map[string]string) PackagesMetadata {
	p := new(packages)
	p.m = make(map[string]packageMeta)
	for alias, packageName := range packagesAliases {
		version, ok := packagesNameVersion[packageName]
		if !ok {
			continue
		}
		p.m[alias] = newPackage(packageName, version)
	}
	return p
}

/*
	Чтение списка всех пакетов используемых в приложении из файла
	Формат данных в файле должен соответствовать выводу команды "go list -m all"
	т.е.: package/name v1.0.0
*/
func GetPackagesVersion(reader io.Reader) (map[string]string, error) {
	fileContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read packages versions file failed: %w", err)
	}
	// Ответ в виде мапы название пакета и версия
	result := make(map[string]string)
	depsList := strings.Split(string(fileContent), "\n")

	for _, d := range depsList {
		pac := strings.Split(strings.Trim(d, " "), " ")
		// Оставляем только те пакеты, которые имеют версии
		if len(pac) != 2 {
			continue
		}
		result[pac[0]] = pac[1]
	}
	return result, nil
}
