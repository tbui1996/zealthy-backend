INSERT INTO feature_flags.flags
            (created_by,
             updated_by,
             key,
             name)
VALUES      ('SYSTEM',
             'SYSTEM',
             'testFlag',
             'Test Flag')
ON CONFLICT DO NOTHING;