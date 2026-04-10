-- +goose Up
-- Indexes for better query performance
CREATE UNIQUE INDEX IF NOT EXISTS idx_tokens_refresh_token ON tokens(refresh_token);
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_boards_user_id ON boards(user_id);
CREATE INDEX IF NOT EXISTS idx_notes_board_id ON notes(board_id);
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);

-- Update foreign key to CASCADE deletion
ALTER TABLE notes 
DROP CONSTRAINT IF EXISTS notes_board_id_fkey,
ADD CONSTRAINT notes_board_id_fkey 
    FOREIGN KEY (board_id) 
    REFERENCES boards(id) 
    ON DELETE CASCADE;

-- Function to sync boards.notes counter
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_board_note_counter()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        IF NEW.board_id IS NOT NULL THEN
            UPDATE boards SET notes = notes + 1 WHERE id = NEW.board_id;
        END IF;
    ELSIF (TG_OP = 'DELETE') THEN
        IF OLD.board_id IS NOT NULL THEN
            UPDATE boards SET notes = notes - 1 WHERE id = OLD.board_id;
        END IF;
    ELSIF (TG_OP = 'UPDATE') THEN
        IF OLD.board_id IS DISTINCT FROM NEW.board_id THEN
            IF OLD.board_id IS NOT NULL THEN
                UPDATE boards SET notes = notes - 1 WHERE id = OLD.board_id;
            END IF;
            IF NEW.board_id IS NOT NULL THEN
                UPDATE boards SET notes = notes + 1 WHERE id = NEW.board_id;
            END IF;
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function on notes table changes
DROP TRIGGER IF EXISTS trg_note_counter ON notes;
CREATE TRIGGER trg_note_counter
AFTER INSERT OR DELETE OR UPDATE ON notes
FOR EACH ROW EXECUTE FUNCTION update_board_note_counter();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_note_counter ON notes;
DROP FUNCTION IF EXISTS update_board_note_counter();
DROP INDEX IF EXISTS idx_notes_user_id;
DROP INDEX IF EXISTS idx_notes_board_id;
DROP INDEX IF EXISTS idx_boards_user_id;
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_refresh_token;
