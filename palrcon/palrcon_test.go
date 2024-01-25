package palrcon

import (
	"fmt"
	"testing"
)

var endpoint = "endpointHere"
var password = "adminPasswordHere"

func TestGetPlayers(t *testing.T) {

	pal := NewPalRCON(endpoint, password)

	data, err := pal.GetPlayers()
	if err != nil {
		fmt.Println(err.Error())
	}

	for k := range data {
		fmt.Println(data[k].Name, data[k].PlayerUID, data[k].SteamID)
	}
}

func TestInfo(t *testing.T) {
	pal := NewPalRCON(endpoint, password)

	info, err := pal.Info()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(info)
}

func TestSave(t *testing.T) {
	pal := NewPalRCON(endpoint, password)

	err := pal.Save()
	if err != nil {
		fmt.Println(err.Error())
	}

}
