-- Migration Up: Create fingerprint_devices table and seed initial 10 Solution X100C machines

CREATE TABLE IF NOT EXISTS fingerprint_devices (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    port INTEGER NOT NULL DEFAULT 80,
    com_key INTEGER NOT NULL DEFAULT 0,
    location VARCHAR(150),
    is_online BOOLEAN DEFAULT TRUE,
    last_sync TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed initial 10 machines from field script
INSERT INTO fingerprint_devices (code, name, ip_address, port, location) VALUES
('t41', 'Mesin Solution X100C T41', '192.168.150.172', 80, 'Pos T41'),
('t42', 'Mesin Solution X100C T42', '192.168.150.179', 80, 'Pos T42'),
('t43', 'Mesin Solution X100C T43', '192.168.150.188', 80, 'Pos T43'),
('t44', 'Mesin Solution X100C T44', '192.168.150.189', 80, 'Pos T44'),
('t61', 'Mesin Solution X100C T61', '192.168.150.173', 80, 'Pos T61'),
('t62', 'Mesin Solution X100C T62', '192.168.150.174', 80, 'Pos T62'),
('t63', 'Mesin Solution X100C T63', '192.168.150.175', 80, 'Pos T63'),
('t64', 'Mesin Solution X100C T64', '192.168.150.176', 80, 'Pos T64'),
('t65', 'Mesin Solution X100C T65', '192.168.150.177', 80, 'Pos T65'),
('t66', 'Mesin Solution X100C T66', '192.168.150.178', 80, 'Pos T66')
ON CONFLICT (code) DO NOTHING;
