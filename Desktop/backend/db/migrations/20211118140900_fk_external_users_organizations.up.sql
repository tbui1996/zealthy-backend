ALTER TABLE users.external_users
ADD COLUMN external_user_organization_id int;

ALTER TABLE users.external_users
ADD CONSTRAINT fk_external_users_external_user_organizations FOREIGN KEY (external_user_organization_id) REFERENCES users.external_user_organizations (id);