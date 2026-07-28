-- Migration Seed: Master data for per-category tables (matches FE master-data.ts)

-- ============================================================================
-- 1. EGI TYPES
-- ============================================================================
INSERT INTO master_egi_types (code, name, is_active) VALUES
('HD785', 'HD785', TRUE),
('HD777', 'HD777', TRUE),
('HD405', 'HD405', TRUE),
('EX2600', 'EX2600', TRUE),
('EX2000', 'EX2000', TRUE),
('PC2000', 'PC2000', TRUE),
('PC200', 'PC200', TRUE),
('DT777', 'DT777', TRUE),
('DT785', 'DT785', TRUE),
('GR140', 'GR140', TRUE),
('WT773', 'WT773', TRUE),
('BUS', 'BUS', TRUE),
('SKT105', 'SKT105', TRUE),
('SPARE', 'SPARE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 2. PRODUCTS / BRANDS
-- ============================================================================
INSERT INTO master_products (code, name, is_active) VALUES
('Caterpillar', 'Caterpillar', TRUE),
('Komatsu', 'Komatsu', TRUE),
('Hitachi', 'Hitachi', TRUE),
('Liebherr', 'Liebherr', TRUE),
('Volvo', 'Volvo', TRUE),
('Scania', 'Scania', TRUE),
('Hino', 'Hino', TRUE),
('Isuzu', 'Isuzu', TRUE),
('Mitsubishi', 'Mitsubishi', TRUE),
('Toyota', 'Toyota', TRUE),
('UD Trucks', 'UD Trucks', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 3. EQUIPMENT CLASSES
-- ============================================================================
INSERT INTO master_eq_classes (code, name, description, is_active) VALUES
('HD 785 / 777', 'HD 785 / 777', 'Haul Truck 785 dan 777 series', TRUE),
('PC 2000', 'PC 2000', 'Excavator PC2000', TRUE),
('PC 2600', 'PC 2600', 'Excavator PC2600', TRUE),
('EX 2000', 'EX 2000', 'Excavator EX2000', TRUE),
('EX 2600', 'EX 2600', 'Excavator EX2600', TRUE),
('PC 200', 'PC 200', 'Excavator PC200', TRUE),
('SANY SYZ 440', 'SANY SYZ 440', 'Haul Truck SANY', TRUE),
('WATER TRUCK', 'WATER TRUCK', 'Water Truck', TRUE),
('SPARE', 'SPARE', 'Spare unit', TRUE),
('SKT105', 'SKT105', 'Skid Steer Loader', TRUE),
('BUS', 'BUS', 'Bus antar jemput', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 4. AREAS
-- ============================================================================
INSERT INTO master_areas (code, name, category, is_active) VALUES
('Pit Utara', 'Pit Utara', 'Mining', TRUE),
('Pit Selatan', 'Pit Selatan', 'Mining', TRUE),
('Panel East Tengah', 'Panel East Tengah', 'Mining', TRUE),
('Pit Service', 'Pit Service', 'Mining', TRUE),
('Workshop Plant', 'Workshop Plant', 'Non-Mining', TRUE),
('Kantor Operation', 'Kantor Operation', 'Non-Mining', TRUE),
('Kantor SDI', 'Kantor SDI', 'Non-Mining', TRUE),
('KM 31', 'KM 31', 'Mining', TRUE),
('Hauling Road A', 'Hauling Road A', 'Mining', TRUE),
('Hauling Road B', 'Hauling Road B', 'Mining', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 5. TEMPUDO POINTS
-- ============================================================================
INSERT INTO master_tempudo (code, name, location, pickup_type, is_active) VALUES
('TP-01', 'TP-01', 'Workshop', 'Pickup', TRUE),
('TP-02', 'TP-02', 'Panel East Tengah', 'Pickup & Drop', TRUE),
('TP-03', 'TP-03', 'Kasturi Tengah', 'Pickup & Drop', TRUE),
('TP-04', 'TP-04', 'Parkiran T6', 'Pickup', FALSE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 6. BUSES
-- ============================================================================
INSERT INTO master_buses (code, name, egi_type, departure_time, is_active) VALUES
('UDBU001', 'UDBU001', 'BUS', '05:00', TRUE),
('UDBU002', 'UDBU002', 'BUS', '05:15', TRUE),
('UDBU003', 'UDBU003', 'BUS', '05:30', TRUE),
('UDBU004', 'UDBU004', 'BUS', '05:45', TRUE),
('UDBU005', 'UDBU005', 'BUS', '05:00', TRUE),
('UDBU006', 'UDBU006', 'BUS', '05:15', TRUE),
('UDBU007', 'UDBU007', 'BUS', '05:30', TRUE),
('UDBU008', 'UDBU008', 'BUS', '05:45', TRUE),
('UDBU009', 'UDBU009', 'BUS', '05:00', TRUE),
('UDBU010', 'UDBU010', 'BUS', '05:15', TRUE),
('UDBU011', 'UDBU011', 'BUS', '05:30', TRUE),
('UDBU012', 'UDBU012', 'BUS', '05:45', TRUE),
('UDBU013', 'UDBU013', 'BUS', '05:00', TRUE),
('UDBU014', 'UDBU014', 'BUS', '05:15', TRUE),
('UDBU015', 'UDBU015', 'BUS', '05:30', TRUE),
('UDBU016', 'UDBU016', 'BUS', '05:45', TRUE),
('UDBU017', 'UDBU017', 'BUS', '05:00', TRUE),
('UDBU018', 'UDBU018', 'BUS', '05:15', TRUE),
('UDBU019', 'UDBU019', 'BUS', '05:30', TRUE),
('UDBU020', 'UDBU020', 'BUS', '05:45', TRUE),
('UDBU021', 'UDBU021', 'BUS', '05:00', TRUE),
('UDBU022', 'UDBU022', 'BUS', '05:15', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 7. EXCAVATOR LOCATIONS
-- ============================================================================
INSERT INTO master_locations_ex (code, name, bus_code, tempudo_code, is_active) VALUES
('EX7007', 'EX7007', 'UD-BU06', 'TP-02', TRUE),
('EX6001', 'EX6001', 'UD-BU07', 'TP-02', TRUE),
('EX5001', 'EX5001', 'UD-BU07', 'TP-03', TRUE),
('EX7003', 'EX7003', 'UD-BU08', 'TP-01', TRUE),
('EX7004', 'EX7004', 'UD-BU06', 'TP-04', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 8. MESS
-- ============================================================================
INSERT INTO master_mess (code, name, block, is_active) VALUES
('Mess 31 Blok A', 'Mess 31', 'Blok A', TRUE),
('Mess 31 Blok C', 'Mess 31', 'Blok C', TRUE),
('Mess 31 Blok D', 'Mess 31', 'Blok D', TRUE),
('Mess KM 12 Blok B', 'Mess KM 12', 'Blok B', TRUE),
('Mess KM 12 Blok A', 'Mess KM 12', 'Blok A', FALSE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 9. RUNNING TEXTS
-- ============================================================================
INSERT INTO master_running_texts (code, name, target_display, text_color, is_active) VALUES
('rt-1', 'Utamakan keselamatan — patuhi batas kecepatan 40 km/jam di jalan hauling.', 'Semua kiosk', 'Cyan', TRUE),
('rt-2', 'Wajib P2H sebelum mengoperasikan unit.', 'Display Fleet', 'Oranye', TRUE),
('rt-3', 'Rapat P5M setiap pergantian shift di front masing-masing.', 'Display Attendance', 'Putih', TRUE),
('rt-4', 'Musim hujan: waspadai jalan licin di ramp Pit Tempudo.', 'Semua kiosk', 'Merah', TRUE)
ON CONFLICT (code) DO NOTHING;