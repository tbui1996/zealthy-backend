INSERT INTO feature_flags.flags
            (created_by,
             updated_by,
             key,
             name)
VALUES      ('SYSTEM',
             'SYSTEM',
             'isStarredSectionEnabled',
             'Starred Section')
ON CONFLICT DO NOTHING;