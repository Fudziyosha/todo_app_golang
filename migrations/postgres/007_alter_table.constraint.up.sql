ALTER TABLE todo
DROP CONSTRAINT todo_created_in_list_fkey,
ADD CONSTRAINT todo_created_in_list_fkey
FOREIGN KEY(created_in_list) REFERENCES list(id)
on delete cascade;