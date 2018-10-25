package main

import (
	"fmt"
	"time"

	r "github.com/go-vgo/robotgo"
)

//var cBitmapsDontGCme = make([]interface{}, 0)

// http://world-of-talesworth.wikia.com/wiki/Equipment

func main() {
	bmp, free := loadBitmaps("bag_button.png", "bag_full.png", "upgrade.png", "vendor_trash.png", "vendor_trash_faded.png", "done.png", "done2.png",
		"next.png", "hp_bar.png", "chest_simple.png", "chest_blue.png", "arena1-4.png", "arena5-14.png", "arena15-24.png", "arena25-34.png",
		"arena35-44.png", "arena45-54.png", "arena55-60.png", "zone_easy.png", "zone_faceroll.png", "town.png", "map.png", "eat.png", "eat2.png",
		"money.png", "computer.png", "hungry.png", "goto_work.png", "work.png", "sleep.png", "truck.png")
	defer free()

	for {
		eatIfHungry(bmp)
		workIfAvailable(bmp)
		nextZoneIfEasy(bmp)
		clickChests(bmp)
		clickTruck(bmp)

		//vendorAndEquip(bmp)
		//clickNextDoneButtons(bmp)
		//clickEnemies(bmp)

		time.Sleep(500 * time.Millisecond)
	}
}

func clickTruck(bmp map[string]r.Bitmap) {
	for clickBitmap(bmp["truck.png"], 10, -3) {
		for i := 0; i < 20; i++ {
			r.MouseClick("left", true)
		}
	}
}

func clickEnemies(bmp map[string]r.Bitmap) {
	for clickBitmap(bmp["hp_bar.png"], 10, -3) {
		for i := 0; i < 10; i++ {
			r.MouseClick("left", true)
		}
	}
}

func clickChests(bmp map[string]r.Bitmap) {
	clickBitmap(bmp["chest_blue.png"])
	clickBitmap(bmp["chest_simple.png"])
}

func doStuffComputerRoomThenExit(bmp map[string]r.Bitmap) {
	if clickBitmap(bmp["eat.png"]) {
		time.Sleep(3 * time.Second)
	} else if clickBitmap(bmp["eat2.png"]) {
		time.Sleep(3 * time.Second)
	}
	if clickBitmap(bmp["goto_work.png"]) {
		time.Sleep(3 * time.Second)
	}
	if clickBitmap(bmp["sleep.png"]) {
		time.Sleep(3 * time.Second)
	}
	for i := 0; i < 10; i++ {
		if clickBitmap(bmp["computer.png"]) {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("could not exit computer room")
}

func eatIfHungry(bmp map[string]r.Bitmap) {
	if clickBitmap(bmp["hungry.png"]) && clickBitmap(bmp["money.png"]) {
		time.Sleep(3 * time.Second)
		doStuffComputerRoomThenExit(bmp)
	}
}

func vendorAndEquip(bmp map[string]r.Bitmap) {
	if clickBitmap(bmp["bag_full.png"], 20, -30) {
		time.Sleep(1000 * time.Millisecond)
		for clickBitmap(bmp["upgrade.png"], 0, 0, 1) {
			time.Sleep(500 * time.Millisecond)
		}
		if clickBitmap(bmp["vendor_trash.png"], 20, 5) || clickBitmap(bmp["vendor_trash_faded.png"], 20, 5) {
			time.Sleep(500 * time.Millisecond)
		}
		r.KeyTap("escape")
	}
}

func nextZoneIfEasy(bmp map[string]r.Bitmap) {
	if clickBitmap(bmp["zone_easy.png"]) || clickBitmap(bmp["zone_faceroll.png"]) {
		if clickBitmap(bmp["town.png"]) {
			time.Sleep(3 * time.Second)
			r.MouseClick("left", true)
			clickBitmap(bmp["arena55-60.png"])
			clickBitmap(bmp["arena45-54.png"])
			clickBitmap(bmp["arena35-44.png"])
			clickBitmap(bmp["arena25-34.png"])
			clickBitmap(bmp["arena15-24.png"])
			clickBitmap(bmp["arena5-14.png"])
			clickBitmap(bmp["arena1-4.png"])
		} else {
			fmt.Println("town.png not found")
		}
	}
}

func workIfAvailable(bmp map[string]r.Bitmap) {
	if clickBitmap(bmp["work.png"]) {
		time.Sleep(3 * time.Second)
		doStuffComputerRoomThenExit(bmp)
	}
}

func clickNextDoneButtons(bmp map[string]r.Bitmap) {
	if clickBitmap(bmp["done.png"]) || clickBitmap(bmp["done2.png"]) {
	}
	if clickBitmap(bmp["next.png"]) {
		r.MouseClick("left", true)
	}
}

func clickBitmap(bmp r.Bitmap, offsets ...int) bool {
	//fmt.Printf("\n%+v", bmp)
	xOffset := 0
	yOffset := 0
	if len(offsets) == 2 {
		xOffset = offsets[0]
		yOffset = offsets[1]
	}
	doubleclick := false
	if len(offsets) == 3 && offsets[2] != 0 {
		doubleclick = true
	}
	fx, fy := r.FindBitmap(r.ToCBitmap(bmp))
	if fx != -1 && fy != -1 {
		r.MoveMouse(fx+xOffset, fy+yOffset)
		r.MouseClick("left", doubleclick)
		return true
	}
	return false
}

func loadBitmaps(files ...string) (bitmaps map[string]r.Bitmap, free func()) {
	freeFuncs := make([]func(), 0)
	bitmaps = make(map[string]r.Bitmap)
	for _, f := range files {
		bitmap, freeFunc := loadBitmap(f)
		bitmaps[f] = bitmap
		freeFuncs = append(freeFuncs, freeFunc)
	}
	free = func() {
		for key := range freeFuncs {
			freeFuncs[key]()
		}
	}
	return bitmaps, free
}

func loadBitmap(file string) (bitmap r.Bitmap, free func()) {
	cBitmap := r.OpenBitmap(file)
	//cBitmapsDontGCme = append(cBitmapsDontGCme, cBitmap)
	bitmap = r.ToBitmap(cBitmap)
	free = func() {
		r.FreeBitmap(cBitmap)
	}
	return bitmap, free
}
