CREATE index idx_users_id ON users (id);
CREATE index idx_users_email ON users (email);
CREATE index idx_cats_id ON cats (id);
CREATE index idx_cats_owner_id ON cats (owner_id);
CREATE index idx_match_cats_id ON match_cats (id);
CREATE index idx_match_cats_match_cat_id ON match_cats (match_cat_id);
CREATE index idx_match_cats_user_cat_id ON match_cats (user_cat_id);
