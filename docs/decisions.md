## 2026-08-29 — Gün 1

**Alt komut vs flag:** `sysprobe status` seçildi. Sistem araçlarının ortak dili bu ve ileride komut eklemek kolay. Maliyeti: Go'nun flag paketi alt komutu desteklemiyor, FlagSet ile elle yönetmek gerekiyor.

**main/run ayrımı:** os.Exit defer'ları atladığı için Exit tek noktada tutuldu. Yan fayda: run() test edilebilir.

**status şimdilik exit 3 dönüyor:** Bir izleme aracı "bilmiyorum" durumunda "iyi" dememeli.
