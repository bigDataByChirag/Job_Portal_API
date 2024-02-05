CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    jobRole TEXT,
    salary INTEGER,
    companyId SERIAL,
    FOREIGN KEY (companyId) REFERENCES companies (id)
);
