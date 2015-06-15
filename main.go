package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
        "path/filepath"
        "errors"
	"regexp"
	"strconv"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <GPIO for fan controller> <temperature limit>")
	}
	gpio, _ := strconv.Atoi(os.Args[1])
        tempLimit, _ := strconv.ParseFloat(os.Args[2], 32)
        var limit = float32(tempLimit) 

	var pin = rpio.Pin(gpio)
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	pin.Output()

	temperature, err := readTemperature()
	if err != nil {
		log.Fatal(err)
	}
        
        log.Printf("Limit:%f, Now:%f\n", limit, temperature)
        if temperature >= limit {
                log.Println("Cooling fan is turned ON.")
                pin.High()
        } else {
                log.Println("Cooling fan is turned OFF.")
                pin.Low()
        }
}

func readTemperature() (float32, error) {
	thermDevice, err := getThermDevice("/sys/bus/w1/devices")
	if err != nil {
		return 0, err
	}
	dat, err := ioutil.ReadFile(thermDevice)
	if err != nil {
		return 0, err
	}
	s := string(dat)

	pattern := regexp.MustCompile("t=(\\d+)")
	matches := pattern.FindStringSubmatch(s)
	temperature, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return 0, err
	}

	return float32(temperature / 1000), nil
}

func getThermDevice(baseDir string) (devPath string, err error) {
        subDirInfos, err := ioutil.ReadDir(baseDir)
        if err != nil {
                return "", err
        }

        for _, fileInfo := range subDirInfos {
                var subDir = filepath.Join(baseDir, (fileInfo).Name())
                dirInfo, err := os.Stat(subDir)
                if err == nil && dirInfo.IsDir() {
                        var devfile = filepath.Join(subDir, "w1_slave")
                        _, err := os.Stat(devfile)
                        if !os.IsNotExist(err) {
                                return devfile, nil
                        }
                }
        }
        return "", errors.New("w1_slave is not found.")
}

