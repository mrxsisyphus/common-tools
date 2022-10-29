package convert_helper

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TestUnMarshalToAny(t *testing.T) {
	p := Person{
		Id:   1,
		Name: "p-1",
	}
	marshal, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	}
	marshalStr := BytesToStr(marshal)
	fmt.Println(marshalStr)
	var p2 Person
	p3, err := JsonUnMarshalToAny[Person](StrToBytes(marshalStr), &p2)
	if err != nil {
		panic(err)
	}
	//err = json.Unmarshal(StrToBytes(marshalStr), &p2)
	//if err != nil {
	//	panic(err)
	//}
	fmt.Printf("%p\n", &p2)
	fmt.Printf("%p\n", p3)

}
