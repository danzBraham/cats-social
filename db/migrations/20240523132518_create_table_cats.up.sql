DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cat_races') THEN
    CREATE TYPE cat_races AS ENUM (
      'Persian',
			'Maine Coon',
			'Siamese',
			'Ragdoll',
			'Bengal',
			'Sphynx',
			'British Shorthair',
			'Abyssinian',
			'Scottish Fold',
			'Birman'
      );
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cat_sex') THEN
    CREATE TYPE cat_sex AS ENUM ('male', 'female');
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS cats (
  id VARCHAR(26) PRIMARY KEY NOT NULL,
  name VARCHAR(30) NOT NULL,
  race cat_races NOT NULL,
  sex cat_sex NOT NULL,
  age_in_month INT NOT NULL CHECK (age_in_month BETWEEN 1 AND 120082),
  descrption VARCHAR(200) NOT NULL,
  img_urls TEXT[] NOT NULL CHECK (array_length(img_urls, 1) >= 1),
  owner_id VARCHAR(26) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
