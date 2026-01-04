-- Initial database setup for Pet Service
-- This file will be automatically executed when the PostgreSQL container starts

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Init roles
INSERT INTO roles (id, name, created_at, updated_at, created_by, updated_by, is_active)
VALUES 
    ('DAA6B933DAF84FBB99477508A4DAC571', 'Admin', NOW(), NOW(), '', '', true),
    ('DAA6B933DAF84FBB99477508A4DAC572', 'Editor', NOW(), NOW(), '', '', true),
    ('DAA6B933DAF84FBB99477508A4DAC573', 'User', NOW(), NOW(), '', '', true);

-- Init permissions
INSERT INTO permissions (id, name, created_at, updated_at, created_by, updated_by, is_active)
VALUES
    ('A3CBE2F6B7CD9A34D0FA23593D0E42FA', 'view_pet', NOW(), NOW(), '', '', true),
    ('B02AC61D6F5B8D83111F42A7C75C1D2A', 'add_pet', NOW(), NOW(), '', '', true),
    ('E0B57C9A19E0412399B72391FC5D50CC', 'edit_pet', NOW(), NOW(), '', '', true),
    ('F53A2B89FFB1B3C42B9E6814E5869338', 'delete_pet', NOW(), NOW(), '', '', true),
    ('6D02A9C9378D1D86F39B3BDA48DA9F8E', 'view_user', NOW(), NOW(), '', '', true),
    ('6D02A9C9378D1D86F39B3BDA48DA9F8F', 'add_user', NOW(), NOW(), '', '', true),
    ('6D02A9C9378D1D86F39B3BDA48DA9F8D', 'edit_user', NOW(), NOW(), '', '', true),
    ('6D02A9C9378D1D86F39B3BDA48DA9F8G', 'delete_user', NOW(), NOW(), '', '', true);

-- Init role permissions
INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at, created_by, updated_by, is_active)
VALUES
    -- Editor permissions
    ('1D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC572', 'A3CBE2F6B7CD9A34D0FA23593D0E42FA', NOW(), NOW(), '', '', true), -- view_pet
    ('2D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC572', 'B02AC61D6F5B8D83111F42A7C75C1D2A', NOW(), NOW(), '', '', true), -- add_pet
    ('3D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC572', 'E0B57C9A19E0412399B72391FC5D50CC', NOW(), NOW(), '', '', true), -- edit_pet
    ('4D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC572', 'F53A2B89FFB1B3C42B9E6814E5869338', NOW(), NOW(), '', '', true), -- delete_pet
    ('5D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC572', '6D02A9C9378D1D86F39B3BDA48DA9F8E', NOW(), NOW(), '', '', true), -- view_user
    -- User permissions
    ('6D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC573', '6D02A9C9378D1D86F39B3BDA48DA9F8E', NOW(), NOW(), '', '', true), -- view_user
    ('7D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC573', '6D02A9C9378D1D86F39B3BDA48DA9F8F', NOW(), NOW(), '', '', true), -- add_user
    ('8D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC573', '6D02A9C9378D1D86F39B3BDA48DA9F8D', NOW(), NOW(), '', '', true); -- edit_user

-- Init users (password: 123456)
INSERT INTO users (id, first_name, last_name, email, phone, username, gender, password, is_admin, created_at, updated_at, created_by, updated_by, is_active)
VALUES 
    ('9D02A9C9378D1D86F39B3BDA48DA9F8G', 'Admin', 'User', 'admin@example.com', '0123456789', 'admin', true, '$2b$12$1hw8y1rphV0JGChNge5/Lu.YMLCdOl/CX/q9YwLfgXi4Z1.u3qJ.i', true, NOW(), NOW(), '', '', true),
    ('1102A9C9378D1D86F39B3BDA48DA9F8G', 'Editor', 'User', 'editor@example.com', '0123456789', 'editor', true, '$2b$12$1hw8y1rphV0JGChNge5/Lu.YMLCdOl/CX/q9YwLfgXi4Z1.u3qJ.i', false, NOW(), NOW(), '', '', true),
    ('1002A9C9378D1D86F39B3BDA48DA9F8G', 'User', 'Test', 'user@example.com', '0123456789', 'user', true, '$2b$12$1hw8y1rphV0JGChNge5/Lu.YMLCdOl/CX/q9YwLfgXi4Z1.u3qJ.i', false, NOW(), NOW(), '', '', true);

-- Init user roles
INSERT INTO user_roles (id, user_id, role_id, created_at, updated_at, created_by, updated_by, is_active)
VALUES
    ('AD02A9C9378D1D86F39B3BDA48DA9F8G', '9D02A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC571', NOW(), NOW(), '', '', true), -- admin -> Admin role
    ('BD02A9C9378D1D86F39B3BDA48DA9F8G', '1102A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC572', NOW(), NOW(), '', '', true), -- editor -> Editor role
    ('CD02A9C9378D1D86F39B3BDA48DA9F8G', '1002A9C9378D1D86F39B3BDA48DA9F8G', 'DAA6B933DAF84FBB99477508A4DAC573', NOW(), NOW(), '', '', true); -- user -> User role
