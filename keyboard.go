package main

import (
	"fmt"
	"github.com/micmonay/keybd_event"
	"log"
	"runtime"
	"time"
)

// Странно работает и не постоянна, нам нужна более быстрая работа виртуальной клавиатуры
func main() {
	fun()
	// Here, the program will generate "ABAB" as if they were pressed on the keyboard.
}

func fun() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatalln(err)
	}

	start := time.Now()
	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(40 * time.Millisecond)
	}
	fmt.Println(time.Since(start))

	// Select keys to be pressed
	kb.SetKeys(keybd_event.VK_A, keybd_event.VK_B)
	fmt.Println(time.Since(start))

	// Set shift to be pressed
	//kb.HasSHIFT(true)

	// Press the selected keys
	err = kb.Launching()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(time.Since(start))

	// Or you can use Press and Release
	kb.Press()
	time.Sleep(time.Millisecond)
	kb.Release()
	fmt.Println(time.Since(start))
}
