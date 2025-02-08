package utils

const (
	maxDigit = 9 // Максимальное значение цифры
)

// ValidateLuhn проверяет номер заказа по алгоритму Луна
func ValidateLuhn(number string) bool {
	sum := 0
	isSecond := false

	// Проходим по цифрам справа налево
	for i := len(number) - 1; i >= 0; i-- {
		d := int(number[i] - '0')

		if isSecond {
			d *= 2
			if d > maxDigit {
				d -= maxDigit
			}
		}

		sum += d
		isSecond = !isSecond
	}

	return sum%10 == 0
}
