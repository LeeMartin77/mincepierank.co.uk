// https://next-auth.js.org/getting-started/introduction
import NextAuth from "next-auth";
import { authOptions } from "../../../system/authconfig";

export default NextAuth(authOptions);
