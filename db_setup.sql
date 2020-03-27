create table rail_journeys
(
    id int auto_increment,
    journey_type text null,
    departing text null,
    destination text null,
    ticket_name text null,
    date date null,
    railcard_used boolean null,
    cost double null,
    total_cost double null,
    constraint rail_journeys_pk
        primary key (id)
);

