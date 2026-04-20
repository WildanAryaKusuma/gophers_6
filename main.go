package main

import (
	"belajar-golang/inventory"
	// "fmt"
	// "github.com/fatih/color"
	// "belajar-golang/validation"
	// "belajar-golang/cashier"
) 

type Mobil struct {
	Ban string
	Pintu int
	Ac bool

}

func main() {
	// ** Session 1 **
	// hello := "Hello World!"

	// fmt.Println(hello)
	// color.Red("Hellow World in Red!")

	// validation.CheckFever()
  // cashier.Payment()

	// * Session 2 *

	// ? declare array normal way
	// var buah [4]string

	// buah[0] = "Mangga"
	// buah[1] = "Apel"
	// buah[2] = "Kiwuy"
	// buah[3] = "Jeruk"
	// fmt.Print(buah)

	// ? sliceeees
	// var buah []string

	// buah = append(buah, "Mangga")
	// buah = append(buah, "Apel")
	// buah = append(buah, "Kiwuy")
	// buah = append(buah, "Jeruk")

	// fmt.Print(buah)

	// ? structsss
	// honda := Mobil {
	// 	Ban: "Dunlop", 
	// 	Pintu: 4, 
	// 	Ac: true, 

	// }

	// fmt.Print(honda)

	// ? maps

	// ? Looping
	// traditional for loop
	// for i := 0; i < 3; i++ {
	// 	fmt.Println("Halo", i)
	// }

	// for range loop
	// for i, data := range buah {
	// 	fmt.Println("Buah ke", i, "adalah", data);
	// }

	// ! TASK SESSION 2: INVENTORY MANAGER
	inventory.TokoKelontong()
}