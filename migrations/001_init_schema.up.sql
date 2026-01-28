-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);

COMMENT ON TABLE users IS 'Таблица пользователей системы';
COMMENT ON COLUMN users.user_id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN users.first_name IS 'Имя пользователя';
COMMENT ON COLUMN users.last_name IS 'Фамилия пользователя';
COMMENT ON COLUMN users.login IS 'Логин (email) для входа в систему';
COMMENT ON COLUMN users.password IS 'Хэшированный пароль пользователя';
COMMENT ON COLUMN users.last_login IS 'Дата и время последнего входа';