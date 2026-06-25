CREATE SCHEMA IF NOT EXISTS app;
CREATE TABLE app.word_lists (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    words TEXT NOT NULL
);
CREATE TABLE app.words_picked (
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL,
    word_list_id INT NOT NULL REFERENCES app.word_lists(id),
    picked_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE app.words_picked_history(
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL,
    word_list_id INT NOT NULL REFERENCES app.word_lists(id),
    picked_at TIMESTAMP NOT NULL
);

CREATE OR REPLACE FUNCTION app.words_picked_history_insert()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO app.words_picked_history (
        word,
        word_list_id,
        picked_at
    )
    VALUES (
        NEW.word,
        NEW.word_list_id,
        NEW.picked_at
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_words_picked_history
AFTER INSERT ON app.words_picked
FOR EACH ROW
EXECUTE FUNCTION app.words_picked_history_insert();