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
    ('Style FE', '1' ),
    ('Add CTAs w/o functionality', '1' ),
    ('dockerize dev', '1' ),
    ('replace mock API responses data with real data', '1' ),
    ('Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea', 0 ),
    ('improve project (goApi) structure', '0' ),
    ('implement missing API functionality', '0' ),
    ('make CTAs in FE functional', '0' ),
    ('add APIdoc', '0' ),
    ('find and use migration tooling', '0' ),
    ('prod setup', '0' );
