## 2026-08-29 — Gün 1

**Alt komut vs flag:** `sysprobe status` seçildi. Sistem araçlarının ortak dili bu ve ileride komut eklemek kolay. Maliyeti: Go'nun flag paketi alt komutu desteklemiyor, FlagSet ile elle yönetmek gerekiyor.

**main/run ayrımı:** os.Exit defer'ları atladığı için Exit tek noktada tutuldu. Yan fayda: run() test edilebilir.

**status şimdilik exit 3 dönüyor:** Bir izleme aracı "bilmiyorum" durumunda "iyi" dememeli.

## Gün 3

**testdata klasörü:** Go toolchain bu ismi build'den hariç tutuyor. Gerçek `/proc/meminfo` kopyası buraya alındı; testler artık makineden bağımsız ve CI'da da aynı sonucu veriyor.

**Tablo bazlı test:** Senaryolar veri olarak tutuluyor, tek döngüyle çalışıyor. Yeni senaryo eklemek bir satır — sürtünme düşük olunca gerçekten ekliyorum.

**Float karşılaştırma:** Yüzdeler epsilon toleransıyla karşılaştırılıyor. Doğrudan `!=` kullanmak, float aritmetiğinin kesin olmaması nedeniyle rastgele patlamalara yol açar.

**Testlerin bulduğu buglar:** Tablo bazlı testler yazılırken boş girdilerdeki `MemTotal` kontrolsüzlüğü ve sıfır swap durumlarındaki potansiyel hatalar önceden yakalanıp güvenlik korumaları eklendi.
