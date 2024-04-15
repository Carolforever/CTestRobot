-- SQLite
SELECT DISTINCT name, c FROM packages WHERE header > 0 AND c > ("all" * 0.4) AND c < 250000