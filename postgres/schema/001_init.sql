-- +goose Up
CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    name VARCHAR UNIQUE,
    bed_number SMALLINT,
    patient_number SMALLINT,
    create_at TIMESTAMP,
    update_at TIMESTAMP
);

CREATE TABLE use_room(
    id UUID PRIMARY KEY,
    id_doctor VARCHAR,
    id_room UUID,
    CONSTRAINT fk_room FOREIGN KEY(id_room) REFERENCES rooms(id)
);

CREATE TABLE beds (
    id UUID PRIMARY KEY,
    id_room UUID,
    name VARCHAR,
    status VARCHAR,
    create_at TIMESTAMP,
    update_at TIMESTAMP,
    CONSTRAINT fk_room FOREIGN KEY(id_room) REFERENCES rooms(id)
);

CREATE TABLE patients (
    id UUID PRIMARY KEY,
    patient_code VARCHAR,
    name VARCHAR,
    phone  VARCHAR,
    address VARCHAR,
    create_at TIMESTAMP,
    update_at TIMESTAMP
);

CREATE TABLE records (
    id UUID PRIMARY KEY,
    id_patient UUID,
    id_doctor UUID,
    status VARCHAR,
    create_at TIMESTAMP,
    update_at TIMESTAMP,
    CONSTRAINT fk_patient FOREIGN KEY(id_patient) REFERENCES patients(id)
);

CREATE TABLE devices (
    id UUID PRIMARY KEY,
    serial VARCHAR UNIQUE ,
    warranty SMALLINT,
    status VARCHAR,
    create_at TIMESTAMP,
    update_at TIMESTAMP
);



CREATE TABLE use_bed(
    id UUID PRIMARY KEY,
    id_bed UUID,
    id_record UUID,
    status VARCHAR,
    create_at TIMESTAMP,
    end_at TIMESTAMP,
    CONSTRAINT fk_bed FOREIGN KEY(id_bed) REFERENCES beds(id),
    CONSTRAINT fk_record FOREIGN KEY (id_record) REFERENCES records(id)
);

CREATE TABLE use_device(
    id UUID PRIMARY KEY,
    id_device UUID,
    id_record UUID,
    status VARCHAR,
    create_at TIMESTAMP,
    end_at TIMESTAMP,
    CONSTRAINT fk_device FOREIGN KEY(id_device) REFERENCES devices(id),
    CONSTRAINT fk_record FOREIGN KEY (id_record) REFERENCES records(id)
);


-- +goose Down