package validation

import (
	"fmt"
	"github.com/fatih/color"
)

func CheckFever() {
	var suhu int

	fmt.Print("Masukkan suhu badan Anda: ")
	fmt.Scan(&suhu)

	if suhu > 37 {
		color.Red("Kamu demam!")
	} else {
		color.Green("Kamu sehat!")
	}
}