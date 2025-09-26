-- Pact Broker PostgreSQL Initialization Script
-- Optimized for cross-platform contract testing with 7 microservices

-- Create additional indexes for performance optimization
-- These improve query performance for contract verification workflows

-- Index for consumer/provider lookups (frequent in cross-platform testing)
CREATE INDEX IF NOT EXISTS idx_pact_publications_consumer_provider
  ON pact_publications(consumer_id, provider_id);

-- Index for version lookups (critical for CI/CD workflows)
CREATE INDEX IF NOT EXISTS idx_pact_versions_number_order
  ON pact_versions(pacticipant_id, "number", "order");

-- Index for webhook delivery tracking
CREATE INDEX IF NOT EXISTS idx_webhook_executions_created_at
  ON webhook_executions(created_at);

-- Performance tuning for large-scale contract testing
-- Adjust shared_buffers for better performance
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET work_mem = '4MB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';

-- Connection settings for high concurrency
ALTER SYSTEM SET max_connections = '200';

-- Logging configuration for debugging
ALTER SYSTEM SET log_statement = 'mod';
ALTER SYSTEM SET log_duration = 'on';
ALTER SYSTEM SET log_min_duration_statement = 1000;  -- Log slow queries

-- Apply settings
SELECT pg_reload_conf();

-- Create custom functions for Tchat-specific contract management
CREATE OR REPLACE FUNCTION get_consumer_provider_stats()
RETURNS TABLE(
  consumer_name text,
  provider_name text,
  latest_pact_version text,
  verification_count bigint,
  last_verified_at timestamp with time zone
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    c.name as consumer_name,
    p.name as provider_name,
    pp.revision as latest_pact_version,
    COALESCE(verification_stats.count, 0) as verification_count,
    verification_stats.last_verified_at
  FROM pacticipants c
  JOIN pact_publications pp ON c.id = pp.consumer_id
  JOIN pacticipants p ON pp.provider_id = p.id
  LEFT JOIN (
    SELECT
      pact_version_id,
      COUNT(*) as count,
      MAX(created_at) as last_verified_at
    FROM verifications
    GROUP BY pact_version_id
  ) verification_stats ON pp.pact_version_id = verification_stats.pact_version_id
  WHERE pp.id IN (
    SELECT MAX(id)
    FROM pact_publications
    GROUP BY consumer_id, provider_id
  )
  ORDER BY c.name, p.name;
END;
$$ LANGUAGE plpgsql;

-- Grant permissions
GRANT EXECUTE ON FUNCTION get_consumer_provider_stats() TO pact_broker_user;

-- Create a view for quick contract health overview
CREATE OR REPLACE VIEW contract_health_overview AS
SELECT
  c.name as consumer,
  p.name as provider,
  pp.revision as pact_version,
  CASE
    WHEN v.success = true THEN 'VERIFIED'
    WHEN v.success = false THEN 'FAILED'
    ELSE 'PENDING'
  END as status,
  v.created_at as last_verification_time
FROM pact_publications pp
JOIN pacticipants c ON pp.consumer_id = c.id
JOIN pacticipants p ON pp.provider_id = p.id
LEFT JOIN verifications v ON pp.pact_version_id = v.pact_version_id
WHERE pp.id IN (
  SELECT MAX(id)
  FROM pact_publications
  GROUP BY consumer_id, provider_id
);

-- Grant view access
GRANT SELECT ON contract_health_overview TO pact_broker_user;

-- Insert initial data for Tchat microservices
INSERT INTO pacticipants (name, repository_url, created_at, updated_at) VALUES
  ('tchat-web', 'https://github.com/tchat/tchat-web', NOW(), NOW()),
  ('tchat-ios', 'https://github.com/tchat/tchat-mobile', NOW(), NOW()),
  ('tchat-android', 'https://github.com/tchat/tchat-mobile', NOW(), NOW()),
  ('auth-service', 'https://github.com/tchat/tchat-backend', NOW(), NOW()),
  ('content-service', 'https://github.com/tchat/tchat-backend', NOW(), NOW()),
  ('commerce-service', 'https://github.com/tchat/tchat-backend', NOW(), NOW()),
  ('messaging-service', 'https://github.com/tchat/tchat-backend', NOW(), NOW()),
  ('payment-service', 'https://github.com/tchat/tchat-backend', NOW(), NOW()),
  ('notification-service', 'https://github.com/tchat/tchat-backend', NOW(), NOW()),
  ('gateway-service', 'https://github.com/tchat/tchat-backend', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Create environments for proper deployment tracking
INSERT INTO environments (name, display_name, production, created_at, updated_at) VALUES
  ('development', 'Development', false, NOW(), NOW()),
  ('staging', 'Staging', false, NOW(), NOW()),
  ('production', 'Production', true, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

COMMIT;