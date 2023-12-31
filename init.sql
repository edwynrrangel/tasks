CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA "usr";

CREATE SCHEMA "task";

CREATE TABLE "usr"."roles" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" VARCHAR,
  "description" TEXT
);

CREATE TABLE "usr"."users" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "username" VARCHAR UNIQUE,
  "password" VARCHAR,
  "role_id" UUID,
  "first_login" BOOLEAN DEFAULT true,
  "created_at" TIMESTAMP DEFAULT (now()),
  "updated_at" TIMESTAMP DEFAULT (now()),
  "last_login_at" TIMESTAMP,
  "active" BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE "usr"."tokens" (
  "user_id" UUID DEFAULT (uuid_generate_v4()),
  "jwt_token" VARCHAR,
  "expiration" TIMESTAMP,
  "created_at" TIMESTAMP DEFAULT (now()),
  "first_login" BOOLEAN
);

CREATE TABLE "task"."tasks" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "title" VARCHAR,
  "description" TEXT,
  "due_date" DATE,
  "status" UUID,
  "assigned_user" UUID,
  "created_at" TIMESTAMP DEFAULT (now()),
  "updated_at" TIMESTAMP DEFAULT (now())
);

CREATE TABLE "task"."states" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" VARCHAR UNIQUE
);

CREATE TABLE "task"."state_transitions" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "current_state_id" UUID,
  "next_state_id" UUID
);

CREATE TABLE "task"."comments" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "task_id" UUID,
  "comment" TEXT,
  "created_by" UUID,
  "created_at" TIMESTAMP DEFAULT (now())
);

COMMENT ON COLUMN "task"."states"."id" IS 'Unique identifier for each state';

COMMENT ON COLUMN "task"."states"."name" IS 'Name of the state, e.g., "Asignado", "Iniciado"';

COMMENT ON COLUMN "task"."state_transitions"."id" IS 'Unique identifier for each state transition';

COMMENT ON COLUMN "task"."state_transitions"."current_state_id" IS 'State ID from which a transition can occur';

COMMENT ON COLUMN "task"."state_transitions"."next_state_id" IS 'State ID to which a transition can occur';

ALTER TABLE "usr"."users" ADD FOREIGN KEY ("role_id") REFERENCES "usr"."roles" ("id");

ALTER TABLE "usr"."tokens" ADD FOREIGN KEY ("user_id") REFERENCES "usr"."users" ("id");

ALTER TABLE "task"."tasks" ADD FOREIGN KEY ("status") REFERENCES "task"."states" ("id");

ALTER TABLE "task"."tasks" ADD FOREIGN KEY ("assigned_user") REFERENCES "usr"."users" ("id");

ALTER TABLE "task"."state_transitions" ADD FOREIGN KEY ("current_state_id") REFERENCES "task"."states" ("id");

ALTER TABLE "task"."state_transitions" ADD FOREIGN KEY ("next_state_id") REFERENCES "task"."states" ("id");

ALTER TABLE "task"."comments" ADD FOREIGN KEY ("task_id") REFERENCES "task"."tasks" ("id");

ALTER TABLE "task"."comments" ADD FOREIGN KEY ("created_by") REFERENCES "usr"."users" ("id");

-- Inserción de roles
INSERT INTO usr.roles (name, description)
VALUES 
('Administrator', 'Usuario con derechos administrativos'),
('Executor', 'Usuario responsable de la ejecución de las tareas'),
('Auditor', 'Usuario responsable de las tareas de auditoría');

-- Selección del ID del role de administrador
DO $$ 
DECLARE 
    admin_role_id UUID;
BEGIN 
    SELECT id INTO admin_role_id FROM usr.roles WHERE name = 'Administrator' LIMIT 1;
    -- Inserción del primer usuario administrador
    INSERT INTO usr.users (username, password, role_id)
    VALUES 
    ('admin', null, admin_role_id);
END $$;

-- Inserción de estados de tarea
INSERT INTO task.states (name) VALUES 
('Asignado'), 
('Iniciado'), 
('En espera'), 
('Finalizado Éxito'), 
('Finalizado Error');


-- Del estado Asignado solo puede pasar a Iniciado
INSERT INTO task.state_transitions (current_state_id, next_state_id)
SELECT s1.id, s2.id 
FROM task.states s1, task.states s2 
WHERE s1.name = 'Asignado' AND s2.name = 'Iniciado';

-- Del estado Iniciado puede pasar a Finalizado Exito, Finalizado Error, o En espera
INSERT INTO task.state_transitions (current_state_id, next_state_id)
SELECT s1.id, s2.id 
FROM task.states s1, task.states s2 
WHERE s1.name = 'Iniciado' AND s2.name IN ('Finalizado Éxito', 'Finalizado Error', 'En espera');

INSERT INTO task.state_transitions (current_state_id, next_state_id)
SELECT s1.id, s2.id 
FROM task.states s1, task.states s2 
WHERE s1.name = 'En espera' AND s2.name IN ('Finalizado Éxito', 'Finalizado Error');