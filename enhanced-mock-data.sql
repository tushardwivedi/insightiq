-- Enhanced mock data for comprehensive analytics queries
-- Supporting customer segments, demographics, weekly trends, and marketing campaigns

-- Create customer demographics table
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(200),
    age INTEGER,
    gender VARCHAR(20),
    income_bracket VARCHAR(50),
    education_level VARCHAR(50),
    occupation VARCHAR(100),
    city VARCHAR(100),
    state VARCHAR(50),
    country VARCHAR(50),
    registration_date DATE,
    customer_segment VARCHAR(50),
    lifetime_value DECIMAL(12,2),
    last_purchase_date DATE
);

-- Create customer segments performance table
CREATE TABLE IF NOT EXISTS customer_segments (
    id SERIAL PRIMARY KEY,
    segment_name VARCHAR(50),
    segment_description TEXT,
    total_customers INTEGER,
    avg_order_value DECIMAL(10,2),
    purchase_frequency DECIMAL(5,2),
    customer_lifetime_value DECIMAL(12,2),
    churn_rate DECIMAL(5,2),
    acquisition_cost DECIMAL(10,2),
    month_year VARCHAR(10),
    updated_date DATE
);

-- Create marketing campaigns table
CREATE TABLE IF NOT EXISTS marketing_campaigns (
    id SERIAL PRIMARY KEY,
    campaign_id VARCHAR(50) UNIQUE NOT NULL,
    campaign_name VARCHAR(200),
    campaign_type VARCHAR(50),
    channel VARCHAR(50),
    start_date DATE,
    end_date DATE,
    budget DECIMAL(12,2),
    spend DECIMAL(12,2),
    impressions INTEGER,
    clicks INTEGER,
    conversions INTEGER,
    revenue_generated DECIMAL(12,2),
    target_segment VARCHAR(50),
    status VARCHAR(20),
    week_ending DATE
);

-- Create weekly sales trends table
CREATE TABLE IF NOT EXISTS weekly_sales_trends (
    id SERIAL PRIMARY KEY,
    week_ending DATE,
    week_number INTEGER,
    year INTEGER,
    total_orders INTEGER,
    total_revenue DECIMAL(12,2),
    avg_order_value DECIMAL(10,2),
    new_customers INTEGER,
    returning_customers INTEGER,
    top_selling_category VARCHAR(50),
    region VARCHAR(50)
);

-- Create product performance table
CREATE TABLE IF NOT EXISTS product_performance (
    id SERIAL PRIMARY KEY,
    product_id VARCHAR(50),
    product_name VARCHAR(200),
    category VARCHAR(50),
    units_sold INTEGER,
    revenue DECIMAL(12,2),
    profit_margin DECIMAL(5,2),
    inventory_turnover DECIMAL(5,2),
    customer_rating DECIMAL(3,2),
    return_rate DECIMAL(5,2),
    week_ending DATE,
    region VARCHAR(50)
);

-- Insert customer demographics data
INSERT INTO customers (customer_id, first_name, last_name, email, age, gender, income_bracket, education_level, occupation, city, state, country, registration_date, customer_segment, lifetime_value, last_purchase_date) VALUES
-- Premium Segment
('CUST001', 'Sarah', 'Johnson', 'sarah.j@email.com', 34, 'Female', '$75K-$100K', 'Masters', 'Software Engineer', 'San Francisco', 'CA', 'USA', '2023-01-15', 'Premium', 4850.00, '2024-09-20'),
('CUST002', 'Michael', 'Chen', 'michael.c@email.com', 42, 'Male', '$100K+', 'Bachelors', 'Marketing Director', 'Seattle', 'WA', 'USA', '2023-02-20', 'Premium', 6200.00, '2024-09-18'),
('CUST003', 'Emily', 'Davis', 'emily.d@email.com', 29, 'Female', '$75K-$100K', 'Masters', 'Product Manager', 'Austin', 'TX', 'USA', '2023-03-10', 'Premium', 3900.00, '2024-09-15'),
('CUST004', 'James', 'Wilson', 'james.w@email.com', 38, 'Male', '$100K+', 'PhD', 'Data Scientist', 'Boston', 'MA', 'USA', '2023-01-25', 'Premium', 5500.00, '2024-09-22'),
('CUST005', 'Lisa', 'Martinez', 'lisa.m@email.com', 31, 'Female', '$75K-$100K', 'Masters', 'UX Designer', 'Portland', 'OR', 'USA', '2023-04-12', 'Premium', 4100.00, '2024-09-19'),

