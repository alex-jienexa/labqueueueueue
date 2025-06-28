-- Пользователи-студенты : student.go
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE
);

-- Очереди : queue.go
CREATE TABLE IF NOT EXISTS queues (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    starts_at TIMESTAMP NOT NULL,  -- время начала записи
    ends_at TIMESTAMP NOT NULL,    -- время окончания записи
    conflict_resolution_method VARCHAR(50) DEFAULT 'move_after', -- 'move_after', 'first_free', 'to_end'
    created_at TIMESTAMP DEFAULT NOW()
);

-- Запись в очередь : queue_entry.go
CREATE TABLE IF NOT EXISTS queue_entries (
    id SERIAL PRIMARY KEY,
    queue_id INTEGER REFERENCES queues(id) ON DELETE CASCADE,
    student_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    is_conflict BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Доступ к очередям : queue_student_access.go
CREATE TABLE IF NOT EXISTS student_queue_access (
    student_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    queue_id INTEGER REFERENCES queues(id) ON DELETE CASCADE,
    can_view BOOLEAN DEFAULT TRUE,
    can_join BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (student_id, queue_id)
);

-- Конйликтные ситуации : conflict.go
CREATE TABLE IF NOT EXISTS conflicts (
    id SERIAL PRIMARY KEY,
    queue_entry_id INTEGER REFERENCES queue_entries(id) ON DELETE CASCADE,
    student_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    resolved BOOLEAN DEFAULT FALSE,
    resolution_method VARCHAR(50)  -- 'dice_roll', 'cs2_duel', 'manual'
);

-- Ситуация продажи места : sale.go
CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    queue_entry_id INTEGER REFERENCES queue_entries(id) ON DELETE CASCADE,
    seller_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    buyer_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    price INTEGER NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE
);