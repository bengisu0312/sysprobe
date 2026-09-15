package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Pointer kullanıyoruz çünkü 0 değeri ile "değer girilmedi" (nil) ayrımını yapmalıyız.
type Thresholds struct {
	CPUWarning   *float64 `json:"cpu_warning"`
	CPUCritical  *float64 `json:"cpu_critical"`
	DiskWarning  *float64 `json:"disk_warning"`
	DiskCritical *float64 `json:"disk_critical"`
}

type Config struct {
	Thresholds Thresholds `json:"thresholds"`
}

// DefaultConfig gömülü (embedded) varsayılan değerleri döndürür
func DefaultConfig() *Config {
	cpuW, cpuC := 80.0, 90.0
	diskW, diskC := 80.0, 90.0
	return &Config{
		Thresholds: Thresholds{
			CPUWarning:   &cpuW,
			CPUCritical:  &cpuC,
			DiskWarning:  &diskW,
			DiskCritical: &diskC,
		},
	}
}

// Load, arama hiyerarşisine (Flag -> Env -> /etc -> Default) göre config dosyasını yükler.
func Load(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	// 1. Flag ile değer gelmediyse Environment Variable'a bak (SYSPROBE_CONFIG)
	path := configPath
	if path == "" {
		path = os.Getenv("SYSPROBE_CONFIG")
	}

	// 2. Env de boşsa, standart Linux konfigürasyon dizinini (/etc) kullan
	if path == "" {
		path = "/etc/sysprobe/config.json"
	}

	// 3. Dosyayı oku
	data, err := os.ReadFile(path)
	if err != nil {
		// Dosya yoksa sorun değil, program varsayılanlarla devam etmeli
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("config dosyasi okunamadi: %w", err)
	}

	// 4. Dosya varsa, varsayılanların üzerine yaz (Merge)
	// Pointer kullandığımız için JSON'da olmayan alanlar nil gelir ve cfg'deki mevcut varsayılanları bozmaz!
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("json parse hatasi: %w", err)
	}

	// 5. Doğrulama (Validation)
	t := cfg.Thresholds
	if *t.CPUWarning >= *t.CPUCritical {
		return nil, fmt.Errorf("hata: cpu_warning (%.1f), cpu_critical'dan (%.1f) kucuk olmalidir", *t.CPUWarning, *t.CPUCritical)
	}
	if *t.DiskWarning >= *t.DiskCritical {
		return nil, fmt.Errorf("hata: disk_warning (%.1f), disk_critical'dan (%.1f) kucuk olmalidir", *t.DiskWarning, *t.DiskCritical)
	}

	return cfg, nil
}
