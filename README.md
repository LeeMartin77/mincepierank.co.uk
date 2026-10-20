# MincePieRank.co.uk

Because someone has to try and figure out who has the best mince pies.

## Container Tags:

Other than versions, the intent is to have two tags:

- `latest` will be the latest published release
- `bleed` will be the latest build of main

## Development

```
docker compose up -d
```

this will bring up a local stack with some dummy data

then

```
cp -r data/images .localimages
```

to put the "generic" images in

then 

```
task watch
```

to run the application with gow

then, login as user1 from the dex config, go to /admin, and bulk upload the file at data/2024/bulk/bulk.zip

congrats, you have mincepierank setup for local dev with 2024 as the year, and 2024 data.

### User-Submittable Pies

Users can submit, rate, and manage their own mince pie entries. The feature is controlled via the `config` database table:

```sql
INSERT INTO config (key, value) VALUES ('userPiesEnabled', 'true');
```

Set to `'false'` (or omit the key) to disable. The system fails closed when the key is missing or unparseable.

Admin moderation routes:
- `/admin/userpies` — list all user pies
- `/admin/userpies/{year}/{id}/{ownerid}` — view/edit a single pie
- `/admin/userpies/duplicates` — find duplicate submissions

### Previewing templates without the stack

For template/UI work you can skip all of the above:

```
task preview
```

and open http://127.0.0.1:8099/ — this renders pages with sample data and no database. See [cmd/template-preview/README.md](cmd/template-preview/README.md).
