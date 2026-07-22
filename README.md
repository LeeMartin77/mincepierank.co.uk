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