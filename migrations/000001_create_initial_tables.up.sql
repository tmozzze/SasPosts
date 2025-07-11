CREATE TABLE if NOT EXISTS posts (
    id             VARCHAR(255) PRIMARY KEY,
    title          TEXT NOT NULL,
    content        TEXT NOT NULL,
    author         VARCHAR(255) NOT NULL,
    allow_comments BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE if NOT EXISTS comments (
    id         VARCHAR(255) PRIMARY KEY,
    post_id    VARCHAR(255) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id  VARCHAR(255) REFERENCES comments(id) ON DELETE CASCADE,
    author     VARCHAR(255) NOT NULL,
    content    TEXT NOT NULL,
    path       TEXT,
    depth      INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_path ON comments(path);
