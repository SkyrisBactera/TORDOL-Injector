package main

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-humble/locstor"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jQuery = jquery.NewJQuery
var chance = 10
var truthInjectees []string
var dareInjectees []string
var truthInjecteesSave []string
var dareInjecteesSave []string
var alreadyDone = false

const (
	enter = 13
)

func start(truthInjects, dareInjects []string) {
	truthInjectees = make([]string, len(truthInjects))
	copy(truthInjectees, truthInjects)
	truthInjecteesSave = make([]string, len(truthInjectees))
	copy(truthInjecteesSave, truthInjectees)
	dareInjectees = make([]string, len(dareInjects))
	copy(dareInjectees, dareInjects)
	dareInjecteesSave = make([]string, len(dareInjectees))
	copy(dareInjecteesSave, dareInjectees)
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Save("truthOriginal", truthInjecteesSave); err != nil {
		println("Couldn't save truthOriginal!")
	}
	if err := store.Save("dareOriginal", dareInjecteesSave); err != nil {
		println("Couldn't save dareOriginal!")
	}
	load()
	js.Global.Call("addEventListener", "keyup", func(event *js.Object) {
		keycode := event.Get("keyCode").Int()
		if keycode == enter && !alreadyDone {
			alreadyDone = true
			gucciGang()
		}
	}, false)
}

func gucciGang() {
	if determine() {
		if strings.Contains(js.Global.Get("location").Get("href").String(), "truth") && len(truthInjectees) > 0 {
			println("Truth")
			setMessage(truthInjectees[0])
			truthInjectees = append(truthInjectees[:0], truthInjectees[0+1:]...)
		} else if strings.Contains(js.Global.Get("location").Get("href").String(), "dare") && len(dareInjectees) > 0 {
			println("Dare")
			setMessage(dareInjectees[0])
			dareInjectees = append(dareInjectees[:0], dareInjectees[0+1:]...)
		} else {
			println("Out of truths/dares")
		}
	}
	save()
	os.Exit(0)
}

func main() {
	js.Global.Set("ti", map[string]interface{}{
		"start": start,
	})
}

func determine() (returnVal bool) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randNum := r1.Intn(chance)
	if randNum == 0 {
		println("Showing \"injectee\"!")
		returnVal = true
	} else {
		println("Chances dictated " + strconv.Itoa(randNum) + " where 0 was needed. Chance is " + strconv.Itoa(chance))
		returnVal = false
		chance--
	}
	return
}

func setMessage(message string) {
	jQuery(".mainContent").SetText(message)
}

func save() {
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Delete("chance"); err != nil {
		println("Couldn't delete chance!")
	}
	if err := store.Save("chance", chance); err != nil {
		println("Couldn't save chance!")
	}
	if err := store.Save("truthInjectees", truthInjectees); err != nil {
		println("Couldn't save truthInjectees!")
	}
	if err := store.Save("dareInjectees", dareInjectees); err != nil {
		println("Couldn't save dareInjectees!")
	}
}

func load() {
	var truthOriginal []string
	var dareOriginal []string
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Find("chance", &chance); err != nil {
		println("Couldn't load chance!", err)
	}
	if err := store.Find("truthInjectees", &truthInjectees); err != nil {
		println("Couldn't load truthInjectees!", err)
	}
	if err := store.Find("dareInjectees", &dareInjectees); err != nil {
		println("Couldn't load truthInjectees!", err)
	}
	if err := store.Find("truthOriginal", &truthOriginal); err != nil {
		println("Couldn't load truthOriginal!", err)
	} else {
		if !testEq(truthOriginal, truthInjecteesSave) {
			println(truthOriginal, truthInjecteesSave)
			println("Injectees has changed! Resetting truthOriginal and resetting all values to defaults")
			reset()
		}
	}
	if err := store.Find("dareOriginal", &dareOriginal); err != nil {
		println("Couldn't load dareOriginal!", err)
	} else {
		if !testEq(dareOriginal, dareInjecteesSave) {
			println(dareOriginal, dareInjecteesSave)
			println("Injectees has changed! Resetting dareOriginal and resetting all values to defaults")
			reset()
		}
	}
}

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func reset() {
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Delete("truthOriginal"); err != nil {
		println("Couldn't delete truthOriginal!")
	}
	if err := store.Delete("dareOriginal"); err != nil {
		println("Couldn't delete dareOriginal!")
	}
	if err := store.Save("truthOriginal", truthInjecteesSave); err != nil {
		println("Couldn't save truthOriginal!")
	}
	if err := store.Save("dareOriginal", dareInjecteesSave); err != nil {
		println("Couldn't save dareOriginal!")
	}
	chance = 10
	save()
}
