package main

import (
	"fmt"
	"os"

	"github.com/NinnjA254/httpfromtcp/internal/request"
	"github.com/NinnjA254/httpfromtcp/internal/server"
)

type Apartment struct {
	Name                 string
	rentPerMonth         int
	ConstructionComplete bool
	Location             string
	Amenities            []string
}

type filter string

const (
	ALL        filter = "all"
	FINISHED   filter = "finished"
	UNFINISHED filter = "unfinished"
)

func getApartments(f filter) []Apartment {
	finishedUnits := []Apartment{
		{"based apartments", 6000, true, "Kilifi town", []string{"wifi", "fresh water"}},
		{"Baobab beach flats", 8000, true, "Baobab beach", []string{"wifi", "fresh water", "beach view"}},
		{"Barizi apartments", 10000, true, "Lamu", []string{"wifi", "fresh water", "air conditioning"}},
	}
	unfinishedUnits := []Apartment{
		{"creek apartments", 0, false, "Kilifi creek", []string{"wifi", "balcony", "fresh water", "beach view"}},
		{"Maisha Mnarani apartments", 0, false, "Mnarani", []string{"wifi", "fresh water", "air conditioning"}},
		{"Baobab life apartments", 0, false, "Baobab beach", []string{"wifi", "fresh water", "beach view"}},
	}
	switch f {
	case FINISHED:
		return finishedUnits
	case UNFINISHED:
		return unfinishedUnits
	case ALL:
		allUnits := []Apartment{}
		allUnits = append(allUnits, finishedUnits...)
		allUnits = append(allUnits, unfinishedUnits...)
		return allUnits
	default:
		return nil
	}
}
func main() {
	s := server.InitServer()
	s.Route("/finishedunits", func(req *request.Request, res *server.Response) {
		fmt.Println("/finished")
		finishedUnits := getApartments(FINISHED)
		res.SendJSON(finishedUnits)
	})
	s.Route("/allunits", func(req *request.Request, res *server.Response) {
		fmt.Println("/allunits")
		allUnits := getApartments(ALL)
		res.SendJSON(allUnits)
	})
	s.Route("/", func(req *request.Request, res *server.Response) {
		fmt.Println("/")
		res.SendHTML(
			"<h1>Tsavo apartments</h1>" +
				"<p>site under construction</p>")
	})
	s.Route("/favicon.ico", func(req *request.Request, res *server.Response) {
		fmt.Println("/favicon.ico")
		f, err := os.Open("./cmd/blog/img.jpg")
		if err != nil {
			res.SendText("image sending error")
			return
		}
		defer f.Close()
		fileInfo, err := f.Stat()
		if err != nil {
			res.SendText("image sending error")
			return
		}
		imgSize := fileInfo.Size()
		imgData := make([]byte, imgSize)
		n, err := f.Read(imgData)
		if err != nil {
			res.SendText("image sending error")
			return
		}
		res.SendImage(imgData[:n])
	})

	s.Listen(42069)
}
