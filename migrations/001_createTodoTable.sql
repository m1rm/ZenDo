DROP TABLE IF EXISTS todos;
CREATE TABLE todos (
                      id          INT AUTO_INCREMENT NOT NULL,
                      description VARCHAR(255),
                      status      TINYINT NOT NULL DEFAULT 0,
                      PRIMARY KEY (`id`)
);

INSERT INTO todos
(description, status)
VALUES
    ('update UI after api calls', '0'),
    ('implement edit todo in API', '0' ),
    ('implement edit todo in FE', '0' ),
    ('add status change functionality in FE', '0'),
    ('improve project structure (goApi)', '1' ),
    ('improve project structure (svelteKit app)', '0' ),
    ('add pagination BE', '0'),
    ('add pagination FE', '0');
