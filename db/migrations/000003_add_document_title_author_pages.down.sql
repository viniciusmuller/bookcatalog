alter table Documents drop column title VARCHAR(250);

alter table Documents drop column author varchar(250);

alter table Documents drop column pages int;

alter table Documents rename filename to name;
