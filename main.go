package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hanwen/go-mtpfs/mtp"
)

func GetDevice() (*mtp.Device, error) {

	dev, err := mtp.SelectDevice("GoPro")
	if err != nil {
		return nil, err
	}

	dev.Configure()
	return dev, nil
}

func main() {
	fmt.Println("--------------------")
	dev, err := GetDevice()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer dev.Close()
	fmt.Println(dev.ID())

	sids := mtp.Uint32Array{}
	err = dev.GetStorageIDs(&sids)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, id := range sids.Values {
		var info mtp.StorageInfo
		err = dev.GetStorageInfo(id, &info)
		if err != nil {
			fmt.Printf(err.Error())
		} else {
			displayHandles(dev, id)
		}
	}
}

func displayHandles(dev *mtp.Device, id uint32) {
	hs := mtp.Uint32Array{}
	err := dev.GetObjectHandles(id, 0x0, 0x0, &hs)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, handle := range hs.Values {
		var oi mtp.ObjectInfo
		dev.GetObjectInfo(handle, &oi)
		if strings.Contains(oi.Filename, "MP4") {
			fmt.Printf("%d:%s\n", oi.ObjectFormat, oi.Filename)
			writeFile(dev, handle, oi.Filename)
		}
	}
}

func writeFile(dev *mtp.Device, handle uint32, name string) {
	fs, err := os.Create(name)
	if err != nil {
		return
	}
	defer fs.Close()
	writer := bufio.NewWriter(fs)
	err = dev.GetObject(handle, writer)
	if err != nil {
		fmt.Println(err.Error())
	}
}
