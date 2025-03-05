package utils

import (
	"calc_service/models"
	"errors"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// Генерация уникального ID
func GenerateID() string {
	return uuid.New().String()
}

func ToRPN(expression string) ([]string, error) {
	// Удаляем все пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")

	// Приоритет операторов
	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	var output []string     // Выходной массив (RPN)
	var operators []rune    // Стек для операторов и скобок
	var numberBuffer string // Буфер для сбора многоразрядных чисел

	for i, char := range expression {
		// Если символ — цифра или точка, добавляем его в буфер
		if unicode.IsDigit(char) || char == '.' {
			numberBuffer += string(char)
			// Если это последний символ, добавляем число в выходной массив
			if i == len(expression)-1 {
				output = append(output, numberBuffer)
				numberBuffer = ""
			}
			continue
		}

		// Если буфер числа не пуст, добавляем число в выходной массив
		if numberBuffer != "" {
			output = append(output, numberBuffer)
			numberBuffer = ""
		}

		// Если символ — оператор
		if char == '+' || char == '-' || char == '*' || char == '/' {
			// Переносим операторы с более высоким приоритетом из стека в выходной массив
			for len(operators) > 0 && operators[len(operators)-1] != '(' &&
				precedence[operators[len(operators)-1]] >= precedence[char] {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			// Добавляем текущий оператор в стек
			operators = append(operators, char)
			continue
		}

		// Если символ — открывающая скобка, добавляем её в стек
		if char == '(' {
			operators = append(operators, char)
			continue
		}

		// Если символ — закрывающая скобка
		if char == ')' {
			// Переносим операторы из стека в выходной массив до открывающей скобки
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			// Удаляем открывающую скобку из стека
			if len(operators) == 0 {
				return nil, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
			continue
		}

		// Если символ не распознан, возвращаем ошибку
		return nil, errors.New("invalid character in expression")
	}

	// Если буфер числа не пуст, добавляем число в выходной массив
	if numberBuffer != "" {
		output = append(output, numberBuffer)
	}

	// Переносим оставшиеся операторы из стека в выходной массив
	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func ParseExpression(expr models.Expr) ([]models.Task, error) {
	rpn, err := ToRPN(expr.Expr)
	if err != nil {
		return nil, err
	}
	log.Printf("Parsed RPN: %v", rpn)

	var stack []models.Task // Стек для хранения задач
	var tasks []models.Task // Список всех задач

	for _, token := range rpn {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			task := models.Task{
				ID:           GenerateID(), // Не создаем отдельную задачу для константы
				ExpressionID: expr.ID,
				Arg1:         num,
				Status:       "completed",
				Result:       num,
			}
			stack = append(stack, task)
			continue
		}

		if token == "+" || token == "-" || token == "*" || token == "/" {
			if len(stack) < 2 {
				return nil, errors.New("invalid expression: not enough operands for operator " + token)
			}

			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			task := models.Task{
				ID:           GenerateID(),
				ExpressionID: expr.ID,
				Arg1TaskID:   arg1.ID,
				Arg2TaskID:   arg2.ID,
				Arg1:         arg1.Result,
				Arg2:         arg2.Result,
				Operation:    token,
				Status:       "waiting",
			}
			tasks = append(tasks, task)
			stack = append(stack, task)
		}
	}

	if len(stack) != 1 {
		return nil, errors.New("invalid expression: malformed RPN")
	}

	return tasks, nil
}

// Проверка валидности выражения
func IsValidExpression(expr string) bool {
	if len(expr) == 0 {
		return false
	}

	// Проверка на допустимые символы
	for _, char := range expr {
		if !unicode.IsDigit(char) && !IsOperator(char) && char != '(' && char != ')' {
			return false
		}
	}

	// Проверка на сбалансированность скобок
	return AreParenthesesBalanced(expr)
}

// Проверка, является ли символ оператором
func IsOperator(char rune) bool {
	return strings.ContainsRune("+-*/", char)
}

// Проверка сбалансированности скобок
func AreParenthesesBalanced(expr string) bool {
	var stack []rune
	for _, char := range expr {
		if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			if len(stack) == 0 {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}
