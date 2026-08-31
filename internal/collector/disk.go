package collector

import (
	"syscall"
)

// Disk, disk kullanım metriklerini tutar.
type Disk struct {
	Path  string
	Total uint64
	Free  uint64
	Used  uint64
}

// CollectDisk syscall (sistem çağrısı) yaparak diskin anlık durumunu kernel'den çeker.
// Bu fonksiyon doğrudan sisteme bağlı olduğu için birim testi yapılamaz.
func CollectDisk(path string) (Disk, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return Disk{}, err
	}

	bsize := uint64(stat.Bsize)

	total := stat.Blocks * bsize
	free := stat.Bavail * bsize // Normal kullanıcının görebildiği boş alan

	// df komutunun hesaplama mantığı: Toplam - Gerçek Boş (Bfree)
	used := (stat.Blocks - stat.Bfree) * bsize

	return Disk{
		Path:  path,
		Total: total,
		Free:  free,
		Used:  used,
	}, nil
}

// UsedPercent disk kullanım yüzdesini hesaplar.
// Bu saf bir hesaplama fonksiyonu olduğu için birim testleri (unit test) yazılabilir.
func (d Disk) UsedPercent() float64 {
	// df komutunun yüzde formülü: Kullanılan / (Kullanılan + Bavail)
	denominator := d.Used + d.Free
	if denominator == 0 {
		return 0.0
	}
	return (float64(d.Used) / float64(denominator)) * 100.0
}
type DiskCollector struct {
	Path string
}

func (c DiskCollector) Name() string { return "disk" }

func (c DiskCollector) Collect() ([]Metric, error) {
	disk, err := CollectDisk(c.Path)
	if err != nil {
		return nil, err
	}
	return []Metric{
		{Name: "disk_percent", Value: disk.UsedPercent(), Unit: "%"},
	}, nil
}
