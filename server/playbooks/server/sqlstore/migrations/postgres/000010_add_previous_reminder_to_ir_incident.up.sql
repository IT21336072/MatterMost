ALTER TABLE IR_Incident ADD COLUMN IF NOT EXISTS PreviousReminder BIGINT NOT NULL DEFAULT 0;
