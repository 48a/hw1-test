package domain

const percentScale = 100

// MissingPolicy — стратегия допустимости пропусков при расчёте итога.
type MissingPolicy interface {
	// Check проверяет допустимость текущего количества оценок.
	Check(graded, criteria int) error
}

// MissingAsZero разрешает пропуски, оставляя их вклад равным нулю.
type MissingAsZero struct{}

// Check разрешает текущий итог при любом количестве оценок.
func (MissingAsZero) Check(_, _ int) error { return nil }

// RequireAll разрешает окончательный итог только при полном наборе оценок.
type RequireAll struct{}

// Check проверяет полноту оценивания.
func (RequireAll) Check(graded, criteria int) error {
	if graded != criteria {
		return ErrIncomplete
	}
	return nil
}

// Total — вычисленный итог, который не хранится в агрегате.
type Total struct {
	// Percent — взвешенный процент без округления.
	Percent float64
	// Graded — число выставленных оценок.
	Graded int
	// Criteria — число критериев в рубрикаторе.
	Criteria int
}

// Calculator — доменный сервис взвешенного итога со сменяемой стратегией.
// Нулевое значение считает пропуски нулём.
type Calculator struct {
	policy MissingPolicy
}

// NewCalculator создаёт сервис с заданной стратегией учёта пропусков.
func NewCalculator(policy MissingPolicy) Calculator {
	return Calculator{policy: policy}
}

// Calculate вычисляет итог по текущим оценкам и всем критериям курса.
func (s Calculator) Calculate(course *Course, studentID string) (Total, error) {
	if _, err := course.Student(studentID); err != nil {
		return Total{}, err
	}
	criteria := course.Criteria()
	if len(criteria) == 0 {
		return Total{}, ErrEmptyRubric
	}
	var largestWeight float64
	for _, criterion := range criteria {
		largestWeight = max(largestWeight, criterion.weight)
	}
	var weighted, weights float64
	result := Total{Criteria: len(criteria)}
	for _, criterion := range criteria {
		// Масштабирование предотвращает переполнение суммы больших весов.
		weight := criterion.weight / largestWeight
		weights += weight
		if evaluation, exists := course.Evaluation(studentID, criterion.id); exists {
			weighted += evaluation.score / criterion.maximum * weight
			result.Graded++
		}
	}
	policy := s.policy
	if policy == nil {
		policy = MissingAsZero{}
	}
	if err := policy.Check(result.Graded, result.Criteria); err != nil {
		return Total{}, err
	}
	result.Percent = min(percentScale, weighted/weights*percentScale)
	return result, nil
}
