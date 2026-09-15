## Gün 6
**Interface neden gerekti:** main.go'da dört ayrı blok vardı, her collector farklı tip döndürdüğü için tek listede toplanamıyorlardı. Interface bir tip değil, davranış tarifi — "Name() ve Collect() metodları olan her şey".
**Neden şimdi, önce değil:** Dört collector'ı yazmadan ortak şeklin ne olduğunu bilemezdim. Soyutlamayı somut kodu görmeden tasarlamak, iki kez baştan yazmak demek. Interface'i kullanan taraf tanımlar.
**Collect() neden []Metric döndürüyor:** Tek float64 yetmezdi — bellek hem RAM hem swap veriyor, load üç değer veriyor, disk birden fazla mount için çalışabilir.
## Config Yönetimi ve Kararlar
**Neden JSON (YAML veya TOML değil)?**
Go'nun dahili paketlerinde (stdlib) JSON desteği bulunurken, YAML ve TOML için dış bağımlılık (third-party package) indirmek gerekiyordu. "Sıfır bağımlılık" mimarisi hedefimiz doğrultusunda JSON tercih edildi. 
**Pointer Kullanımı:** JSON içerisindeki `0` değeri ile "değer girilmediği" durumunu (nil) ayırt edebilmek ve eksik ayarlarda gömülü varsayılanları (defaults) ezmemek için config yapılarında pointer `*float64` kullanıldı.
**Arama Sırası:** `--config` flag'i -> `SYSPROBE_CONFIG` env değişkeni -> `/etc/sysprobe/config.json` -> Gömülü varsayılanlar.
## Kurallar ve Çıkış Kodları (Nagios Standardı)
**Neden Worst-of Aggregation?** İzleme sistemlerinde yanlış negatif (false negative), yanlış pozitiften (false positive) çok daha tehlikelidir. Disk %100 doluysa, CPU'nun %2'de çalışıyor olmasının bir önemi yoktur ve sistem tehlikededir. Bu yüzden genel durum, her zaman değerlendirilen metrikler arasındaki en kötü duruma eşitlenir.
**os.Exit ve defer tuzağı:** Go'da `os.Exit()` çağrıldığında program anında kapanır ve `defer` blokları (dosya kapatma vb.) çalışmaz. Bu yüzden `run()` fonksiyonu doğrudan çıkış yapmak yerine int (exit code) döndürecek şekilde tasarlandı ve `os.Exit()` yalnızca `main()` içinde bir kez çağrıldı.