-- Standard Segment
('CUST006', 'David', 'Brown', 'david.b@email.com', 45, 'Male', '$50K-$75K', 'Bachelors', 'Operations Manager', 'Denver', 'CO', 'USA', '2023-05-20', 'Standard', 2800.00, '2024-09-10'),
('CUST007', 'Jennifer', 'Garcia', 'jennifer.g@email.com', 33, 'Female', '$50K-$75K', 'Bachelors', 'Teacher', 'Phoenix', 'AZ', 'USA', '2023-06-15', 'Standard', 2200.00, '2024-09-12'),
('CUST008', 'Robert', 'Taylor', 'robert.t@email.com', 39, 'Male', '$50K-$75K', 'Associates', 'Sales Representative', 'Miami', 'FL', 'USA', '2023-07-08', 'Standard', 2650.00, '2024-09-14'),
('CUST009', 'Amanda', 'Anderson', 'amanda.a@email.com', 27, 'Female', '$50K-$75K', 'Bachelors', 'Nurse', 'Chicago', 'IL', 'USA', '2023-08-22', 'Standard', 1950.00, '2024-09-11'),
('CUST010', 'Christopher', 'Moore', 'chris.m@email.com', 35, 'Male', '$50K-$75K', 'Bachelors', 'Accountant', 'Atlanta', 'GA', 'USA', '2023-09-18', 'Standard', 2350.00, '2024-09-13'),

-- Budget Segment
('CUST011', 'Jessica', 'Thompson', 'jessica.t@email.com', 22, 'Female', '$25K-$50K', 'Bachelors', 'Student', 'College Station', 'TX', 'USA', '2023-10-05', 'Budget', 850.00, '2024-09-08'),
('CUST012', 'Daniel', 'White', 'daniel.w@email.com', 26, 'Male', '$25K-$50K', 'High School', 'Retail Associate', 'Sacramento', 'CA', 'USA', '2023-11-12', 'Budget', 650.00, '2024-09-05'),
('CUST013', 'Ashley', 'Harris', 'ashley.h@email.com', 24, 'Female', '$25K-$50K', 'Associates', 'Administrative Assistant', 'Orlando', 'FL', 'USA', '2023-12-20', 'Budget', 750.00, '2024-09-07'),
('CUST014', 'Matthew', 'Clark', 'matthew.c@email.com', 28, 'Male', '$25K-$50K', 'Bachelors', 'Customer Service Rep', 'Kansas City', 'MO', 'USA', '2024-01-08', 'Budget', 920.00, '2024-09-09'),
('CUST015', 'Brittany', 'Lewis', 'brittany.l@email.com', 23, 'Female', '<$25K', 'High School', 'Part-time Worker', 'Detroit', 'MI', 'USA', '2024-02-14', 'Budget', 480.00, '2024-09-03');

-- Insert customer segment performance data
INSERT INTO customer_segments (segment_name, segment_description, total_customers, avg_order_value, purchase_frequency, customer_lifetime_value, churn_rate, acquisition_cost, month_year, updated_date) VALUES
-- January 2024
('Premium', 'High-value customers with premium bikes and accessories', 1250, 1850.00, 3.2, 5920.00, 5.5, 145.00, '2024-01', '2024-01-31'),
('Standard', 'Mid-range customers with regular purchase patterns', 3200, 950.00, 2.1, 1995.00, 12.8, 85.00, '2024-01', '2024-01-31'),
('Budget', 'Price-sensitive customers seeking value deals', 2100, 420.00, 1.4, 588.00, 22.5, 35.00, '2024-01', '2024-01-31'),

-- February 2024
('Premium', 'High-value customers with premium bikes and accessories', 1320, 1920.00, 3.3, 6336.00, 4.8, 150.00, '2024-02', '2024-02-29'),
('Standard', 'Mid-range customers with regular purchase patterns', 3450, 985.00, 2.2, 2167.00, 11.2, 88.00, '2024-02', '2024-02-29'),
('Budget', 'Price-sensitive customers seeking value deals', 2280, 445.00, 1.5, 667.50, 21.0, 38.00, '2024-02', '2024-02-29'),

