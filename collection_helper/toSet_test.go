package collection_helper

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSet(t *testing.T) {
	people := make([]Person, 0, 10)
	for i := 10; i > 0; i-- {
		people = append(people, Person{
			id:   i,
			name: "p" + strconv.Itoa(i),
			age:  10 + (i % 2),
		})
	}
	fmt.Printf("%+v\n", people)

	//Set 直接把people去重(people对象直接hash)
	peopleSet := ToMapSet[Person, Person](&people, func(person Person) Person {
		return person
	})
	fmt.Printf("%+v\n", peopleSet)
	//类似
	peopleSet2 := ToItSelfMapSet[Person](&people)
	fmt.Printf("%+v\n", peopleSet2)

	// 用户按照age维度组mapSet
	peopleSet3 := ToMapSet[Person, int](&people, func(person Person) int {
		// 按照age 维度
		return person.age
	})
	fmt.Printf("%+v\n", peopleSet3)

}
