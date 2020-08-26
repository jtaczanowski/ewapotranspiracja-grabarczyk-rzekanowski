#### Program obliczający ewapotranspirację potencjalną i rzeczywistą wg. wzorów Grabarczyka i Rzekanowskiego dla drzew jabłoni, gruszy, wiśni, śliwy.
##### Info 
Program obliczający ewapotranspirację potencjalną (ETp) i rzeczywistą (ETr) wg. wzorów Grabarczyka i Rzekanowskiego dla drzew jabłoni, gruszy, wiśni, śliwy.
Obliczenia na podstawie danych ze stacji IMGW w Toruniu (353180250) i współczynnika K z http://www.nawadnianie.inhort.pl.

Wzór Grabarczyka:
![wspolczynnik-ETp.png](/readme-obrazki/wspolczynnik-ETp.png)
Wzór Rzekanowskiego:
![wspolczynnik-ETr.png](/readme-obrazki/wspolczynnik-ETr.png)
Współcznik K (http://www.nawadnianie.inhort.pl)
![wspolczynnik-K-wg-inhort_pl.png](/readme-obrazki/wspolczynnik-K-wg-inhort_pl.png)

##### Uruchomienie programu
Program nalezy pobrać 

Aby uruchomić program na Windows 10 nalezy najpierw go odblokowac:

![windows-odblokowanie](/readme-obrazki/windows-odblokowanie.PNG)

Następnie nalezy kliknąć dwa razy na ikone programu, program sie uruchomi, sciągnie i przeliczy dane, następnie zapisze do pliku csvs:

![uruchomiony-program](/ewapotranspiracja-grabarczyk-rzekanowski-program.PNG)

Plik csv mozna otworzyć w arkuszu kalkulacyjnym i pracować na nim:

![praca-w-arkuszu-kalkulacyjnym](/ewapotranspiracja-grabarczyk-rzekanowski-program.PNG)

##### Samodzielna kompilacja kodu programu:
Nalezy ściągnąć język programowania Golang z https://golang.org/dl/, następnie zainstalować, a potem w katalogu projektu uruchomić kompilację:
```
go build .
```
