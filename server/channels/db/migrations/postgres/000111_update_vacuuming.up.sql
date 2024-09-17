DO $$
DECLARE
 version text;
BEGIN
SELECT version() INTO version;
 -- YugabyteDB has no vaccuum
IF version !~ '-YB-' THEN
 execute $DDL$
alter table posts set (autovacuum_vacuum_scale_factor = 0.1, autovacuum_analyze_scale_factor = 0.05);
alter table threadmemberships set (autovacuum_vacuum_scale_factor = 0.1, autovacuum_analyze_scale_factor = 0.05);
alter table fileinfo set (autovacuum_vacuum_scale_factor = 0.1, autovacuum_analyze_scale_factor = 0.05);
alter table preferences set (autovacuum_vacuum_scale_factor = 0.1, autovacuum_analyze_scale_factor = 0.05);
 $DDL$;
END IF;
END $$;
