drop database if exists hospital;
create database hospital;

create user compfestadmin with password 'compfestadmin';
revoke connect on database hospital from public;
grant connect on database hospital to compfestadmin;

revoke all on all tables in schema public from PUBLIC;
grant select, insert, update, delete on all tables in schema public to compfestadmin;
alter default privileges for user compfestadmin in schema public grant select, insert, update, delete on tables to compfestadmin;