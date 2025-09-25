-- Create bike sales data table
CREATE TABLE IF NOT EXISTS bike_sales (
    id SERIAL PRIMARY KEY,
    sale_date DATE,
    quarter VARCHAR(10),
    month_year VARCHAR(10),
    bike_model VARCHAR(100),
    bike_category VARCHAR(50),
    price DECIMAL(10,2),
    quantity INTEGER,
    total_revenue DECIMAL(12,2),
    customer_type VARCHAR(50),
    region VARCHAR(50),
    store_id VARCHAR(20)
);

-- Create monthly active users table
CREATE TABLE IF NOT EXISTS monthly_active_users (
    id SERIAL PRIMARY KEY,
    month_year VARCHAR(10),
    active_users INTEGER,
    new_users INTEGER,
    returning_users INTEGER,
    region VARCHAR(50),
    platform VARCHAR(50),
    recorded_date DATE
);

-- Insert quarterly bike sales data for 2024
INSERT INTO bike_sales (sale_date, quarter, month_year, bike_model, bike_category, price, quantity, total_revenue, customer_type, region, store_id) VALUES
-- Q1 2024 - January
('2024-01-15', 'Q1 2024', '2024-01', 'Mountain Explorer Pro', 'Mountain', 1299.99, 25, 32499.75, 'Individual', 'North', 'ST001'),
('2024-01-18', 'Q1 2024', '2024-01', 'City Cruiser Deluxe', 'Urban', 849.99, 35, 29749.65, 'Individual', 'South', 'ST002'),
('2024-01-22', 'Q1 2024', '2024-01', 'Road Racer Elite', 'Road', 1899.99, 15, 28499.85, 'Individual', 'East', 'ST003'),
('2024-01-28', 'Q1 2024', '2024-01', 'Electric Commuter', 'Electric', 2199.99, 20, 43999.80, 'Corporate', 'West', 'ST004'),

-- Q1 2024 - February
('2024-02-05', 'Q1 2024', '2024-02', 'Mountain Explorer Pro', 'Mountain', 1299.99, 30, 38999.70, 'Individual', 'North', 'ST001'),
('2024-02-12', 'Q1 2024', '2024-02', 'City Cruiser Deluxe', 'Urban', 849.99, 42, 35699.58, 'Individual', 'South', 'ST002'),
('2024-02-18', 'Q1 2024', '2024-02', 'Road Racer Elite', 'Road', 1899.99, 18, 34199.82, 'Individual', 'East', 'ST003'),
('2024-02-24', 'Q1 2024', '2024-02', 'Electric Commuter', 'Electric', 2199.99, 28, 61599.72, 'Corporate', 'West', 'ST004'),

-- Q1 2024 - March
('2024-03-08', 'Q1 2024', '2024-03', 'Mountain Explorer Pro', 'Mountain', 1299.99, 38, 49399.62, 'Individual', 'North', 'ST001'),
('2024-03-15', 'Q1 2024', '2024-03', 'City Cruiser Deluxe', 'Urban', 849.99, 45, 38249.55, 'Individual', 'South', 'ST002'),
('2024-03-22', 'Q1 2024', '2024-03', 'Road Racer Elite', 'Road', 1899.99, 22, 41799.78, 'Individual', 'East', 'ST003'),
('2024-03-28', 'Q1 2024', '2024-03', 'Electric Commuter', 'Electric', 2199.99, 32, 70399.68, 'Corporate', 'West', 'ST004'),

-- Q2 2024 - April
('2024-04-10', 'Q2 2024', '2024-04', 'Mountain Explorer Pro', 'Mountain', 1299.99, 45, 58499.55, 'Individual', 'North', 'ST001'),
('2024-04-15', 'Q2 2024', '2024-04', 'City Cruiser Deluxe', 'Urban', 849.99, 52, 44199.48, 'Individual', 'South', 'ST002'),
('2024-04-20', 'Q2 2024', '2024-04', 'Road Racer Elite', 'Road', 1899.99, 28, 53199.72, 'Individual', 'East', 'ST003'),
('2024-04-25', 'Q2 2024', '2024-04', 'Electric Commuter', 'Electric', 2199.99, 38, 83599.62, 'Corporate', 'West', 'ST004'),

