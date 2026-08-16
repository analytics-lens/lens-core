CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    body TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ,

    sender_role TEXT NOT NULL
        CHECK (sender_role IN ('ai', 'user'))
);

CREATE INDEX IF NOT EXISTS messages_created_at_idx
    ON messages(created_at);


CREATE TABLE IF NOT EXISTS support_metrics (
    id BIGSERIAL PRIMARY KEY,

    date DATE NOT NULL,
    channel TEXT NOT NULL,
    agent TEXT NOT NULL,

    tickets INTEGER NOT NULL,
    tickets_resolved INTEGER NOT NULL,

    average_response_time NUMERIC(10, 2),
    first_response_time NUMERIC(10, 2),

    resolution_rate NUMERIC(5, 2),
    automation_rate NUMERIC(5, 2),
    escalation_rate NUMERIC(5, 2),

    csat NUMERIC(3, 2),

    bot_conversions INTEGER,
    revenue_attributed_to_bot NUMERIC(12, 2),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS support_metrics_date_idx
    ON support_metrics(date);

CREATE INDEX IF NOT EXISTS support_metrics_channel_idx
    ON support_metrics(channel);

CREATE INDEX IF NOT EXISTS support_metrics_agent_idx
    ON support_metrics(agent);

CREATE INDEX IF NOT EXISTS support_metrics_date_channel_idx
    ON support_metrics(date, channel);