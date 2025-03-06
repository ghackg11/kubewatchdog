CREATE TABLE kubernetes_events (
                                   id TEXT,
                                   event_time TIMESTAMPTZ NOT NULL ,
                                   event_type TEXT,
                                   reason TEXT,
                                   message TEXT,
                                   namespace TEXT,
                                   resource TEXT,
                                   resource_name TEXT,
                                   PRIMARY KEY (id, event_time)
);


SELECT create_hypertable('kubernetes_events', 'event_time');
