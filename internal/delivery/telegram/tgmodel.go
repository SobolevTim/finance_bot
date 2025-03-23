package telegram

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/SobolevTim/finance_bot/internal/pkg/calc"
)

var (
	StatusBudget = "budget" // Статус установки бюджета "budget"
)

type Transaction struct {
	Date        time.Time // Дата транзакции
	Category    string    // Категория транзакции
	Amount      string    // Сумма транзакции
	Description string    // Примечание для траты
	Result      float64   // Результат подсчета
	Error       error     // Ошибка
}

func ParseInput(input string) Transaction {
	t := Transaction{}
	input = strings.TrimSpace(input)

	// Шаг 1: Извлечение даты (если присутствует)
	// Дата может быть в формате "01.02" или "01.02.2015"
	dateRe := regexp.MustCompile(`^(\d{1,2}\.\d{1,2}(?:\.\d{2,4})?)\s+`)
	if dateMatch := dateRe.FindStringSubmatch(input); len(dateMatch) > 0 {
		d := dateMatch[1]
		// Если указаны только день и месяц, добавляем текущий год
		if len(d) < 8 {
			d = fmt.Sprintf("%s.%d", d, time.Now().Year())
		}
		date, err := time.Parse("02.01.2006", d)
		if err != nil {
			t.Error = fmt.Errorf("invalid date format")
			return t
		}
		t.Date = date
		input = strings.TrimSpace(input[len(dateMatch[0]):])
	} else {
		// Если дата не указана, используем сегодняшний день
		t.Date = time.Now()
	}

	// Шаг 2: Извлечение математического выражения и примечания
	// Математическое выражение может состоять только из цифр, точек, знаков +, -, *, /, %, ^ и пробелов.
	// Всё, что идёт после этого выражения, считается примечанием.
	exprRe := regexp.MustCompile(`^([0-9+\-*/%^. ]+)(.*)$`)
	exprMatch := exprRe.FindStringSubmatch(input)
	if len(exprMatch) < 2 {
		t.Error = fmt.Errorf("invalid input format")
		return t
	}

	t.Amount = strings.TrimSpace(exprMatch[1])
	t.Description = strings.TrimSpace(exprMatch[2])

	// Проверка наличия математического выражения
	if t.Amount == "" {
		t.Error = fmt.Errorf("empty expression")
		return t
	}

	// Подготовка выражения для вычисления
	exprStr := prepareExpression(t.Amount)
	result, err := calc.Calculate(exprStr)
	if err != nil {
		t.Error = fmt.Errorf("expression error: %v", err)
		return t
	}
	t.Result = result

	return t
}

func prepareExpression(expr string) string {
	// Удаляем все пробелы
	expr = strings.ReplaceAll(expr, " ", "")
	// Заменяем ^ на ** для возведения в степень (если calc.Calculate ожидает такой синтаксис)
	// expr = strings.ReplaceAll(expr, "^", "**")
	// Если внутри чисел встречаются лишние пробелы (например, "10 000"), удаляем их
	return regexp.MustCompile(`(\d)\s+(\d)`).ReplaceAllString(expr, "$1$2")
}
