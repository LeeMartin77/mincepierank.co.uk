import promclient from 'prom-client';
import http from 'http';

import { SvelteKitAuth } from '@auth/sveltekit';
import CredentialsProvider from '@auth/core/providers/credentials';
import GoogleProvider from '@auth/core/providers/google';
import { env } from '$env/dynamic/private';
import { sequence } from '@sveltejs/kit/hooks';

import { runMigrations } from '$lib/storage/migrations';

// This is because for whatever reason, vite tries
// to talk to the DB during the "build"
if (env.VITE_BUILDING !== 'building') {
  await runMigrations();
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const providers: any[] = [];

if (env.AUTHENTICATION === 'development') {
  providers.push(
    CredentialsProvider({
      name: 'Development',
      credentials: {
        username: { label: 'Username', type: 'text', placeholder: 'john' },
        password: { label: 'Password', type: 'password', placeholder: 'doe' }
      },
      async authorize(credentials) {
        if (credentials?.username === 'john' && credentials?.password === 'doe') {
          return {
            id: 'TEST_AUTH_john.doe@example.com',
            email: 'TEST_AUTH_john.doe@example.com',
            name: 'TEST_AUTH_LeeMartin'
          };
        }
        return null;
      }
    })
  );
}

if (env.GOOGLE_CLIENT_ID && env.GOOGLE_CLIENT_SECRET) {
  providers.push(
    GoogleProvider({
      clientId: env.GOOGLE_CLIENT_ID,
      clientSecret: env.GOOGLE_CLIENT_SECRET
    })
  );
}

const promprefix = process.env.PROMETHEUS_PREFIX || 'mpr_';

const promport = parseInt(process.env.PROMETHEUS_PORT || '9091');

promclient.collectDefaultMetrics({
  prefix: promprefix
});

const server = http.createServer(async (req, res) => {
  res.writeHead(200, { 'Content-Type': promclient.register.contentType });
  res.write(await promclient.register.metrics());
  res.end();
});

server.listen(promport, () => {
  console.log('Prometheus is listening on port ' + promport);
});

const requestsCounter = new promclient.Counter({
  name: `${promprefix}requests`,
  help: 'Counter of page requests',
  labelNames: ['type', 'method', 'generic_route', 'specific_route']
});

export const handle = sequence(
  ({ event, resolve }) => {
    requestsCounter
      .labels(
        event.url.pathname.startsWith('/api/') ? 'api' : 'page',
        event.request.method,
        event.route.id ?? 'N/A',
        event.url.pathname
      )
      .inc();
    return resolve(event);
  },
  SvelteKitAuth({
    providers: providers
  })
);
