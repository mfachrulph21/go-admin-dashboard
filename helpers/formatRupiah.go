package helpers

import "github.com/yudapc/go-rupiah"

func FormatUang(nominal int) string {
	floatValue := float64(nominal)
	formatRupiah := rupiah.FormatRupiah(floatValue)

	return formatRupiah
}
