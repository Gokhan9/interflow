CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
    
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_value VARCHAR(255) UNIQUE NOT NULL, 
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE usage_logs (
    id SERIAL PRIMARY KEY,
    api_key_id INTEGER NOT NULL REFERENCES api_keys(id), -- Foreign key to api_keys table
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    prompt_tokens INTEGER DEFAULT 0, -- API isteğinde kullanılan prompt token sayısı
    completion_tokens INTEGER DEFAULT 0, -- API isteğinde kullanılan tamamlanan token sayısı
    total_tokens INTEGER DEFAULT 0, -- API isteğinde kullanılan toplam token sayısı
    latency_ms INTEGER, -- İstek işleme süresi(ms cinsinden)
    status_code INTEGER, -- API isteğinin HTTP durum kodu
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);