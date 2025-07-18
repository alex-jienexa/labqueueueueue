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
	GetAll() ([]models.Queue, error)
	GetEntries(queueID int) ([]models.QueueEntry, error)
	GetByID(queueID int) (*models.Queue, error)

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

	// Проверяет, занята ли позиция position в очереди queueID.
	IsPositionBusy(queueID int, position int) (bool, error)
	// Перепроверяет все элементы очереди и назначает IsActive = true для тех,
	// у которых начинается время записи. Ровно и наоборот, если время записи уже
	// закончилось, ставит IsActive как false
	ManageActive() error

	// Выбирает уступившего/проигравшего в конфликтной ситуации и делает с ним то,
	// что указано в Queue.resolution_method
	ResolveConflict(loser *models.QueueEntry, queue *models.Queue) error
}

type QueueEntryRepository interface {
	// Создание нового элемента очереди в БД
	Create(queueEntry *models.QueueEntry) error
	// Обновить конфликтные ситуации для этого объекта и для других объектах данной позиции
	UpdateConflict(queueEntry *models.QueueEntry) error

	// Получение позиции в очереди по его ID
	GetByID(id int) (*models.QueueEntry, error)

	// Добавляет элемент в очередь. Позиция элемента зависит от значения queueEntry.Position.
	// Если queueEntry.Position == 0, то элемент добавляется в первый свободный слот очереди
	AddToQueue(queueEntry *models.QueueEntry, queueRepo queueRepository) error
}

type StudentQueueAccessRepository interface {
	Create(studentQueueAccess *models.StudentQueueAccess) error
	GetByStudent(studentEntry *models.Student) (*models.StudentQueueAccess, error)
	// Получает access для всех студентов в очереди
	GetByQueue(queueEntry *models.Queue) ([]models.StudentQueueAccess, error)

	// Обновляет доступ к очереди согласно данным из queueEntry.
	// До времени записи читать и писать запрещено.
	// Во время записи можно и читать, и писать.
	// После времени записи можно читать, но не писать.
	UpdateByQueueState(queueEntry *models.QueueEntry) error

	// Устанавливает значения флагов для конкретного студента в конкретной очереди
	SetFlags(accessEntry *models.StudentQueueAccess, canRead, canWrite bool) error
}
