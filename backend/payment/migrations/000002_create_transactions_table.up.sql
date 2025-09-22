-- Migration: create_transactions_table
-- Created at: 2024-12-20T10:07:00Z

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) UNIQUE,
    from_wallet_id UUID REFERENCES wallets(id),
    to_wallet_id UUID REFERENCES wallets(id),
    from_user_id UUID,
    to_user_id UUID,
    type VARCHAR(20) NOT NULL CHECK (type IN ('transfer', 'topup', 'withdraw', 'payment', 'refund', 'fee', 'bonus')),
    amount DECIMAL(20,2) NOT NULL CHECK (amount > 0),
    fee DECIMAL(20,2) DEFAULT 0.00 CHECK (fee >= 0),
    total_amount DECIMAL(20,2) GENERATED ALWAYS AS (amount + fee) STORED,
    currency VARCHAR(10) NOT NULL,
    exchange_rate DECIMAL(20,8) DEFAULT 1.00000000,
    original_amount DECIMAL(20,2),
    original_currency VARCHAR(10),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled', 'refunded')),
    payment_method VARCHAR(50),
    gateway_provider VARCHAR(50),
    gateway_transaction_id VARCHAR(255),
    gateway_response JSONB,
    description TEXT,
    metadata JSONB DEFAULT '{}',
    reference_id VARCHAR(255),
    parent_transaction_id UUID REFERENCES transactions(id),
    expires_at TIMESTAMP WITH TIME ZONE,
    processed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_transactions_external_id ON transactions(external_id);
CREATE INDEX idx_transactions_from_wallet_id ON transactions(from_wallet_id);
CREATE INDEX idx_transactions_to_wallet_id ON transactions(to_wallet_id);
CREATE INDEX idx_transactions_from_user_id ON transactions(from_user_id);
CREATE INDEX idx_transactions_to_user_id ON transactions(to_user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_currency ON transactions(currency);
CREATE INDEX idx_transactions_payment_method ON transactions(payment_method);
CREATE INDEX idx_transactions_gateway_provider ON transactions(gateway_provider);
CREATE INDEX idx_transactions_gateway_transaction_id ON transactions(gateway_transaction_id);
CREATE INDEX idx_transactions_reference_id ON transactions(reference_id);
CREATE INDEX idx_transactions_parent_transaction_id ON transactions(parent_transaction_id);
CREATE INDEX idx_transactions_expires_at ON transactions(expires_at);
CREATE INDEX idx_transactions_processed_at ON transactions(processed_at);
CREATE INDEX idx_transactions_completed_at ON transactions(completed_at);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- Create composite indexes for common queries
CREATE INDEX idx_transactions_user_status ON transactions(from_user_id, status, created_at DESC);
CREATE INDEX idx_transactions_wallet_status ON transactions(from_wallet_id, status, created_at DESC);
CREATE INDEX idx_transactions_pending ON transactions(status, expires_at)
    WHERE status IN ('pending', 'processing');

-- Create GIN index for metadata search
CREATE INDEX idx_transactions_metadata ON transactions USING gin(metadata);

-- Create updated_at trigger
CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();