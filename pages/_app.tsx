import "../styles/globals.css";
import type { AppProps } from "next/app";
import { useEffect, useMemo, useState } from "react";
import {
  Box,
  createTheme,
  CssBaseline,
  Divider,
  ThemeProvider,
  useMediaQuery,
  Link,
} from "@mui/material";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { Container } from "@mui/system";
import {
  SideNavigationComponent,
  BottomNavigationComponent,
} from "../components/navigation";
import { SessionProvider } from "next-auth/react";
import nxLink from "next/link";
import Head from "next/head";

export default function App({
  Component,
  pageProps: { session, ...pageProps },
}: AppProps) {
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");
  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode: prefersDarkMode ? "dark" : "light",
        },
      }),
    [prefersDarkMode]
  );

  const [isDesktop, setDesktop] = useState(false);

  useEffect(() => {
    const updateMedia = () => {
      setDesktop(window.innerWidth > theme.breakpoints.values.sm);
    };
    updateMedia();
    window.addEventListener("resize", updateMedia);
    return () => window.removeEventListener("resize", updateMedia);
  }, [theme.breakpoints.values.sm]);

  const containerClassName = isDesktop
    ? "main-container-nonmobile"
    : "main-container-mobile";

  return (
    <SessionProvider session={session}>
      <ThemeProvider theme={theme}>
        <CssBaseline enableColorScheme />
        <Head>
          <link rel="icon" href="/favicon.ico" />
        </Head>
        <Box sx={{ display: "flex" }}>
          {isDesktop && <SideNavigationComponent />}
          <Box component="main" style={{ width: "100%" }}>
            <Container className={containerClassName}>
              <Component {...pageProps} />
            </Container>
            <Divider style={{ marginTop: "1em", marginBottom: "1em" }} />
            <div
              style={{
                width: "100%",
                textAlign: "center",
                marginBottom: isDesktop ? "1em" : "5em",
              }}
            >
              <Link component={nxLink} href={"/about"}>
                About
              </Link>
            </div>
          </Box>
          {!isDesktop && <BottomNavigationComponent />}
        </Box>
      </ThemeProvider>
    </SessionProvider>
  );
}
