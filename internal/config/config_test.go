package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Olmayan bir dosya verelim, varsayılan değerler dönmeli
	cfg, err := Load("olmayan_dosya.json")
	if err != nil {
		t.Fatalf("Beklenmeyen hata: %v", err)
	}
	if *cfg.Thresholds.CPUWarning != 80.0 {
		t.Errorf("Beklenen 80.0, alinan %f", *cfg.Thresholds.CPUWarning)
	}
}

func TestLoad_Validation(t *testing.T) {
	// Geçersiz bir JSON dosyası oluşturalım (Warning > Critical)
	content := `{"thresholds": {"cpu_warning": 95.0, "cpu_critical": 90.0}}`
	
	tmpfile, _ := os.CreateTemp("", "config_*.json")
	defer os.Remove(tmpfile.Name())
	
	tmpfile.Write([]byte(content))
	tmpfile.Close()

	_, err := Load(tmpfile.Name())
	if err == nil {
		t.Fatal("cpu_warning > cpu_critical iken hata fırlatılmalıydı, ancak fırlatılmadı")
	}
}
