package env

import (
	"fmt"
	"reflect"

	"github.com/Netflix/go-env"
)

// Загрузить переменные среды в структуру с тэгами "env" из os.Environ
func InitEnvironment(environment interface{}) {
	if _, err := env.UnmarshalFromEnviron(environment); err != nil {
		t := reflect.TypeOf(environment)
		v := reflect.ValueOf(environment)
		errMsg := ""
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			errMsg += fmt.Sprintf("%s,type=%s. Value provided=%v\n", field.Tag.Get("env"), field.Type, v.Field(i).Interface())
		}
		panic(fmt.Errorf("%w. Loaded environment: \n%s", err, errMsg))
	}
}
