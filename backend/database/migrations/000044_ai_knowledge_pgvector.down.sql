DROP INDEX IF EXISTS idx_ai_knowledge_embedding_vec;
ALTER TABLE ai_knowledge DROP COLUMN IF EXISTS embedding_vec;
-- Extension intentionally left installed (other objects may depend on it).
