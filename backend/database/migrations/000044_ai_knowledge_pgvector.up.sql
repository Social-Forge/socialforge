-- Enable pgvector and add a vector column for semantic knowledge retrieval.
-- The legacy `embedding` JSONB column is kept for backward-compat; `embedding_vec`
-- is the queryable vector (OpenAI text-embedding-3-small = 1536 dims).
CREATE EXTENSION IF NOT EXISTS vector;

ALTER TABLE ai_knowledge ADD COLUMN IF NOT EXISTS embedding_vec vector(1536);

-- HNSW index for cosine-distance ANN search (pgvector >= 0.5). Builds fine on an
-- empty table (no training needed, unlike ivfflat).
CREATE INDEX IF NOT EXISTS idx_ai_knowledge_embedding_vec
    ON ai_knowledge USING hnsw (embedding_vec vector_cosine_ops);
