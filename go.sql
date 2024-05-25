SELECT m.id,
      u.name, u.email, u.created_at,
      mc.id AS match_cat_id, mc.name, mc.race, mc.sex, mc.description, mc.age_in_month, mc.image_urls, mc.has_matched, mc.created_at,
      uc.id AS user_cat_id, uc.name, uc.race, uc.sex, uc.description, uc.age_in_month, uc.image_urls, uc.has_matched, uc.created_at,
      m.message, m.created_at
FROM match_cats m
JOIN cats mc ON m.match_cat_id = mc.id
JOIN cats uc ON m.user_cat_id = uc.id
JOIN users u ON m.issued_by = u.id;