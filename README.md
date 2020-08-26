#### Program obliczający ewapotranspirację potencjalną i rzeczywistą wg. wzorów Grabarczyka i Rzekanowskiego dla drzew jabłoni, gruszy, wiśni, śliwy.
### Info 
Program obliczający ewapotranspirację potencjalną (ETp) i rzeczywistą (ETr) wg. wzorów Grabarczyka i Rzekanowskiego dla drzew jabłoni, gruszy, wiśni, śliwy.
Obliczenia na podstawie danych ze stacji IMGW w Toruniu (353180250) i współczynnika K z portalu http://www.nawadnianie.inhort.pl.

### Uruchomienie programu:
Program należy pobrać w odpowiedniej wersji dla swojego systemu operacyjnego:
Przykładowo dla Windows 64bit:

[Pobierz program](https://github.com/jtaczanowski/ewapotranspiracja-grabarczyk-rzekanowski/raw/master/program-binarki/Windows/ewapotranspiracja-grabarczyk-rzekanowski-amd64.exe)

Wszystkie dostępne wersje programu w katalogu program-binarki:

[Inne systemy operacyjne](https://github.com/jtaczanowski/ewapotranspiracja-grabarczyk-rzekanowski/tree/master/program-binarki)


Aby uruchomić program na Windows 10 należy najpierw go odblokować (kliknąć na plik prawym klawiszem myszy, wybrać "właściwości" i zaznaczyć "Odblokuj"):

<img src="/readme-obrazki/windows-odblokowanie.PNG" alt="windows-odblokowanie" height="45%" width="45%">

Następnie należy kliknąć dwa razy na ikonę programu, program się uruchomi, ściągnie i przeliczy dane, następnie zapisze do pliku csv:

![uruchomiony-program](/readme-obrazki/ewapotranspiracja-grabarczyk-rzekanowski-program.PNG)

Plik csv mozna otworzyć w arkuszu kalkulacyjnym i pracować na nim:

![praca-w-arkuszu-kalkulacyjnym](/readme-obrazki/ewapotranspiracja-grabarczyk-rzekanowski-excel.PNG)

### Wzory użyte do obliczeń:
#### Wzór Grabarczyka:
<img src="/readme-obrazki/wspolczynnik-ETp.png" alt="wspolczynnik-ETp" height="70%" width="70%">

#### Wzór Rzekanowskiego:
<img src="/readme-obrazki/wspolczynnik-ETr.png" alt="wspolczynnik-ETr" height="70%" width="70%">

#### Współcznik K (http://www.nawadnianie.inhort.pl)
<img src="/readme-obrazki/wspolczynnik-K-wg-inhort_pl.png" alt="wspolczynnik-K-wg-inhort_pl" height="70%" width="70%">

### Samodzielna kompilacja kodu programu:
Trzeba pobrać język programowania Golang z https://golang.org/dl/, następnie zainstalować, a potem w katalogu projektu uruchomić kompilację:
```
go build .
```
