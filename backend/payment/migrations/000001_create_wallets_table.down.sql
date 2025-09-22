-- Rollback migration: create_wallets_table
-- Created at: 2024-12-20T10:06:00Z

-- Drop trigger
DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;

-- Drop indexes
DROP INDEX IF EXISTS idx_wallets_user_active;
DROP INDEX IF EXISTS idx_wallets_created_at;
DROP INDEX IF EXISTS idx_wallets_last_transaction_at;
DROP INDEX IF EXISTS idx_wallets_balance;
DROP INDEX IF EXISTS idx_wallets_kyc_level;
DROP INDEX IF EXISTS idx_wallets_status;
DROP INDEX IF EXISTS idx_wallets_currency;
DROP INDEX IF EXISTS idx_wallets_user_id;

-- Drop table
DROP TABLE IF EXISTS wallets;