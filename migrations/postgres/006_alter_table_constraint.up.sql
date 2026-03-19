ALTER TABLE list
DROP CONSTRAINT list_created_by_fkey,
ADD CONSTRAINT list_created_by_fkey
FOREIGN KEY(created_by) REFERENCES users(id)
on delete cascade;

