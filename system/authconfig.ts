// https://next-auth.js.org/getting-started/introduction
import CredentialsProvider from "next-auth/providers/credentials";
import GitHubProvider from "next-auth/providers/github";
import GoogleProvider from "next-auth/providers/google";

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

if (process.env.GITHUB_ID && process.env.GITHUB_SECRET) {
  providers.push(
    GitHubProvider({
      clientId: process.env.GITHUB_ID,
      clientSecret: process.env.GITHUB_SECRET,
      profile: (profile) => {
        return {
          ...profile,
          id: "GITHUB_AUTH_" + profile.id,
          name: "GITHUB_AUTH_" + profile.login,
        };
      },
    })
  );
}

if (process.env.GOOGLE_CLIENT_ID && process.env.GOOGLE_CLIENT_SECRET) {
  providers.push(
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET,
    })
  );
}

export const authOptions = {
  // Configure one or more authentication providers
  providers,
};
