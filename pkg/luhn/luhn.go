package luhn

func Valid(n string) bool {

	total := 0
	// флаг, чтобы трекать 2 цифру справа
	isSecondDigit := false

	// проходимся по цифрам в обратном порядке
	for i := len(n) - 1; i >= 0; i-- {
		// конвертируем ASCII в инт
		digit := int(n[i] - '0')

		if isSecondDigit {
			// умножаем цифру на 2 за каждую вторую цифру справа
			digit *= 2
			if digit > 9 {
				// если получилось двузначное число, то вычитаем 9
				digit -= 9
			}
		}

		// добавляем в общую сумму
		total += digit

		// возвращаем флаг в начальное состояние
		isSecondDigit = !isSecondDigit
	}

	// возвращаем кратно ли число 10
	return total%10 == 0
}
