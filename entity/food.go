package entity

import (
	"math/rand"
	"time"
)

// randSource — локальный источник случайных чисел, инициализируется один раз.
// Используется вместо глобального math/rand для предсказуемости в тестах
// и отсутствия гонок данных при параллельном вызове.
var randSource *rand.Rand

func init() {
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// SetRandSource позволяет установить детерминированный источник для тестов.
// Вызывается перед тестами, требующими воспроизводимого поведения.
func SetRandSource(r *rand.Rand) {
	randSource = r
}

// Movable описывает объект, который может двигаться по полю
type Movable interface {
	// Move выполняет одно движение
	Move()
	// GetHead возвращает текущую позицию головы
	GetHead() Point
}

// Renderable описывает объект, который можно отрисовать на экране
type Renderable interface {
	// GetPos возвращает позицию для рендеринга
	GetPos() Point
}

// Food представляет еду на поле
type Food struct {
	Position Point
}

// NewFood создаёт новую еду в случайной позиции
func NewFood(snakeBody []Point) *Food {
	f := &Food{}
	f.Spawn(snakeBody)
	return f
}

// GetPos возвращает позицию еды (реализует Renderable)
func (f *Food) GetPos() Point {
	return f.Position
}

// Spawn размещает еду в случайной пустой клетке.
// Использует rand.Intn для выбора позиции; избегает клеток, занятых змейкой.
// Оптимизированная версия: при длинной змейке генерирует список свободных клеток.
func (f *Food) Spawn(snakeBody []Point) {
	occupied := make(map[Point]bool, len(snakeBody))
	for _, p := range snakeBody {
		occupied[p] = true
	}

	// Если змейка занимает больше половины поля, используем оптимизированный подход
	freeCells := GridWidth * GridHeight - len(snakeBody)
	if freeCells <= 0 {
		return // поле полностью заполнено
	}

	if freeCells < GridWidth*GridHeight/2 {
		// Оптимизация: генерируем список свободных клеток и выбираем случайный индекс
		free := make([]Point, 0, freeCells)
		for x := 0; x < GridWidth; x++ {
			for y := 0; y < GridHeight; y++ {
				pos := Point{X: x, Y: y}
				if !occupied[pos] {
					free = append(free, pos)
				}
			}
		}
		if len(free) > 0 {
			f.Position = free[randSource.Intn(len(free))]
			return
		}
	}

	// Стандартный подход: случайный выбор с проверкой
	for {
		x := randSource.Intn(GridWidth)
		y := randSource.Intn(GridHeight)
		pos := Point{X: x, Y: y}
		if !occupied[pos] {
			f.Position = pos
			return
		}
	}
}
