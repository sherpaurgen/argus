CREATE TABLE users (
                       user_id char(27) PRIMARY KEY,
                       fname varchar(50)  NOT NULL,
                       lname varchar(50)  NOT NULL,
                       email varchar(255) UNIQUE NOT NULL,
                       secret varchar(255) NOT NULL,
                       created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
                       updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
                       last_login timestamp with time zone,
                       is_active boolean DEFAULT true,
                       role varchar(20) DEFAULT 'user'
);

CREATE TABLE children (
                          child_id char(27) PRIMARY KEY,
                          user_id char(27) REFERENCES users(user_id),
                          name VARCHAR(100) NOT NULL,
                          date_of_birth DATE NOT NULL
);

CREATE TABLE sleep_patterns (
                                
                                sleep_id char(27) PRIMARY KEY,
                                child_id char(27) REFERENCES children(child_id),
                                sleep_start TIMESTAMP WITH TIME ZONE NOT NULL,
                                sleep_end TIMESTAMP WITH TIME ZONE NOT NULL,
                                sleep_quality INTEGER CHECK (sleep_quality BETWEEN 1 AND 5),
                                device_id char(27)
);

CREATE TABLE temperature_readings (
                                      id char(27) PRIMARY KEY,
                                      child_id char(27) REFERENCES children(child_id),
                                      device_id char(27),
                                      timestamp TIMESTAMP WITH TIME ZONE,
                                      temperature DECIMAL(5,2)
);
