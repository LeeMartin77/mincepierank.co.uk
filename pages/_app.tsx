import "../styles/globals.css";
import type { AppProps } from "next/app";
import { useEffect, useMemo, useState } from "react";
import {
  Box,
  createTheme,
  CssBaseline,
  ThemeProvider,
  useMediaQuery,
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

export default function App({ Component, pageProps }: AppProps) {
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
    <ThemeProvider theme={theme}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: "flex" }}>
        {isDesktop && <SideNavigationComponent />}
        <Box component="main" style={{ width: "100%" }}>
          <Container className={containerClassName}>
            <Component {...pageProps} />
          </Container>
        </Box>
        {!isDesktop && <BottomNavigationComponent />}
      </Box>
    </ThemeProvider>
  );
}
