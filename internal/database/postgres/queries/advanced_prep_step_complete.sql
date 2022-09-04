UPDATE advanced_prep_steps SET completed_at = extract(epoch from NOW()) WHERE id = $1;
