package collection_helper

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

type Person struct {
	name string
	id   int
	age  int
}

func TestMap(t *testing.T) {
	people := make([]Person, 0, 10)
	for i := 10; i > 0; i-- {
		people = append(people, Person{
			id:   i,
			name: "p" + strconv.Itoa(i),
			age:  10 + (i % 2),
		})
	}
	fmt.Printf("%+v\n", people)
	// toIndexMap
	indexMap := ToIndexMap[Person](&people)
	fmt.Printf("%+v\n", indexMap)
	// toItselfMap key 为 id value 为 people
	idPersonMap := ToItselfMap[Person, int](&people, func(person Person) int {
		return person.id
	})
	fmt.Printf("%+v\n", idPersonMap)
	// toSimpleMap key 为 id value 为 people
	idPersonMap2 := ToSimpleMap[Person, int, Person](&people, func(person Person) int {
		return person.id
	}, func(person Person) Person {
		return person
	})
	fmt.Printf("%+v\n", idPersonMap2)
	//toGroup groupKey 为 age values 为 []Person
	agePersonGroup := ToItselfGroup[Person, int](&people, func(person Person) int {
		return person.age
	})
	fmt.Printf("%+v\n", agePersonGroup)

	// 按照 person的age 分map,最后需要的map的value是name,如果person的age重复了,走重复处理逻辑
	toMap, err := ToMap[Person, int, string](&people, func(person Person) int {
		// 处理key 按照什么来分map
		return person.age
	}, func(person Person) string {
		//处理value,最后的值是什么
		return person.name
	}, func(x, y string) (string, error) {
		////重复处理
		//fmt.Println(x, y)
		//// 后者替换前者
		//return y, nil
		return "", errors.New("dulicate!")
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", toMap)

}
