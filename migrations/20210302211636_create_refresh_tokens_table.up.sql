CREATE TABLE refresh_tokens(
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    ua varchar(200) NOT NULL, /* user-agent */
    fingerprint varchar(200) NOT NULL,
    ip varchar(39) NOT NULL,
    expires_in bigint NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);