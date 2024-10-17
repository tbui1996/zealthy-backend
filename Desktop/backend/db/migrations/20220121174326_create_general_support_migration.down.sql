DELETE FROM feature_flags.flags f
WHERE f.key = 'isGeneralSupportEnabled' AND f.name = 'General Support';