-- March 2024
('Premium', 'High-value customers with premium bikes and accessories', 1385, 1995.00, 3.4, 6783.00, 4.2, 155.00, '2024-03', '2024-03-31'),
('Standard', 'Mid-range customers with regular purchase patterns', 3680, 1015.00, 2.3, 2334.50, 10.5, 90.00, '2024-03', '2024-03-31'),
('Budget', 'Price-sensitive customers seeking value deals', 2420, 465.00, 1.6, 744.00, 19.8, 40.00, '2024-03', '2024-03-31'),

-- April 2024
('Premium', 'High-value customers with premium bikes and accessories', 1450, 2080.00, 3.5, 7280.00, 3.8, 160.00, '2024-04', '2024-04-30'),
('Standard', 'Mid-range customers with regular purchase patterns', 3920, 1045.00, 2.4, 2508.00, 9.8, 92.00, '2024-04', '2024-04-30'),
('Budget', 'Price-sensitive customers seeking value deals', 2580, 485.00, 1.7, 824.50, 18.5, 42.00, '2024-04', '2024-04-30'),

-- May 2024
('Premium', 'High-value customers with premium bikes and accessories', 1520, 2165.00, 3.6, 7794.00, 3.5, 165.00, '2024-05', '2024-05-31'),
('Standard', 'Mid-range customers with regular purchase patterns', 4180, 1075.00, 2.5, 2687.50, 9.2, 95.00, '2024-05', '2024-05-31'),
('Budget', 'Price-sensitive customers seeking value deals', 2750, 505.00, 1.8, 909.00, 17.2, 45.00, '2024-05', '2024-05-31'),

-- June 2024
('Premium', 'High-value customers with premium bikes and accessories', 1595, 2250.00, 3.7, 8325.00, 3.2, 170.00, '2024-06', '2024-06-30'),
('Standard', 'Mid-range customers with regular purchase patterns', 4450, 1105.00, 2.6, 2873.00, 8.8, 98.00, '2024-06', '2024-06-30'),
('Budget', 'Price-sensitive customers seeking value deals', 2920, 525.00, 1.9, 997.50, 16.5, 48.00, '2024-06', '2024-06-30'),

-- July 2024
('Premium', 'High-value customers with premium bikes and accessories', 1675, 2335.00, 3.8, 8873.00, 2.9, 175.00, '2024-07', '2024-07-31'),
('Standard', 'Mid-range customers with regular purchase patterns', 4720, 1135.00, 2.7, 3064.50, 8.4, 100.00, '2024-07', '2024-07-31'),
('Budget', 'Price-sensitive customers seeking value deals', 3100, 545.00, 2.0, 1090.00, 15.8, 50.00, '2024-07', '2024-07-31'),

-- August 2024
('Premium', 'High-value customers with premium bikes and accessories', 1760, 2420.00, 3.9, 9438.00, 2.6, 180.00, '2024-08', '2024-08-31'),
('Standard', 'Mid-range customers with regular purchase patterns', 5010, 1165.00, 2.8, 3262.00, 8.0, 105.00, '2024-08', '2024-08-31'),
('Budget', 'Price-sensitive customers seeking value deals', 3290, 565.00, 2.1, 1186.50, 15.1, 52.00, '2024-08', '2024-08-31'),

-- September 2024
('Premium', 'High-value customers with premium bikes and accessories', 1850, 2505.00, 4.0, 10020.00, 2.4, 185.00, '2024-09', '2024-09-30'),
('Standard', 'Mid-range customers with regular purchase patterns', 5320, 1195.00, 2.9, 3465.50, 7.6, 108.00, '2024-09', '2024-09-30'),
('Budget', 'Price-sensitive customers seeking value deals', 3480, 585.00, 2.2, 1287.00, 14.5, 55.00, '2024-09', '2024-09-30');

-- Insert marketing campaigns data
INSERT INTO marketing_campaigns (campaign_id, campaign_name, campaign_type, channel, start_date, end_date, budget, spend, impressions, clicks, conversions, revenue_generated, target_segment, status, week_ending) VALUES
-- Week ending 2024-01-07
('CAMP001', 'New Year Bike Sale', 'Promotional', 'Social Media', '2024-01-01', '2024-01-07', 15000.00, 14200.00, 285000, 8550, 342, 68400.00, 'All', 'Completed', '2024-01-07'),
('CAMP002', 'Premium Bike Launch', 'Product Launch', 'Email', '2024-01-01', '2024-01-14', 25000.00, 23800.00, 125000, 3750, 188, 94000.00, 'Premium', 'Completed', '2024-01-07'),

