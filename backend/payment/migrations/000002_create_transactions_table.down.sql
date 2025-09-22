-- Rollback migration: create_transactions_table
-- Created at: 2024-12-20T10:07:00Z

-- Drop trigger
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;

-- Drop indexes
DROP INDEX IF EXISTS idx_transactions_metadata;
DROP INDEX IF EXISTS idx_transactions_pending;
DROP INDEX IF EXISTS idx_transactions_wallet_status;
DROP INDEX IF EXISTS idx_transactions_user_status;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_completed_at;
DROP INDEX IF EXISTS idx_transactions_processed_at;
DROP INDEX IF EXISTS idx_transactions_expires_at;
DROP INDEX IF EXISTS idx_transactions_parent_transaction_id;
DROP INDEX IF EXISTS idx_transactions_reference_id;
DROP INDEX IF EXISTS idx_transactions_gateway_transaction_id;
DROP INDEX IF EXISTS idx_transactions_gateway_provider;
DROP INDEX IF EXISTS idx_transactions_payment_method;
DROP INDEX IF EXISTS idx_transactions_currency;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_to_user_id;
DROP INDEX IF EXISTS idx_transactions_from_user_id;
DROP INDEX IF EXISTS idx_transactions_to_wallet_id;
DROP INDEX IF EXISTS idx_transactions_from_wallet_id;
DROP INDEX IF EXISTS idx_transactions_external_id;

-- Drop table
DROP TABLE IF EXISTS transactions;