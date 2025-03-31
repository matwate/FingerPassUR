create table Usuarios (
  id serial primary key,
  correo varchar(255),
  nombre varchar(255),
  programa varchar(255)
);


create table Images (
  id serial primary key,
  path varchar(512),
  user_id integer references Usuarios(id)
)
