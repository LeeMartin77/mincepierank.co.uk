CREATE TABLE IF NOT EXISTS admins (
    id text NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS categories (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    slug text NULL,
    label text NULL
  );

CREATE TABLE IF NOT EXISTS config (
    key text NOT NULL PRIMARY KEY,
    value text NULL
);

CREATE TABLE IF NOT EXISTS maker_pie_yearly
  (
      year int,
      makerid text,
      id text,
      displayname text,
      fresh boolean,
      labels text[],
      image_file text,
      web_link text,
      pack_count int,
      pack_price_in_pence int,
      validated boolean,
      oid uuid DEFAULT gen_random_uuid(),
      PRIMARY KEY (year, makerId, id)
  );

  CREATE TABLE IF NOT EXISTS maker_pie_ranking_yearly
  (
      year int,
      makerid text,
      pieid text,
      userid text,
      pastry int,
      filling int,
      topping int,
      looks int,
      value int,
      notes text,
      last_updated timestamp,
      PRIMARY KEY (year, makerid, pieid, userid)
  );

  CREATE TABLE IF NOT EXISTS user_pie_yearly
  (
      year int,
      id text,
      owner_userid text,
      maker text,
      location text,
      displayname text,
      fresh boolean,
      image_file text,
      web_link text,
      pack_count int,
      pack_price_in_pence int,
      clean boolean,
      PRIMARY KEY (year, id, owner_userid)
  );

  CREATE TABLE IF NOT EXISTS user_pie_ranking_yearly
  (
      year int,
      pieid text,
      userid text,
      pastry int,
      filling int,
      topping int,
      looks int,
      value int,
      notes text,
      last_updated timestamp,
      PRIMARY KEY (year, pieid, userid)
  );

  CREATE TABLE IF NOT EXISTS maker (
		id text PRIMARY KEY,
		name text,
		logo text,
		website text
	);

CREATE TABLE IF NOT EXISTS categories (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		slug text,
		label text
	);
    CREATE TABLE IF NOT EXISTS maker_pie_categories (
		maker_pie_oid uuid,
		category_id uuid,
		PRIMARY KEY (maker_pie_oid, category_id)
	);

    CREATE VIEW maker_pie_yearly_categories AS
		  SELECT
		    mpy.oid as oid,
		    jsonb_agg(jsonb_build_object('id', c.id, 'slug', c.slug, 'label', c.label)) as categories,
			ARRAY_AGG(c.slug) as category_slugs
		  from maker_pie_yearly mpy
		  inner join maker_pie_categories mpc on mpc.maker_pie_oid = mpy.oid
		  inner join categories c on mpc.category_id = c.id
		  group by mpy.oid;