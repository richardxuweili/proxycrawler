package proxy

import (
	"math/rand"
	"reflect"
	"time"
)

func shuffleSlice(slice interface{}) {
	if rSlice := reflect.ValueOf(slice); rSlice.Kind() == reflect.Slice {
		rander := rand.New(rand.NewSource(time.Now().Unix()))
		for i := rSlice.Len() - 1; i >= 0; i-- {
			j := rander.Intn(i + 1)
			tmpValue := reflect.ValueOf(rSlice.Index(i).Interface())
			rSlice.Index(i).Set(rSlice.Index(j))
			rSlice.Index(j).Set(tmpValue)
		}
	}
}
