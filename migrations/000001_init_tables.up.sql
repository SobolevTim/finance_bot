-- Таблица пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id VARCHAR(255) UNIQUE NOT NULL,
    user_name VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    timezone VARCHAR(255) NOT NULL DEFAULT 'UTC+3',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица категорий
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    icon VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Таблица расходов
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_recurring BOOLEAN NOT NULL DEFAULT false,
    recurrence_rule VARCHAR(100),
    description TEXT,
    card_id UUID, -- Для будущей реализации карт
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Таблица бюджетов
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (end_date > start_date)
);

-- Связующая таблица для бюджета и категорий
CREATE TABLE budget_categories (
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    limit_amount NUMERIC(15,2) CHECK (limit_amount > 0),
    PRIMARY KEY (budget_id, category_id)
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_users_telegram_id ON users(telegram_id);
CREATE INDEX idx_expenses_user_date ON expenses(user_id, date);
CREATE INDEX idx_budgets_user_period ON budgets(user_id, start_date, end_date);
CREATE INDEX idx_categories_user ON categories(user_id);