package repositories

import "github.com/alex-jienexa/labqueueueueue/models"

type StudentRepository interface {
	Create(student *models.Student) error
	GetByID(id int) (*models.Student, error)
	GetByUsername(username string) (*models.Student, error)
	GetByUsernameNoPassword(username string) (*models.Student, error)
}

type QueueRepository interface {
	Create(queue *models.Queue) error
	GetActive() (*models.Queue, error)
	GetEntries(queueID int) ([]models.QueueEntry, error)

	//! Вставка элементов в очередь

	// Вставляет уже существующий элемент очереди в позицию position очереди.
	// При этом, меняет позиции ВСЕХ элементов после position на +1
	MoveAndPush(queueEntry *models.QueueEntry, position int) error
	// Вставляет уже созданный элемент насильно в позицию очереди.
	// Если на данной позиции уже имеются элементы, то все элементы в нём отмечают как конфликтные
	MoveForce(entry *models.QueueEntry, position int) error
	// Перемещает уже существующий элемент очереди в первую свободную позицию дальше.
	// Используется в разрешении конфликтных ситуаций
	MoveToNextFree(entry *models.QueueEntry) error
	// Перемещает элемент в заданную позицию position очереди.
	// Если в данной позиции уже имеется элемент, то он становится прямо за ним.
	MoveAndFree(entry *models.QueueEntry, position int) error

	// Создаёт и вставляет новый элемент в очереди в позицию position.
	// Если position == 0, то элемент вставляется в первый свободный слот очереди.
	// Если position попадает на занятый элемент, то вместо этого вставится в первый свободный
	// элемент очереди после.
	//TODO: InsertToNextFree(queueEntry *models.QueueEntry, position int) error

	// Проверяет, занята ли позиция position в очереди queueID.
	IsPositionBusy(queueID int, position int) (bool, error)
}

type QueueEntryRepository interface {
	// Создание нового элемента очереди в БД
	Create(queueEntry *models.QueueEntry) error

	// Получение позиции в очереди по его ID
	GetByID(id int) (*models.QueueEntry, error)

	// Добавляет элемент в очередь. Позиция элемента зависит от значения queueEntry.Position.
	// Если queueEntry.Position == 0, то элемент добавляется в первый свободный слот очереди
	AddToQueue(queueEntry *models.QueueEntry, queueRepo queueRepository) error
}
