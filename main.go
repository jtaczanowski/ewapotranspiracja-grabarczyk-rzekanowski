// Program obliczający ewapotranspirację potencjalną i rzeczywistą wg. wzorów Grabarczyka i Rzekanowskiego dla drzew jabłoni, gruszy, wiśni, śliwy.
// Obliczenia na podstawie danych ze stacji IMGW w Toruniu i współczynnika K z http://www.nawadnianie.inhort.pl.
package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const info = `
############################## Informacja o programie ##############################
# Program obliczający ewapotranspirację potencjalną (ETp) i rzeczywistą (ETr)      #
# wg. wzorów Grabarczyka i Rzekanowskiego dla drzew jabłoni, gruszy, wiśni, śliwy. #
# Obliczenia na podstawie danych ze stacji IMGW w Toruniu (353180250)              #
# i współczynnika K z http://www.nawadnianie.inhort.pl                             #
# Jan Taczanowski                                                                  #
# Poznań 2020                                                                      #
####################################################################################
`

const stacjaMeteorologiczna = "353180250" // Stacja metorelogiczna Torun

var miesiace = map[string]string{
	"01": "styczen",
	"02": "luty",
	"03": "marzec",
	"04": "kwiecien",
	"05": "maj",
	"06": "czerwiec",
	"07": "lipiec",
	"08": "sierpien",
	"09": "wrzesien",
	"10": "pazdziernik",
	"11": "listopad",
	"12": "grudzien",
}
var paraNasycona = map[int]float64{
	0:  0.6168,
	1:  0.6566,
	2:  0.7054,
	3:  0.7575,
	4:  0.8129,
	5:  0.8719,
	6:  0.9346,
	7:  1.0012,
	8:  1.0721,
	9:  1.1473,
	10: 1.2271,
	11: 1.3118,
	12: 1.4015,
	13: 1.4967,
	14: 1.5974,
	15: 1.7041,
	16: 1.8170,
	17: 1.9364,
	18: 2.063,
	19: 2.196,
	20: 2.337,
	21: 2.485,
	22: 2.642,
	23: 2.808,
	24: 2.982,
	25: 3.166,
}

// wspolczynikK wzięty z http://www.nawadnianie.inhort.pl
var wspolczynikK = map[string]map[string]float64{
	"jablon": {"04": 0.5, "05": 0.75, "06": 1.1, "07": 1.2, "08": 1.2, "09": 1.15},
	"grusza": {"04": 0.45, "05": 0.75, "06": 1.05, "07": 1.15, "08": 1.15, "09": 1.1},
	"wisnia": {"04": 0.45, "05": 0.75, "06": 1, "07": 1.1, "08": 1.1, "09": 0.9},
	"sliwa":  {"04": 0.45, "05": 0.75, "06": 1.1, "07": 1.2, "08": 1.15, "09": 1.15},
}

func main() {
	fmt.Println(info + "\n")
	filePath, err := os.Getwd()
	if err != nil {
		fmt.Println("Blad przy pobieraniu aktualnej sciezki w której uruchamiany jest program", err)
		waitForExit()
	}
	var imgwKontent string
	urls := prepareUrls()
	fmt.Printf("--------------------------- Pobieram dane z dane.imgw.pl ---------------------------\n")
	for _, url := range urls {
		respHTTP := getUrl(url)
		imgwKontent = imgwKontent + unZipHTTPresp(respHTTP)
	}
	fmt.Printf("\n----------------------------- Przetwarzam dane ------------------------------------\n")
	timeStart := time.Now()
	dane := obliczDane(imgwKontent)
	timeDuration := time.Now().Sub(timeStart)
	fmt.Println("Przetworzono w czasie: ", timeDuration)
	fmt.Printf("\n------------------------------ Zapisuję do pliku ----------------------------------\n")
	pathToWriteFile := path.Join(filePath, "ewapotranspiracja-grabarczyk-rzekanowski-"+time.Now().Format("2006-01-02_15-04-05")+".csv")
	err = ioutil.WriteFile(pathToWriteFile, []byte(dane), 0644)
	if err != nil {
		fmt.Print("Błąd podczas zapisywania wyników obliczeń do pliku csv: ", err)
		waitForExit()
	}
	fmt.Printf("Zapisano wyniki do pliku: %s \n\n", pathToWriteFile)
	waitForExit()
}

func prepareUrls() []string {
	var urlList []string
	urlList = append(urlList, "https://dane.imgw.pl/data/dane_pomiarowo_obserwacyjne/dane_meteorologiczne/miesieczne/synop/1986_1990/1986_1990_m_s.zip")
	urlList = append(urlList, "https://dane.imgw.pl/data/dane_pomiarowo_obserwacyjne/dane_meteorologiczne/miesieczne/synop/1991_1995/1991_1995_m_s.zip")
	urlList = append(urlList, "https://dane.imgw.pl/data/dane_pomiarowo_obserwacyjne/dane_meteorologiczne/miesieczne/synop/1996_2000/1996_2000_m_s.zip")
	for i := 2001; i <= time.Now().Year(); i++ {
		urlList = append(urlList, fmt.Sprintf("https://dane.imgw.pl/data/dane_pomiarowo_obserwacyjne/dane_meteorologiczne/miesieczne/synop/%v/%v_m_s.zip", i, i))
	}
	return urlList
}

