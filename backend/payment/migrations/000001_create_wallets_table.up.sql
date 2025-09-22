-- Migration: create_wallets_table
-- Created at: 2024-12-20T10:06:00Z

-- Create wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    currency VARCHAR(10) NOT NULL CHECK (currency IN ('THB', 'SGD', 'IDR', 'MYR', 'PHP', 'VND', 'USD')),
    balance DECIMAL(20,2) DEFAULT 0.00 CHECK (balance >= 0),
    frozen_balance DECIMAL(20,2) DEFAULT 0.00 CHECK (frozen_balance >= 0),
    available_balance DECIMAL(20,2) GENERATED ALWAYS AS (balance - frozen_balance) STORED,
    daily_limit DECIMAL(20,2) DEFAULT 50000.00,
    monthly_limit DECIMAL(20,2) DEFAULT 1000000.00,
    daily_spent DECIMAL(20,2) DEFAULT 0.00,
    monthly_spent DECIMAL(20,2) DEFAULT 0.00,
    last_transaction_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'frozen', 'suspended', 'closed')),
    kyc_level INT DEFAULT 0 CHECK (kyc_level BETWEEN 0 AND 3),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, currency)
);

-- Create indexes
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_currency ON wallets(currency);
CREATE INDEX idx_wallets_status ON wallets(status);
CREATE INDEX idx_wallets_kyc_level ON wallets(kyc_level);
CREATE INDEX idx_wallets_balance ON wallets(balance);
CREATE INDEX idx_wallets_last_transaction_at ON wallets(last_transaction_at);
CREATE INDEX idx_wallets_created_at ON wallets(created_at);

-- Create composite index for user's active wallets
CREATE INDEX idx_wallets_user_active ON wallets(user_id, status)
    WHERE status = 'active';

-- Create updated_at trigger
CREATE TRIGGER update_wallets_updated_at
    BEFORE UPDATE ON wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();