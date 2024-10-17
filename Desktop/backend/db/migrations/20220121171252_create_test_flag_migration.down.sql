DELETE FROM feature_flags.flags f
WHERE f.key = 'testFlag' AND f.name = 'Test Flag';