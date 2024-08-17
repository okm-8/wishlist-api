/*
 * Created at: 2024-08-17 15:51:03
 * Description: wishlist model
 */

CREATE TABLE wishlists
(
    id          UUID PRIMARY KEY,
    wisher_id   UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    hidden      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    FOREIGN KEY (wisher_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_wishlist_wisher_id ON wishlists (wisher_id);
CREATE INDEX idx_wishlist_created_at ON wishlists (created_at);

CREATE TABLE wishes
(
    id          UUID PRIMARY KEY,
    wishlist_id UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    hidden      BOOLEAN      NOT NULL DEFAULT FALSE,
    fulfilled   BOOLEAN      NOT NULL DEFAULT FALSE,
    assignee_id UUID         NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000', -- NULL UUID
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    FOREIGN KEY (wishlist_id) REFERENCES wishlists (id) ON DELETE CASCADE,
    FOREIGN KEY (assignee_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_wish_wishlist_id ON wishes (wishlist_id);
CREATE INDEX idx_wish_created_at ON wishes (created_at);
CREATE INDEX idx_wish_assignee_id ON wishes (assignee_id);


CREATE OR REPLACE VIEW wishes_view AS (
    SELECT
        wish.id AS id,
        wish.name AS name,
        wish.description AS description,
        wish.hidden AS hidden,
        wish.fulfilled AS fulfilled,
        wish.created_at AS created_at,
        assignee.id AS assignee_id,
        assignee.email AS assignee_email,
        assignee.name AS assignee_name,
        wishlist.id AS wishlist_id,
        wishlist.name AS wishlist_name,
        wishlist.description AS wishlist_description,
        wishlist.hidden AS wishlist_hidden,
        wishlist.created_at AS wishlist_created_at,
        wisher.id AS wisher_id,
        wisher.email AS wisher_email,
        wisher.name AS wisher_name
    FROM wishes AS wish
    INNER JOIN wishlists AS wishlist ON wish.wishlist_id = wishlist.id
    INNER JOIN users AS wisher ON wishlist.wisher_id = wisher.id
    LEFT JOIN users AS assignee ON wish.assignee_id = assignee.id
);
