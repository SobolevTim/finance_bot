package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expr string) (int, error) {
	expr = strings.ReplaceAll(expr, ",", ".")

	var currentNumber string
	var sumFloat, temp float64
	var lastOperator byte = '+'

	for i := 0; i < len(expr); i++ {
		char := expr[i]

		if unicode.IsDigit(rune(char)) || char == '.' {
			// Собирать цифры и десятичные точки в строку
			currentNumber += string(char)
		}

		// Обработка операторов или завершение строки
		if char == '+' || char == '-' || char == '*' || char == '/' || i == len(expr)-1 {
			if currentNumber != "" {
				number, err := strconv.ParseFloat(currentNumber, 64)
				if err != nil {
					return 0, err
				}

				// Выполняем операцию с предыдущим числом, с учётом приоритета
				switch lastOperator {
				case '+':
					sumFloat += temp
					temp = number
				case '-':
					sumFloat += temp
					temp = -number
				case '*':
					temp *= number
				case '/':
					if number == 0 {
						return 0, fmt.Errorf("деление на ноль")
					}
					temp /= number
				}

				// Очищаем строку для следующего числа
				currentNumber = ""
			}

			// Сохраняем текущий оператор
			lastOperator = char
		}
	}

	// Добавляем последнее значение `temp` к результату
	sumFloat += temp

	// Возвращаем результат в формате int (умножаем на 100 для учёта копеек)
	return int(sumFloat * 100), nil
}
