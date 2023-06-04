import { SvelteKitAuth } from '@auth/sveltekit';
import CredentialsProvider from '@auth/core/providers/credentials';
import GoogleProvider from '@auth/core/providers/google';
import { env } from '$env/dynamic/private';

const providers: any[] = [];

if (env.NODE_ENV === 'development') {
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

export const handle = SvelteKitAuth({
  providers: providers
});
