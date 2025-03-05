package models

type Expression struct {
	ID     string  `json:"id"`
	Expr   string  `json:"-"` // Исходное выражение (не включается в JSON)
	Status string  `json:"status"`
	Result float64 `json:"result,omitempty"`
}

type Task struct {
	ID           string  `json:"id"`
	ExpressionID string  `json:"expression_id"`
	Arg1TaskID   string  `json:"arg1_task_id,omitempty"` // Ссылка на первую задачу (для операций)
	Arg2TaskID   string  `json:"arg2_task_id,omitempty"` // Ссылка на вторую задачу (для операций)
	Arg1         float64 `json:"arg1,omitempty"`         // Число (для констант)
	Arg2         float64 `json:"arg2,omitempty"`         // Число (для констант, не используется)
	Operation    string  `json:"operation,omitempty"`    // Оператор (+, -, *, /)
	Status       string  `json:"status"`                 // Статус задачи
	Result       float64 `json:"result,omitempty"`       // Результат задачи
}

type Expr struct {
	ID   string `json:"id"`
	Expr string `json:"expression"`
}
