package main

import(
	"github.com/albrow/zoom"
	"fmt"
	//"strconv"
	//"strings"
	"time"
)

type Person struct{
	Name string
	Age  int
	zoom.DefaultData
} 

func main(){
	conf:= &zoom.Configuration{
		Address:"localhost:59999",
		Network:"tcp",
	}
	zoom.Init(conf)
	if err:= zoom.Register(&Person{});err!=nil{
		fmt.Println("regiser Error")
	}
	starttime:=time.Now()

	p := &Person{Name: "Alice", Age: 27}
	//for i:=0;i<1000000;i++{
	//	tmp:= []string{"name",strconv.Itoa(i)}
	//	zoom.Save(&Person{Name:strings.Join(tmp,"|"),Age:i+2})
	//}

	result,_ := zoom.FindById("Person",p.Id)
	person, ok := result.(*Person)
	if !ok{
		fmt.Println("Error Person type")
	}
	fmt.Println(result,person.Name,person.Age)
	var persons  []*Person
	if err:=zoom.NewQuery("Person").Scan(&persons);err!=nil{
		panic(err)
	}
	for _,p:= range persons{
		//p.Age =p.Age+20000
		//zoom.Save(p)
		fmt.Println(p.Age,p.Name)
	}
	fmt.Println(len(persons))
	//for i,rs := range results{
	//	if ps, ok := rs.(*Person); !ok {
	//	}
	//	fmt.Println(ps.Age, ps.Name)
	//}
	fmt.Println(time.Now().Sub(starttime))
	defer zoom.Close()
}
