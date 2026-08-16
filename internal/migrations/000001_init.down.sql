
DROP INDEX IF EXISTS messages_created_at_idx;
DROP TABLE IF EXISTS messages;


DROP INDEX IF EXISTS support_metrics_date_idx;
DROP INDEX IF EXISTS support_metrics_channel_idx;
DROP INDEX IF EXISTS support_metrics_agent_idx;
DROP INDEX IF EXISTS support_metrics_date_channel_idx;

DROP TABLE IF EXISTS support_metrics;