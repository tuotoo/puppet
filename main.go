package main

import (
	"github.com/zach-klippenstein/goadb"
	"fmt"
	"log"
	"flag"
	"os"
	"strconv"
	"git.spiritframe.com/tuotoo/utils"
	"time"
	"io"
	"image/png"
	"github.com/astaxie/beego"
	"strings"
	"image"
	"math"
)

var (
	port = flag.Int("p", adb.AdbPort, "")

	client *adb.Adb
)

func main() {
	flag.Parse()

	var err error
	client, err = adb.NewWithConfig(adb.ServerConfig{
		Port: *port,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting server…")
	client.StartServer()

	serverVersion, err := client.ServerVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server version:", serverVersion)

	devices, err := client.ListDevices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Devices:")
	for _, device := range devices {
		fmt.Printf("\t%+v\n", *device)
	}
	firstdevice := devices[0]

	device := client.Device(adb.DeviceWithSerial(firstdevice.Serial))
	for {
		run(device)
	}
}

const centerx, centery = 750, 1308

func run(device *adb.Device) {
	time.Sleep(time.Second)
	dis := Positioning(device)
	beego.Debug(dis)
	press(device,int(float64(dis)*2.8))
}

func Positioning(device *adb.Device) int{
	output, err := device.RunCommand("screencap", []string{"-p"}...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 0
	}
	fi, err := os.Create("pic" + strconv.FormatInt(time.Now().Unix(), 10) + ".png")
	if !utils.Check(err) {
		return 0
	}
	defer fi.Close()
	fi.WriteString(output)
	return Image(strings.NewReader(output))
}

func press(device *adb.Device, ti int) {
	output, err := device.RunCommand("input", []string{"swipe", "720", "2100","720", "2100", strconv.FormatInt(int64(ti), 10)}...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}
	log.Println(output)
	time.Sleep(time.Duration(ti) * time.Millisecond)
}

func Image(r io.Reader) int {
	img, err := png.Decode(r)
	if !utils.Check(err) {
		return 0
	}
	return centery-gety(img)
}

func gety(img image.Image)int{
	l := left(img)
	r := right(img)
	if l < r{
		return l
	}
	return r
}

func left(img image.Image)int{
	var c float64
	for x:=0;x<centerx;x++{
		y := centery - int(float64(centerx-x) /math.Sqrt(3))
		color := img.At(x,y)
		r,g,b,_ := color.RGBA()
		cc := float64(r>>8+g>>8+b>>8 )/3
		if c!=0 {
			if math.Abs(c-cc) >10{
				beego.Debug(r>>8,g>>8,b>>8,c,cc,x,y,int(float64(centerx-x) /math.Sqrt(3)))
				return y
			}
		}
		c = cc
	}
	return 2560
}


func right(img image.Image)int{
	var c float64
	for x:=1440;x>centerx;x--{
		y := centery - int(float64(x-centerx) /math.Sqrt(3))
		color := img.At(x,y)
		r,g,b,_ := color.RGBA()
		cc := float64(r>>8+g>>8+b>>8 )/3
		if c!=0 {
			if math.Abs(c-cc) >5{
				beego.Debug(r>>8,g>>8,b>>8,c,cc,x,y,int(float64(centerx-x) /math.Sqrt(3)))
				return y
			}
		}
		c = cc
	}
	return 2560
}