-- Week ending 2024-01-14
('CAMP003', 'Winter Commuter Special', 'Seasonal', 'Search Ads', '2024-01-08', '2024-01-21', 18000.00, 17500.00, 195000, 5850, 234, 58500.00, 'Standard', 'Completed', '2024-01-14'),
('CAMP004', 'Student Discount Program', 'Demographic', 'Social Media', '2024-01-08', '2024-02-08', 8000.00, 7200.00, 145000, 4350, 174, 26100.00, 'Budget', 'Completed', '2024-01-14'),

-- Week ending 2024-01-21
('CAMP005', 'Mountain Bike Adventure', 'Category Focus', 'Display Ads', '2024-01-15', '2024-01-28', 20000.00, 19200.00, 225000, 6750, 270, 81000.00, 'Premium', 'Completed', '2024-01-21'),

-- Week ending 2024-01-28
('CAMP006', 'Urban Commuter Campaign', 'Lifestyle', 'Influencer', '2024-01-22', '2024-02-05', 12000.00, 11400.00, 165000, 4950, 198, 39600.00, 'Standard', 'Completed', '2024-01-28'),

-- February campaigns
('CAMP007', 'Valentines Day Special', 'Holiday', 'Email', '2024-02-01', '2024-02-14', 10000.00, 9500.00, 120000, 3600, 144, 36000.00, 'All', 'Completed', '2024-02-04'),
('CAMP008', 'Electric Bike Showcase', 'Product Focus', 'Video Ads', '2024-02-05', '2024-02-18', 30000.00, 28500.00, 320000, 9600, 384, 153600.00, 'Premium', 'Completed', '2024-02-11'),
('CAMP009', 'Spring Prep Sale', 'Seasonal', 'Search Ads', '2024-02-12', '2024-02-25', 22000.00, 20900.00, 245000, 7350, 294, 73500.00, 'Standard', 'Completed', '2024-02-18'),
('CAMP010', 'Budget Bike Bonanza', 'Price Focus', 'Social Media', '2024-02-19', '2024-03-04', 9000.00, 8100.00, 155000, 4650, 186, 27900.00, 'Budget', 'Completed', '2024-02-25'),

-- March campaigns
('CAMP011', 'Spring Cycling Festival', 'Event', 'Multi-channel', '2024-03-01', '2024-03-15', 35000.00, 33200.00, 385000, 11550, 462, 184800.00, 'All', 'Completed', '2024-03-03'),
('CAMP012', 'Professional Cyclist Endorsement', 'Influencer', 'Social Media', '2024-03-04', '2024-03-17', 28000.00, 26600.00, 295000, 8850, 354, 141600.00, 'Premium', 'Completed', '2024-03-10'),
('CAMP013', 'Family Bike Package', 'Bundle', 'Email', '2024-03-11', '2024-03-24', 16000.00, 15200.00, 185000, 5550, 222, 55500.00, 'Standard', 'Completed', '2024-03-17'),
('CAMP014', 'College Campus Tour', 'Demographic', 'Campus Events', '2024-03-18', '2024-03-31', 12000.00, 10800.00, 95000, 2850, 114, 17100.00, 'Budget', 'Completed', '2024-03-24'),

-- April campaigns
('CAMP015', 'Earth Day Eco Bikes', 'Cause Marketing', 'Multi-channel', '2024-04-01', '2024-04-22', 25000.00, 23750.00, 275000, 8250, 330, 132000.00, 'Premium', 'Completed', '2024-04-07'),
('CAMP016', 'Spring Training Special', 'Seasonal', 'Sports Media', '2024-04-08', '2024-04-21', 18000.00, 17100.00, 205000, 6150, 246, 61500.00, 'Standard', 'Completed', '2024-04-14'),
('CAMP017', 'First Time Buyer Program', 'Acquisition', 'Search Ads', '2024-04-15', '2024-04-28', 14000.00, 13300.00, 175000, 5250, 210, 31500.00, 'Budget', 'Completed', '2024-04-21'),

