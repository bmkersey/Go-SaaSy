-- +goose Up
CREATE TABLE password_resets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_password_reset_user_id
      FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);

-- +goose Down
DROP TABLE password_resets;