-- Q2 2024 - May
('2024-05-08', 'Q2 2024', '2024-05', 'Mountain Explorer Pro', 'Mountain', 1299.99, 52, 67599.48, 'Individual', 'North', 'ST001'),
('2024-05-15', 'Q2 2024', '2024-05', 'City Cruiser Deluxe', 'Urban', 849.99, 58, 49299.42, 'Individual', 'South', 'ST002'),
('2024-05-22', 'Q2 2024', '2024-05', 'Road Racer Elite', 'Road', 1899.99, 35, 66499.65, 'Individual', 'East', 'ST003'),
('2024-05-28', 'Q2 2024', '2024-05', 'Electric Commuter', 'Electric', 2199.99, 45, 98999.55, 'Corporate', 'West', 'ST004'),

-- Q2 2024 - June
('2024-06-05', 'Q2 2024', '2024-06', 'Mountain Explorer Pro', 'Mountain', 1299.99, 60, 77999.40, 'Individual', 'North', 'ST001'),
('2024-06-12', 'Q2 2024', '2024-06', 'City Cruiser Deluxe', 'Urban', 849.99, 65, 55249.35, 'Individual', 'South', 'ST002'),
('2024-06-18', 'Q2 2024', '2024-06', 'Road Racer Elite', 'Road', 1899.99, 42, 79799.58, 'Individual', 'East', 'ST003'),
('2024-06-25', 'Q2 2024', '2024-06', 'Electric Commuter', 'Electric', 2199.99, 55, 120999.45, 'Corporate', 'West', 'ST004'),

-- Q3 2024 - July
('2024-07-08', 'Q3 2024', '2024-07', 'Mountain Explorer Pro', 'Mountain', 1299.99, 68, 88399.32, 'Individual', 'North', 'ST001'),
('2024-07-15', 'Q3 2024', '2024-07', 'City Cruiser Deluxe', 'Urban', 849.99, 72, 61199.28, 'Individual', 'South', 'ST002'),
('2024-07-20', 'Q3 2024', '2024-07', 'Road Racer Elite', 'Road', 1899.99, 48, 91199.52, 'Individual', 'East', 'ST003'),
('2024-07-28', 'Q3 2024', '2024-07', 'Electric Commuter', 'Electric', 2199.99, 62, 136399.38, 'Corporate', 'West', 'ST004'),

-- Q3 2024 - August
('2024-08-05', 'Q3 2024', '2024-08', 'Mountain Explorer Pro', 'Mountain', 1299.99, 75, 97499.25, 'Individual', 'North', 'ST001'),
('2024-08-12', 'Q3 2024', '2024-08', 'City Cruiser Deluxe', 'Urban', 849.99, 78, 66299.22, 'Individual', 'South', 'ST002'),
('2024-08-18', 'Q3 2024', '2024-08', 'Road Racer Elite', 'Road', 1899.99, 55, 104499.45, 'Individual', 'East', 'ST003'),
('2024-08-25', 'Q3 2024', '2024-08', 'Electric Commuter', 'Electric', 2199.99, 68, 149599.32, 'Corporate', 'West', 'ST004'),

-- Q3 2024 - September
('2024-09-08', 'Q3 2024', '2024-09', 'Mountain Explorer Pro', 'Mountain', 1299.99, 82, 106599.18, 'Individual', 'North', 'ST001'),
('2024-09-15', 'Q3 2024', '2024-09', 'City Cruiser Deluxe', 'Urban', 849.99, 85, 72249.15, 'Individual', 'South', 'ST002'),
('2024-09-20', 'Q3 2024', '2024-09', 'Road Racer Elite', 'Road', 1899.99, 62, 117799.38, 'Individual', 'East', 'ST003'),
('2024-09-25', 'Q3 2024', '2024-09', 'Electric Commuter', 'Electric', 2199.99, 75, 164999.25, 'Corporate', 'West', 'ST004');

