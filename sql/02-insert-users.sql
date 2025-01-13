INSERT INTO users (logon_name, full_name, email)
SELECT * FROM (
                  VALUES
      ('mike1', 'mike', 'mike@email.com'),
      ('bob44', 'bob', 'bob@email.com'),
      ('sarah485', 'sarah', 'sarah@email.com'),
      ('eric2', 'eric', 'eric@email.com'),
      ('susan9', 'susan', 'susan@email.com'),
      ('holly0', 'holly', 'holly@email.com'),
      ('bobby8', 'bobby', 'bobby@email.com'),
      ('clive88', 'clive', 'clive@email.com'),
      ('lorna1', 'lorna', 'lorna@email.com'),
      ('jayne2234', 'jayne', 'jayne@email.com')
) AS new_data(logon_name, full_name, email)
WHERE NOT EXISTS (
    SELECT 1 FROM users
);