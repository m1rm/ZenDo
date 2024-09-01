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
    ('Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea', 1 ),
    ('implement edit todo in API', '0' ),
    ('implement edit todo in FE', '0' ),
    ('add status change functionality in FE', '0'),
    ('improve project structure (goApi)', '0' ),
    ('improve project structure (svelteKit app)', '0' ),
    ('add pagination BE', '0'),
    ('add pagination FE', '0');