-- Insert monthly active users data
INSERT INTO monthly_active_users (month_year, active_users, new_users, returning_users, region, platform, recorded_date) VALUES
-- 2024 data showing growth trend
('2024-01', 12500, 2800, 9700, 'North America', 'Web', '2024-01-31'),
('2024-01', 8200, 1900, 6300, 'North America', 'Mobile', '2024-01-31'),
('2024-01', 6800, 1500, 5300, 'Europe', 'Web', '2024-01-31'),
('2024-01', 4200, 950, 3250, 'Europe', 'Mobile', '2024-01-31'),

('2024-02', 13800, 3200, 10600, 'North America', 'Web', '2024-02-29'),
('2024-02', 9100, 2100, 7000, 'North America', 'Mobile', '2024-02-29'),
('2024-02', 7500, 1700, 5800, 'Europe', 'Web', '2024-02-29'),
('2024-02', 4800, 1100, 3700, 'Europe', 'Mobile', '2024-02-29'),

('2024-03', 15200, 3600, 11600, 'North America', 'Web', '2024-03-31'),
('2024-03', 10300, 2400, 7900, 'North America', 'Mobile', '2024-03-31'),
('2024-03', 8400, 1900, 6500, 'Europe', 'Web', '2024-03-31'),
('2024-03', 5600, 1300, 4300, 'Europe', 'Mobile', '2024-03-31'),

('2024-04', 16800, 4000, 12800, 'North America', 'Web', '2024-04-30'),
('2024-04', 11800, 2700, 9100, 'North America', 'Mobile', '2024-04-30'),
('2024-04', 9200, 2100, 7100, 'Europe', 'Web', '2024-04-30'),
('2024-04', 6200, 1450, 4750, 'Europe', 'Mobile', '2024-04-30'),

('2024-05', 18500, 4500, 14000, 'North America', 'Web', '2024-05-31'),
('2024-05', 13200, 3000, 10200, 'North America', 'Mobile', '2024-05-31'),
('2024-05', 10100, 2300, 7800, 'Europe', 'Web', '2024-05-31'),
('2024-05', 7000, 1600, 5400, 'Europe', 'Mobile', '2024-05-31'),

('2024-06', 20400, 5000, 15400, 'North America', 'Web', '2024-06-30'),
('2024-06', 14800, 3400, 11400, 'North America', 'Mobile', '2024-06-30'),
('2024-06', 11200, 2600, 8600, 'Europe', 'Web', '2024-06-30'),
('2024-06', 7900, 1800, 6100, 'Europe', 'Mobile', '2024-06-30'),

('2024-07', 22600, 5600, 17000, 'North America', 'Web', '2024-07-31'),
('2024-07', 16800, 3900, 12900, 'North America', 'Mobile', '2024-07-31'),
('2024-07', 12500, 2900, 9600, 'Europe', 'Web', '2024-07-31'),
('2024-07', 8900, 2050, 6850, 'Europe', 'Mobile', '2024-07-31'),

('2024-08', 25200, 6300, 18900, 'North America', 'Web', '2024-08-31'),
('2024-08', 19200, 4400, 14800, 'North America', 'Mobile', '2024-08-31'),
('2024-08', 14100, 3250, 10850, 'Europe', 'Web', '2024-08-31'),
('2024-08', 10200, 2350, 7850, 'Europe', 'Mobile', '2024-08-31'),

('2024-09', 27800, 7000, 20800, 'North America', 'Web', '2024-09-30'),
('2024-09', 21800, 5000, 16800, 'North America', 'Mobile', '2024-09-30'),
('2024-09', 15900, 3700, 12200, 'Europe', 'Web', '2024-09-30'),
('2024-09', 11700, 2700, 9000, 'Europe', 'Mobile', '2024-09-30');