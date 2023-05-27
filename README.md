# MincePieRank.co.uk

Because someone has to try and figure out who has the best mince pies.

## Container Tags:

Other than versions, the intent is to have two tags:

- `latest` will be the latest published release
- `bleed` will be the latest build of main

## Development

### Before you run the app

The appliction expects a cassandra database setup to run. The environment variables you can control the connection with are:

- `CASSANDRA_CONTACT_POINTS` (split on ;)
- `CASSANDRA_LOCALDATACENTER`
- `CASSANDRA_USER`
- `CASSANDRA_PASSWORD`

However, if you set none up, you can simply hook up to a local database which can be run with `podman` with:

```
npm run dev:startlocaldb
```

You can then apply the migrations to this database with:

```
npm run dev:migrations
```

And seed some local test data with:

```bash
npm run dev:seedlocaltestdata
```

### Running the app

First, run the development server:

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

[API routes documentation](https://nextjs.org/docs/api-routes/introduction).

The `pages/api` directory is mapped to `/api/*`. Files in this directory are treated as [API routes](https://nextjs.org/docs/api-routes/introduction) instead of React pages.

### Learn More about Next.js

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!
