CREATE TABLE IF NOT EXISTS `User` (
    Email varchar(50) NOT NULL,
    Fullname text NOT NULL,
    Birthdate DATE NOT NULL ,
    Mentor varchar(50),
    Start_Date DATE,
    End_Date DATE,
    PRIMARY KEY (Email)
) ENGINE=INNODB;

ALTER TABLE `User`
    ADD FOREIGN KEY (Mentor) REFERENCES User(Email) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS `Room` (
    name char(4) NOT NULL,
    description text,
    max_capacity int unsigned,
    PRIMARY KEY (name)
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `Team` (
    id varchar(4),
    description text,
    team_lead varchar(50),
    PRIMARY KEY (id)
) ENGINE=INNODB;

ALTER TABLE `Team`
    ADD FOREIGN KEY (team_lead) REFERENCES User(Email);

CREATE TABLE IF NOT EXISTS `Booking` (
    id int AUTO_INCREMENT NOT NULL,
    booker varchar(50) NOT NULL,
    room char(4) NOT NULL,
    start_time datetime NOT NULL,
    end_time datetime NOT NULL,
    CHECK (Booking.start_time < Booking.end_time),
    PRIMARY KEY (id)
) ENGINE=INNODB;

ALTER TABLE `Booking`
    ADD FOREIGN KEY (room) REFERENCES Room(name),
    ADD FOREIGN KEY (booker) REFERENCES User(Email);

CREATE TABLE IF NOT EXISTS `UserBooking` (
    id int NOT NULL,
    tagged_user varchar(50) NOT NULL,
    PRIMARY KEY (id, tagged_user)
) ENGINE=INNODB;

ALTER TABLE `UserBooking`
    ADD FOREIGN KEY (id) REFERENCES Booking(id) ON DELETE CASCADE,
    ADD FOREIGN KEY (tagged_user) REFERENCES User(Email);

CREATE TABLE IF NOT EXISTS `BelongToTeam` (
    staff varchar(50) NOT NULL,
    team varchar(4) NOT NULL,
    Start_Date date NOT NULL,
    End_Date date DEFAULT NULL,
    PRIMARY KEY (staff, team)
) ENGINE=INNODB;

ALTER TABLE `BelongToTeam`
    ADD FOREIGN KEY (staff) REFERENCES User(Email),
    ADD FOREIGN KEY (team) REFERENCES Team(id);