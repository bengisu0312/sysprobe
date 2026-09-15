## Gün 6
**Interface neden gerekti:** main.go'da dört ayrı blok vardı, her collector farklı tip döndürdüğü için tek listede toplanamıyorlardı. Interface bir tip değil, davranış tarifi — "Name() ve Collect() metodları olan her şey".
**Neden şimdi, önce değil:** Dört collector'ı yazmadan ortak şeklin ne olduğunu bilemezdim. Soyutlamayı somut kodu görmeden tasarlamak, iki kez baştan yazmak demek. Interface'i kullanan taraf tanımlar.
**Collect() neden []Metric döndürüyor:** Tek float64 yetmezdi — bellek hem RAM hem swap veriyor, load üç değer veriyor, disk birden fazla mount için çalışabilir.
## Config Yönetimi ve Kararlar
**Neden JSON (YAML veya TOML değil)?**
Go'nun dahili paketlerinde (stdlib) JSON desteği bulunurken, YAML ve TOML için dış bağımlılık (third-party package) indirmek gerekiyordu. "Sıfır bağımlılık" mimarisi hedefimiz doğrultusunda JSON tercih edildi. 
**Pointer Kullanımı:** JSON içerisindeki `0` değeri ile "değer girilmediği" durumunu (nil) ayırt edebilmek ve eksik ayarlarda gömülü varsayılanları (defaults) ezmemek için config yapılarında pointer `*float64` kullanıldı.
**Arama Sırası:** `--config` flag'i -> `SYSPROBE_CONFIG` env değişkeni -> `/etc/sysprobe/config.json` -> Gömülü varsayılanlar.
