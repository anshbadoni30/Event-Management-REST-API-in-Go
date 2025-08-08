create table if not EXISTS attendees (
 id integer primary key AUTOINCREMENT,
 user_id integer not null,
 event_id integer not null,
 foreign key (user_id) references users(id) on delete cascade,
 foreign key (event_id) references events(id) on delete cascade
);

