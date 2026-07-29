package textparser

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseAmount bóc tách số tiền từ một chuỗi văn bản (ví dụ: "Ăn sáng 50k", "Mua đồ 150.000", "Điện nước 1.5 triệu")
func ParseAmount(text string) (float64, error) {
	// Chuyển toàn bộ về chữ thường để dễ xử lý
	text = strings.ToLower(text)

	// Các mẫu regex phổ biến:
	// 1. Số có chữ "k" (ví dụ: 50k, 100k, 50 k) -> nhân 1,000
	// 2. Số có chữ "cành" (ví dụ: 50 cành) -> nhân 1,000
	// 3. Số có chữ "triệu" hoặc "tr" (ví dụ: 1.5 triệu, 2 tr) -> nhân 1,000,000
	// 4. Số thường, có hoặc không có dấu chấm phẩy (ví dụ: 150.000, 150000)

	// Xử lý chữ "k" hoặc "cành"
	reK := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(k|cành)`)
	if match := reK.FindStringSubmatch(text); match != nil {
		val, _ := strconv.ParseFloat(match[1], 64)
		return val * 1000, nil
	}

	// Xử lý chữ "triệu" hoặc "tr"
	reTr := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(triệu|tr)`)
	if match := reTr.FindStringSubmatch(text); match != nil {
		val, _ := strconv.ParseFloat(match[1], 64)
		return val * 1000000, nil
	}

	// Xử lý số thuần túy (VD: 150.000 hoặc 150,000)
	// Bắt các cụm số có chứa dấu chấm/phẩy, hoặc cụm số liên tiếp
	reNum := regexp.MustCompile(`\d{1,3}(?:[.,]\d{3})+|\d+`)
	matches := reNum.FindAllString(text, -1)
	if len(matches) > 0 {
		// Lấy số cuối cùng xuất hiện trong chuỗi (thường số tiền nằm ở cuối)
		lastMatch := matches[len(matches)-1]
		// Xóa các dấu chấm phẩy
		cleanNum := strings.ReplaceAll(lastMatch, ".", "")
		cleanNum = strings.ReplaceAll(cleanNum, ",", "")
		val, _ := strconv.ParseFloat(cleanNum, 64)
		return val, nil
	}

	return 0, nil
}
