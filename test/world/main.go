package main

import (
	"fmt"
	"game-engine/core"
	"game-engine/test"
)

func main() {
	world := core.NewWorld()
	world.RegisterLevel(core.NewLevel[test.Level1]())

	// 홍콩에서 특정 이벤트 발생 시 신규 instance level을 생성 요청
	// static level과 dynmaic level을 분리
	// dynamic level의 경우 ready start가 있어야함
	world.CreateLevel("TestLevel")

	//go func() {
	//	var name string
	//	fmt.Scanln(&name)
	//	world.Stop()
	//}()

	<-world.Run()
	fmt.Println("close world")

}
