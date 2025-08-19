Create table if not EXISTS events (
id integer primary key AUTOINCREMENT,
owner_id integer not null,
name text not null,
description text not null,
date datetime not null,
location text not null,
foreign key (owner_id) references users(id) on delete cascade
);
