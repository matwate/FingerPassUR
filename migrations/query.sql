create table Usuarios (
  id serial primary key,
  Correo varchar(255),
  Nombre varchar(255),
  Programa varchar(255)
);


create table Images (
  id serial primary key,
  path varchar(512),
  user_id integer references Usuarios(id)
)