func getUrl(url string) []byte {
	fmt.Printf("Pobieram: %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Bład w trakcie przy pobieraniu danych z imgw", err)
		waitForExit()
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Bład przy odczycie danych z imgw", err)
		waitForExit()
	}
	return body
}

func unZipHTTPresp(body []byte) string {
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		fmt.Println("Bład przy rozpakowywaniu plików zip z imgw: ", err)
		waitForExit()
	}
	for _, zipFile := range zipReader.File {
		stationFileNamePattern := "s_m_t_"
		match, err := regexp.MatchString(stationFileNamePattern, zipFile.Name)
		if err != nil {
			fmt.Println("Bład przy szukaniu patternu s_m_t_ w nazwach pliów zimgw: ", err)
			waitForExit()
		}
		if match {
			fmt.Println("Wczytuję plik:", zipFile.Name)
			f, err := zipFile.Open()
			if err != nil {
				log.Println(err)
				continue
			}
			defer f.Close()
			unzippedFileBytes, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Println("Bład przy odczycie rozpakowanych danych z imgw", err)
				waitForExit()
			}
			return string(unzippedFileBytes)
		}
	}
	return ""
}

func obliczDane(data string) string {
	scanner := bufio.NewScanner(strings.NewReader(data))
	dane := "rok;miesiac;srednia temp;suma opadow;ETp;ETr Jablon;ETr Grusza;ETr Wisnia;ETr Sliwa;Srednie Cisnienie Pary Wodnej;Wzor\n"
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), ",")
		match, err := regexp.MatchString(stacjaMeteorologiczna, scanner.Text())
		if err != nil {
			fmt.Println("Bład przy szukaniu patternu stacji meteoroligcznej "+stacjaMeteorologiczna+" w danych ", err)
			continue
		}
		if match {
			var rok string
			var miesiac string
			var sredniaTemp float64
			var nocnyOpad float64
			var dziennyOpad float64
			var cisnienie float64
			rok = record[2]
			miesiac = record[3]
			if temp, err := strconv.ParseFloat(record[8], 64); err == nil {
				sredniaTemp = temp
			} else {
				fmt.Println("Blad parsowania float64 ze stringa z pliku imgw: ", err.Error())
				fmt.Println("Kontynuuję pracę")
				continue

			}
			if v, err := strconv.ParseFloat(record[10], 64); err == nil {
				cisnienie = v
			} else {
				fmt.Println("Blad parsowania float64 ze stringa z pliku imgw: ", err.Error())
				fmt.Println("Kontynuuję pracę")
				continue

			}
			if v, err := strconv.ParseFloat(record[18], 64); err == nil {
				nocnyOpad = v
			} else {
				fmt.Println("Blad parsowania float64 ze stringa z pliku imgw: ", err.Error())
				fmt.Println("Kontynuuję pracę")
				continue

			}
			if v, err := strconv.ParseFloat(record[20], 64); err == nil {
				dziennyOpad = v
			} else {
				fmt.Println("Blad parsowania float64 ze stringa z pliku imgw: ", err.Error())
				fmt.Println("Kontynuuję pracę")
				continue

			}
			sumaOpad := nocnyOpad + dziennyOpad
			if sredniaTemp < 0.0 {
				sredniaTemp = 0
			}
			if miesiac == "\"04\"" || miesiac == "\"05\"" || miesiac == "\"06\"" || miesiac == "\"07\"" || miesiac == "\"08\"" || miesiac == "\"09\"" {
				ETp := 0.32 * ((30 * ((paraNasycona[int(sredniaTemp)] * 10) - cisnienie)) + (10 * sredniaTemp))
				ETrJablon := ETp * wspolczynikK["jablon"][strings.ReplaceAll(miesiac, "\"", "")]
				ETrGrusza := ETp * wspolczynikK["grusza"][strings.ReplaceAll(miesiac, "\"", "")]
				ETrWisnia := ETp * wspolczynikK["wisnia"][strings.ReplaceAll(miesiac, "\"", "")]
				ETrSliwa := ETp * wspolczynikK["sliwa"][strings.ReplaceAll(miesiac, "\"", "")]
				Wzor := fmt.Sprintf("0.32 * ((30*((%v*10) - %v)) + (10 * %v))", paraNasycona[int(sredniaTemp)], cisnienie, sredniaTemp)
				dane = dane + rok + ";" + miesiac + ";" + strings.ReplaceAll(fmt.Sprintf("%v", sredniaTemp), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", sumaOpad), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", ETp), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", ETrJablon), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", ETrGrusza), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", ETrWisnia), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", ETrSliwa), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", cisnienie), ".", ",") + ";" + strings.ReplaceAll(fmt.Sprintf("%v", Wzor), ".", ",") + "\n"
			}
		}
	}
	return dane
}

func waitForExit() {
	go func() {
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-termChan
		os.Exit(0)
	}()
	fmt.Println("Aby wyjść z programu kliknij ctrl+c")
	select {}
}
