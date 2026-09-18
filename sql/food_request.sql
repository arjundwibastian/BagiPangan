
CREATE TABLE food_requests (
    request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,

    latitude DECIMAL(9,6) NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'searching',
    radius_km DECIMAL(6,2) NOT NULL DEFAULT 5,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_food_requests_latitude
        CHECK (latitude BETWEEN -90 AND 90),

    CONSTRAINT chk_food_requests_longitude
        CHECK (longitude BETWEEN -180 AND 180),

    CONSTRAINT chk_food_requests_radius
        CHECK (radius_km > 0),

    CONSTRAINT chk_food_requests_status
        CHECK (status IN (
            'searching',
            'claimed',
            'completed',
            'cancelled',
            'expired'
        ))
);

CREATE INDEX idx_food_requests_user_id
    ON food_requests (user_id);

CREATE INDEX idx_food_requests_status
    ON food_requests (status);

CREATE INDEX idx_food_requests_latitude_longitude
    ON food_requests (latitude, longitude);