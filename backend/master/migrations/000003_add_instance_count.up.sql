-- Add instance_count column to nodes table
ALTER TABLE nodes ADD COLUMN instance_count INTEGER DEFAULT 0;
