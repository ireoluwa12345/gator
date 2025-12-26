-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(200) NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_fetched_at TIMESTAMP WITH TIME ZONE NULL,

    foreign key (user_id) references users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feeds;
-- +goose StatementEnd
