## Gün 6
**Interface neden gerekti:** main.go'da dört ayrı blok vardı, her collector farklı tip döndürdüğü için tek listede toplanamıyorlardı. Interface bir tip değil, davranış tarifi — "Name() ve Collect() metodları olan her şey".
**Neden şimdi, önce değil:** Dört collector'ı yazmadan ortak şeklin ne olduğunu bilemezdim. Soyutlamayı somut kodu görmeden tasarlamak, iki kez baştan yazmak demek. Interface'i kullanan taraf tanımlar.
**Collect() neden []Metric döndürüyor:** Tek float64 yetmezdi — bellek hem RAM hem swap veriyor, load üç değer veriyor, disk birden fazla mount için çalışabilir.