-- May campaigns
('CAMP018', 'Mother\'s Day Gift Guide', 'Holiday', 'Email', '2024-05-01', '2024-05-12', 15000.00, 14250.00, 155000, 4650, 186, 46500.00, 'All', 'Completed', '2024-05-05'),
('CAMP019', 'Premium Accessories Launch', 'Product Launch', 'Display Ads', '2024-05-06', '2024-05-19', 32000.00, 30400.00, 345000, 10350, 414, 165600.00, 'Premium', 'Completed', '2024-05-12'),
('CAMP020', 'Memorial Day Weekend Sale', 'Holiday', 'Multi-channel', '2024-05-20', '2024-05-27', 28000.00, 26600.00, 315000, 9450, 378, 113400.00, 'Standard', 'Completed', '2024-05-26'),

-- June campaigns
('CAMP021', 'Summer Kickoff Campaign', 'Seasonal', 'Video Ads', '2024-06-01', '2024-06-15', 40000.00, 38000.00, 425000, 12750, 510, 204000.00, 'All', 'Completed', '2024-06-02'),
('CAMP022', 'Father\'s Day Special', 'Holiday', 'Social Media', '2024-06-10', '2024-06-16', 20000.00, 19000.00, 235000, 7050, 282, 84600.00, 'Premium', 'Completed', '2024-06-16'),
('CAMP023', 'Graduation Gift Campaign', 'Demographic', 'Multi-channel', '2024-06-17', '2024-06-30', 18000.00, 17100.00, 195000, 5850, 234, 35100.00, 'Budget', 'Completed', '2024-06-23'),

-- July campaigns
('CAMP024', 'Mid-Summer Clearance', 'Promotional', 'Email', '2024-07-01', '2024-07-15', 22000.00, 20900.00, 255000, 7650, 306, 91800.00, 'Standard', 'Completed', '2024-07-07'),
('CAMP025', 'Professional Racing Series', 'Sports Marketing', 'Sports Media', '2024-07-08', '2024-07-28', 45000.00, 42750.00, 485000, 14550, 582, 232800.00, 'Premium', 'Completed', '2024-07-14'),
('CAMP026', 'Summer Adventure Package', 'Bundle', 'Influencer', '2024-07-15', '2024-07-31', 25000.00, 23750.00, 285000, 8550, 342, 102600.00, 'All', 'Completed', '2024-07-21'),

-- August campaigns
('CAMP027', 'Back to School Commuter', 'Seasonal', 'Campus Media', '2024-08-01', '2024-08-25', 30000.00, 28500.00, 335000, 10050, 402, 60300.00, 'Budget', 'Completed', '2024-08-04'),
('CAMP028', 'Late Summer Premium Push', 'Category Focus', 'Multi-channel', '2024-08-12', '2024-08-31', 38000.00, 36100.00, 415000, 12450, 498, 199200.00, 'Premium', 'Completed', '2024-08-18'),
('CAMP029', 'End of Summer Sale', 'Promotional', 'Social Media', '2024-08-26', '2024-09-08', 26000.00, 24700.00, 295000, 8850, 354, 88500.00, 'Standard', 'Completed', '2024-09-01'),

-- September campaigns
('CAMP030', 'Fall Preparation Campaign', 'Seasonal', 'Email', '2024-09-01', '2024-09-15', 35000.00, 33250.00, 385000, 11550, 462, 138600.00, 'All', 'Completed', '2024-09-08'),
('CAMP031', 'Harvest Festival Sponsorship', 'Event Marketing', 'Local Events', '2024-09-16', '2024-09-30', 20000.00, 19000.00, 165000, 4950, 198, 59400.00, 'Standard', 'Completed', '2024-09-22');

-- Insert weekly sales trends data
INSERT INTO weekly_sales_trends (week_ending, week_number, year, total_orders, total_revenue, avg_order_value, new_customers, returning_customers, top_selling_category, region) VALUES
-- January 2024 weeks
('2024-01-07', 1, 2024, 145, 187350.00, 1292.07, 58, 87, 'Electric', 'West'),
('2024-01-14', 2, 2024, 162, 209730.00, 1294.63, 65, 97, 'Mountain', 'North'),
('2024-01-21', 3, 2024, 178, 229140.00, 1287.42, 71, 107, 'Urban', 'South'),
('2024-01-28', 4, 2024, 195, 251175.00, 1288.08, 78, 117, 'Road', 'East'),

