-- Up migration: create initial tables for production
-- Enable pgcrypto extension for UUIDs (optional, if using UUIDs)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

