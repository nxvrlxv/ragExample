CREATE EXTENTION IF NOT EXIST vector;

CREATE TABLE documents (
    id      bigserial primary key,
    title   text not null,
    source_path text not null default '',
    created_at timestamptz not null default now()
);

CREATE TABLE chunks (
    id           BIGSERIAL PRIMARY KEY,
    document_id  BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    content      TEXT NOT NULL,
    chunk_index  INT NOT NULL,
    embedding    vector(1024),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_chunks_embedding 
ON chunks USING hnsw (embedding vector_cosine_ops);