-- February 2024 weeks
('2024-02-04', 5, 2024, 210, 270300.00, 1287.14, 84, 126, 'Electric', 'West'),
('2024-02-11', 6, 2024, 225, 289125.00, 1284.00, 90, 135, 'Mountain', 'North'),
('2024-02-18', 7, 2024, 238, 305580.00, 1283.78, 95, 143, 'Urban', 'South'),
('2024-02-25', 8, 2024, 252, 323568.00, 1284.00, 101, 151, 'Road', 'East'),

-- March 2024 weeks
('2024-03-03', 9, 2024, 268, 343904.00, 1283.24, 107, 161, 'Electric', 'West'),
('2024-03-10', 10, 2024, 285, 365175.00, 1281.32, 114, 171, 'Mountain', 'North'),
('2024-03-17', 11, 2024, 302, 387072.00, 1281.70, 121, 181, 'Urban', 'South'),
('2024-03-24', 12, 2024, 318, 407436.00, 1281.32, 127, 191, 'Road', 'East'),
('2024-03-31', 13, 2024, 335, 428795.00, 1279.99, 134, 201, 'Electric', 'West'),

-- April 2024 weeks
('2024-04-07', 14, 2024, 352, 450624.00, 1280.00, 141, 211, 'Mountain', 'North'),
('2024-04-14', 15, 2024, 368, 471104.00, 1280.17, 147, 221, 'Urban', 'South'),
('2024-04-21', 16, 2024, 385, 492675.00, 1279.55, 154, 231, 'Road', 'East'),
('2024-04-28', 17, 2024, 402, 514056.00, 1278.50, 161, 241, 'Electric', 'West'),

-- May 2024 weeks
('2024-05-05', 18, 2024, 418, 534334.00, 1278.40, 167, 251, 'Mountain', 'North'),
('2024-05-12', 19, 2024, 435, 555705.00, 1277.48, 174, 261, 'Urban', 'South'),
('2024-05-19', 20, 2024, 452, 577204.00, 1277.00, 181, 271, 'Road', 'East'),
('2024-05-26', 21, 2024, 468, 597624.00, 1276.99, 187, 281, 'Electric', 'West'),

-- June 2024 weeks
('2024-06-02', 22, 2024, 485, 619015.00, 1276.32, 194, 291, 'Mountain', 'North'),
('2024-06-09', 23, 2024, 502, 640512.00, 1275.52, 201, 301, 'Urban', 'South'),
('2024-06-16', 24, 2024, 518, 661054.00, 1276.37, 207, 311, 'Road', 'East'),
('2024-06-23', 25, 2024, 535, 682045.00, 1274.85, 214, 321, 'Electric', 'West'),
('2024-06-30', 26, 2024, 552, 703536.00, 1274.52, 221, 331, 'Mountain', 'North'),

-- July 2024 weeks
('2024-07-07', 27, 2024, 568, 724544.00, 1275.43, 227, 341, 'Urban', 'South'),
('2024-07-14', 28, 2024, 585, 746025.00, 1275.17, 234, 351, 'Road', 'East'),
('2024-07-21', 29, 2024, 602, 767552.00, 1275.17, 241, 361, 'Electric', 'West'),
('2024-07-28', 30, 2024, 618, 788574.00, 1276.04, 247, 371, 'Mountain', 'North'),

-- August 2024 weeks
('2024-08-04', 31, 2024, 635, 810175.00, 1275.75, 254, 381, 'Urban', 'South'),
('2024-08-11', 32, 2024, 652, 832004.00, 1276.00, 261, 391, 'Road', 'East'),
('2024-08-18', 33, 2024, 668, 853132.00, 1277.04, 267, 401, 'Electric', 'West'),
('2024-08-25', 34, 2024, 685, 875275.00, 1277.74, 274, 411, 'Mountain', 'North'),

-- September 2024 weeks
('2024-09-01', 35, 2024, 702, 897054.00, 1277.85, 281, 421, 'Urban', 'South'),
('2024-09-08', 36, 2024, 718, 918322.00, 1279.02, 287, 431, 'Road', 'East'),
('2024-09-15', 37, 2024, 735, 940725.00, 1279.83, 294, 441, 'Electric', 'West'),
('2024-09-22', 38, 2024, 752, 963504.00, 1281.25, 301, 451, 'Mountain', 'North'),
('2024-09-29', 39, 2024, 768, 985536.00, 1283.27, 307, 461, 'Urban', 'South');

