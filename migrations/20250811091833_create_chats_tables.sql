-- +goose Up
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE chat_users (
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    username TEXT NOT NULL
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- Индексы для ускорения выборки
CREATE INDEX idx_chat_users_chat_id ON chat_users(chat_id);
CREATE INDEX idx_messages_chat_id ON messages(chat_id);

-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chat_users;
DROP TABLE IF EXISTS chats;
