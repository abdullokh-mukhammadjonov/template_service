ALTER TABLE if exists file_track 
        ADD COLUMN if not exists file_name_id uuid DEFAULT NULL;