-- Insert product performance data
INSERT INTO product_performance (product_id, product_name, category, units_sold, revenue, profit_margin, inventory_turnover, customer_rating, return_rate, week_ending, region) VALUES
-- Week ending 2024-01-07
('PROD001', 'Mountain Explorer Pro', 'Mountain', 25, 32499.75, 35.0, 4.2, 4.5, 2.1, '2024-01-07', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 35, 29749.65, 28.5, 5.8, 4.3, 3.2, '2024-01-07', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 15, 28499.85, 42.0, 3.1, 4.7, 1.8, '2024-01-07', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 20, 43999.80, 38.5, 2.9, 4.6, 2.5, '2024-01-07', 'West'),

-- Week ending 2024-02-04
('PROD001', 'Mountain Explorer Pro', 'Mountain', 30, 38999.70, 35.2, 4.3, 4.5, 2.0, '2024-02-04', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 42, 35699.58, 29.0, 6.0, 4.4, 3.1, '2024-02-04', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 18, 34199.82, 42.5, 3.2, 4.7, 1.7, '2024-02-04', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 28, 61599.72, 39.0, 3.0, 4.7, 2.4, '2024-02-04', 'West'),

-- Week ending 2024-03-03
('PROD001', 'Mountain Explorer Pro', 'Mountain', 38, 49399.62, 35.5, 4.4, 4.6, 1.9, '2024-03-03', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 45, 38249.55, 29.5, 6.2, 4.4, 3.0, '2024-03-03', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 22, 41799.78, 43.0, 3.3, 4.8, 1.6, '2024-03-03', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 32, 70399.68, 39.5, 3.1, 4.7, 2.3, '2024-03-03', 'West'),

-- Week ending 2024-04-07
('PROD001', 'Mountain Explorer Pro', 'Mountain', 45, 58499.55, 36.0, 4.5, 4.6, 1.8, '2024-04-07', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 52, 44199.48, 30.0, 6.4, 4.5, 2.9, '2024-04-07', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 28, 53199.72, 43.5, 3.4, 4.8, 1.5, '2024-04-07', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 38, 83599.62, 40.0, 3.2, 4.8, 2.2, '2024-04-07', 'West'),

-- Week ending 2024-05-05
('PROD001', 'Mountain Explorer Pro', 'Mountain', 52, 67599.48, 36.5, 4.6, 4.7, 1.7, '2024-05-05', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 58, 49299.42, 30.5, 6.6, 4.5, 2.8, '2024-05-05', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 35, 66499.65, 44.0, 3.5, 4.9, 1.4, '2024-05-05', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 45, 98999.55, 40.5, 3.3, 4.8, 2.1, '2024-05-05', 'West'),

-- Week ending 2024-06-02
('PROD001', 'Mountain Explorer Pro', 'Mountain', 60, 77999.40, 37.0, 4.7, 4.7, 1.6, '2024-06-02', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 65, 55249.35, 31.0, 6.8, 4.6, 2.7, '2024-06-02', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 42, 79799.58, 44.5, 3.6, 4.9, 1.3, '2024-06-02', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 55, 120999.45, 41.0, 3.4, 4.9, 2.0, '2024-06-02', 'West'),

-- Week ending 2024-07-07
('PROD001', 'Mountain Explorer Pro', 'Mountain', 68, 88399.32, 37.5, 4.8, 4.8, 1.5, '2024-07-07', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 72, 61199.28, 31.5, 7.0, 4.6, 2.6, '2024-07-07', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 48, 91199.52, 45.0, 3.7, 5.0, 1.2, '2024-07-07', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 62, 136399.38, 41.5, 3.5, 4.9, 1.9, '2024-07-07', 'West'),

-- Week ending 2024-08-04
('PROD001', 'Mountain Explorer Pro', 'Mountain', 75, 97499.25, 38.0, 4.9, 4.8, 1.4, '2024-08-04', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 78, 66299.22, 32.0, 7.2, 4.7, 2.5, '2024-08-04', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 55, 104499.45, 45.5, 3.8, 5.0, 1.1, '2024-08-04', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 68, 149599.32, 42.0, 3.6, 5.0, 1.8, '2024-08-04', 'West'),

