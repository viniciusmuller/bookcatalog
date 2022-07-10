alter table Documents add column title VARCHAR(250);

alter table Documents add column author varchar(250);

alter table Documents add column pages int;

alter table Documents rename name to filename;
