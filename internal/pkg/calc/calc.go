package calc

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Calculate — основная функция для вычисления выражения.
// Поддерживаются операции: +, -, *, /, ^ и процент (интерпретируется как X ± (X*Y/100)).
func Calculate(expression string) (float64, error) {
	// Убираем пробелы и нормализуем запись
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(expression, "\t", "")
	expression = strings.ReplaceAll(expression, "\n", "")
	expression = strings.ReplaceAll(expression, "\r", "")
	expression = strings.ToLower(expression)
	expression = strings.ReplaceAll(expression, ",", ".")

	if !isValidExpression(expression) {
		return 0, errors.New("недопустимые символы в выражении")
	}

	// Проверяем соответствие скобок
	openCount := strings.Count(expression, "(")
	closeCount := strings.Count(expression, ")")
	if openCount > closeCount {
		expression += strings.Repeat(")", openCount-closeCount)
	} else if closeCount > openCount {
		return 0, errors.New("несоответствие скобок")
	}

	// Преобразуем процентные выражения
	expression = transformPercentages(expression)

	// Преобразуем инфиксное выражение в постфиксное
	postfix, err := toPostfix(expression)
	if err != nil {
		return 0, err
	}

	// Вычисляем постфиксное выражение
	result, err := evaluatePostfix(postfix)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// transformPercentages преобразует выражения вида:
//
//	<Операнд>[+|–]<число>%
//
// в выражения вида:
//
//	<Операнд> [+|–] ((<Операнд>*<число>)/100)
//
// То есть процент всегда считается от всей части слева от оператора на данном уровне.
// Функция работает в два этапа:
// 1. Рекурсивно обрабатываются подвыражения в скобках.
// 2. На текущем уровне последовательно заменяются вхождения вида [+\-]<число>%.
// При каждом проходе проверяется, изменилось ли выражение – если нет, цикл завершается.
func transformPercentages(expression string) string {
	// 1. Обрабатываем подвыражения в скобках, пока что‑ли есть замены.
	reParen := regexp.MustCompile(`\([^()]*\)`)
	for {
		newExpr := reParen.ReplaceAllStringFunc(expression, func(s string) string {
			inner := s[1 : len(s)-1]
			return "(" + transformPercentages(inner) + ")"
		})
		if newExpr == expression {
			break
		}
		expression = newExpr
	}

	// 2. Обрабатываем процентные шаблоны на текущем уровне.
	rePerc := regexp.MustCompile(`([+\-])(\d+(?:\.\d+)?)%`)
	for {
		if !rePerc.MatchString(expression) {
			break
		}
		newExpr := expression
		// Находим все совпадения и обрабатываем их справа налево
		matches := rePerc.FindAllStringSubmatchIndex(expression, -1)
		if len(matches) == 0 {
			break
		}
		for i := len(matches) - 1; i >= 0; i-- {
			m := matches[i]
			fullStart, fullEnd := m[0], m[1]
			op := expression[m[2]:m[3]]     // оператор: + или -
			numStr := expression[m[4]:m[5]] // число перед %
			base := expression[:fullStart]  // всё, что слева от вхождения
			if strings.TrimSpace(base) == "" {
				// Если база пуста, удаляем это вхождение, чтобы избежать зацикливания
				newExpr = newExpr[:fullStart] + newExpr[fullEnd:]
				continue
			}
			// Формируем замену: <op>((<base>)*(<num>)/100)
			replacement := fmt.Sprintf("%s((%s)*(%s)/100)", op, base, numStr)
			newExpr = newExpr[:fullStart] + replacement + newExpr[fullEnd:]
		}
		if newExpr == expression {
			break
		}
		expression = newExpr
	}
	return expression
}

// isValidExpression проверяет, что в выражении присутствуют только допустимые символы.
func isValidExpression(expression string) bool {
	for _, r := range expression {
		if !unicode.IsDigit(r) && !strings.ContainsRune("+-*/()^.%", r) {
			return false
		}
	}
	return true
}

// toPostfix преобразует инфиксное выражение в постфиксную запись (обратную польскую нотацию).
func toPostfix(expression string) ([]string, error) {
	var result []string
	var stack []rune

	for i := 0; i < len(expression); i++ {
		ch := rune(expression[i])

		if unicode.IsDigit(ch) || ch == '.' {
			start := i
			for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || expression[i+1] == '.') {
				i++
			}
			result = append(result, expression[start:i+1])
		} else if ch == '(' {
			stack = append(stack, ch)
		} else if ch == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return nil, errors.New("несоответствие скобок")
			}
			stack = stack[:len(stack)-1]
		} else if isOperator(ch) {
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(ch) {
				if ch == '^' && stack[len(stack)-1] == '^' {
					break // ^ — правоассоциативный оператор
				}
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, ch)
		} else {
			return nil, errors.New("неизвестный символ в выражении")
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, errors.New("несоответствие скобок")
		}
		result = append(result, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return result, nil
}

// evaluatePostfix вычисляет значение выражения, заданного в постфиксной записи.
func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if len(token) == 1 && isOperator(rune(token[0])) {
			if len(stack) < 2 {
				return 0, errors.New("недостаточно операндов для операции")
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("деление на ноль")
				}
				result = a / b
			case "^":
				result = math.Pow(a, b)
			}
			stack = append(stack, result)
		} else {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("ошибка в вычислениях")
	}
	return stack[0], nil
}

// isOperator возвращает true, если символ является оператором.
func isOperator(ch rune) bool {
	return strings.ContainsRune("+-*/^", ch)
}

// precedence определяет приоритет оператора.
func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '^':
		return 3
	default:
		return 0
	}
}

// FormatNumber форматирует число в строку с округлением до 5 знаков после запятой.
// Если число целое, возвращается без десятичной части.
func FormatNumber(amount float64) string {
	if amount == float64(int(amount)) {
		return fmt.Sprintf("%d", int(amount))
	}
	rounded := math.Round(amount*1e5) / 1e5
	return strconv.FormatFloat(rounded, 'f', -1, 64)
}
