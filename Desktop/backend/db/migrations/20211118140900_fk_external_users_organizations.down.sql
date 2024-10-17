ALTER TABLE IF EXISTS users.external_users
  DROP CONSTRAINT fk_external_users_external_user_organizations;

ALTER TABLE IF EXISTS users.external_users
  DROP COLUMN external_user_organization_id;