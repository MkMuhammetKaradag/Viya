-- Active: 1753095725457@@127.0.0.1@5432@postgres
-- ./init/init.sql

ALTER DATABASE template1 REFRESH COLLATION VERSION;
CREATE DATABASE userdb;

CREATE DATABASE tripdb;

CREATE DATABASE orderdb;

CREATE DATABASE paymentdb;
CREATE DATABASE notificationdb;



-- Ensure required extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

