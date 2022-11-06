// https://next-auth.js.org/getting-started/introduction
import CredentialsProvider, {
  CredentialInput,
  CredentialsConfig,
} from "next-auth/providers/credentials";

const providers: CredentialsConfig<Record<string, CredentialInput>>[] = [];

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

export const authOptions = {
  // Configure one or more authentication providers
  providers,
};
