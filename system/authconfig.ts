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
            id: "test-user-id",
            email: "john.doe@example.com",
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
    })
  );
}

if (process.env.GITHUB_ID && process.env.GITHUB_SECRET) {
  providers.push(
    GitHubProvider({
      clientId: process.env.GITHUB_ID,
      clientSecret: process.env.GITHUB_SECRET,
    })
  );
}

export const authOptions = {
  // Configure one or more authentication providers
  providers,
};
