DELETE FROM feature_flags.flags f
WHERE f.key = 'isStarredSectionEnabled' AND f.name = 'Starred Section';