-- Week ending 2024-09-01
('PROD001', 'Mountain Explorer Pro', 'Mountain', 82, 106599.18, 38.5, 5.0, 4.9, 1.3, '2024-09-01', 'North'),
('PROD002', 'City Cruiser Deluxe', 'Urban', 85, 72249.15, 32.5, 7.4, 4.7, 2.4, '2024-09-01', 'South'),
('PROD003', 'Road Racer Elite', 'Road', 62, 117799.38, 46.0, 3.9, 5.1, 1.0, '2024-09-01', 'East'),
('PROD004', 'Electric Commuter', 'Electric', 75, 164999.25, 42.5, 3.7, 5.0, 1.7, '2024-09-01', 'West');

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_customers_segment ON customers(customer_segment);
CREATE INDEX IF NOT EXISTS idx_customers_age ON customers(age);
CREATE INDEX IF NOT EXISTS idx_customers_income ON customers(income_bracket);
CREATE INDEX IF NOT EXISTS idx_customer_segments_month ON customer_segments(month_year);
CREATE INDEX IF NOT EXISTS idx_marketing_campaigns_week ON marketing_campaigns(week_ending);
CREATE INDEX IF NOT EXISTS idx_marketing_campaigns_segment ON marketing_campaigns(target_segment);
CREATE INDEX IF NOT EXISTS idx_weekly_trends_week ON weekly_sales_trends(week_ending);
CREATE INDEX IF NOT EXISTS idx_product_performance_week ON product_performance(week_ending);
CREATE INDEX IF NOT EXISTS idx_product_performance_category ON product_performance(category);

-- Create views for common analytics queries
CREATE OR REPLACE VIEW customer_segment_analysis AS
SELECT
    cs.segment_name,
    cs.month_year,
    cs.total_customers,
    cs.avg_order_value,
    cs.customer_lifetime_value,
    cs.churn_rate,
    cs.acquisition_cost,
    ROUND(cs.customer_lifetime_value / cs.acquisition_cost, 2) as ltv_cac_ratio,
    CASE
        WHEN cs.churn_rate < 10 THEN 'Low Risk'
        WHEN cs.churn_rate < 20 THEN 'Medium Risk'
        ELSE 'High Risk'
    END as churn_risk_level
FROM customer_segments cs
ORDER BY cs.month_year DESC, cs.customer_lifetime_value DESC;

CREATE OR REPLACE VIEW weekly_performance_trends AS
SELECT
    wst.week_ending,
    wst.total_orders,
    wst.total_revenue,
    wst.avg_order_value,
    wst.new_customers,
    wst.returning_customers,
    ROUND((wst.returning_customers * 100.0 / (wst.new_customers + wst.returning_customers)), 2) as retention_rate,
    wst.top_selling_category,
    LAG(wst.total_revenue) OVER (ORDER BY wst.week_ending) as prev_week_revenue,
    ROUND(((wst.total_revenue - LAG(wst.total_revenue) OVER (ORDER BY wst.week_ending)) * 100.0 /
           LAG(wst.total_revenue) OVER (ORDER BY wst.week_ending)), 2) as revenue_growth_pct
FROM weekly_sales_trends wst
ORDER BY wst.week_ending DESC;

CREATE OR REPLACE VIEW marketing_campaign_roi AS
SELECT
    mc.campaign_name,
    mc.campaign_type,
    mc.channel,
    mc.target_segment,
    mc.spend,
    mc.revenue_generated,
    mc.conversions,
    ROUND((mc.revenue_generated - mc.spend) / mc.spend * 100, 2) as roi_percentage,
    ROUND(mc.revenue_generated / mc.spend, 2) as roas_ratio,
    ROUND(mc.spend / mc.conversions, 2) as cost_per_conversion,
    mc.week_ending
FROM marketing_campaigns mc
WHERE mc.conversions > 0
ORDER BY roi_percentage DESC;

CREATE OR REPLACE VIEW demographic_insights AS
SELECT
    c.customer_segment,
    c.age_group,
    c.gender,
    c.income_bracket,
    c.education_level,
    COUNT(*) as customer_count,
    AVG(c.lifetime_value) as avg_ltv,
    MIN(c.registration_date) as first_registration,
    MAX(c.last_purchase_date) as latest_purchase
FROM (
    SELECT *,
        CASE
            WHEN age < 25 THEN '18-24'
            WHEN age < 35 THEN '25-34'
            WHEN age < 45 THEN '35-44'
            WHEN age < 55 THEN '45-54'
            ELSE '55+'
        END as age_group
    FROM customers
) c
GROUP BY c.customer_segment, c.age_group, c.gender, c.income_bracket, c.education_level
ORDER BY c.customer_segment, avg_ltv DESC;