// https://next-auth.js.org/getting-started/introduction
import CredentialsProvider from "next-auth/providers/credentials";
import FacebookProvider from "next-auth/providers/facebook";
import GitHubProvider from "next-auth/providers/github";

const providers: any[] = [];

if (process.env.NODE_ENV === "development") {
  providers.push(
    CredentialsProvider({
      name: "Development",
      credentials: {
        username: { label: "Username", type: "text", placeholder: "john" },
        password: { label: "Password", type: "password", placeholder: "doe" },
      },
      async authorize(credentials) {
        if (
          credentials?.username === "john" &&
          credentials?.password === "doe"
        ) {
          return {
            id: "TEST_AUTH_test-user-id",
            email: "TEST_AUTH_john.doe@example.com",
            name: "TEST_AUTH_LeeMartin",
          };
        }
        return null;
      },
    })
  );
}

if (process.env.FACEBOOK_CLIENT_ID && process.env.FACEBOOK_CLIENT_SECRET) {
  providers.push(
    FacebookProvider({
      clientId: process.env.FACEBOOK_CLIENT_ID,
      clientSecret: process.env.FACEBOOK_CLIENT_SECRET,
      async profile(profile) {
        return {
          ...profile,
          name: "FACEBOOK_AUTH_" + profile.name,
          email: "FACEBOOK_AUTH_" + profile.email,
          id: "FACEBOOK_AUTH_" + profile.id,
        };
      },
    })
  );
}

if (process.env.GITHUB_ID && process.env.GITHUB_SECRET) {
  providers.push(
    GitHubProvider({
      clientId: process.env.GITHUB_ID,
      clientSecret: process.env.GITHUB_SECRET,
      async profile(profile) {
        return {
          ...profile,
          name: "GITHUB_AUTH_" + profile.name,
          email: "GITHUB_AUTH_" + profile.email,
          id: "GITHUB_AUTH_" + profile.id,
        };
      },
    })
  );
}

export const authOptions = {
  // Configure one or more authentication providers
  providers,
};
