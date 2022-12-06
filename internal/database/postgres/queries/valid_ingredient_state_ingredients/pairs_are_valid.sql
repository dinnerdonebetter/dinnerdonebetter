SELECT
    COUNT(id) = $1
FROM
    valid_ingredient_preparations
WHERE
    (valid_ingredient_id, valid_preparation_id) IN (
    VALUES
    (
        '29e1aTUQlzwyjy5rMb5cQSybaWq',
        '29e2O2aRyQ7ALhd0OY6is3FSHsy'
    ),
    ('2DV0jyi0rejGKoMeAG3kFZKTGxP', '2DV0tMJlspdbTDL7xjReNQGPKyn')
);