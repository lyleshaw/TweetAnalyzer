import { AppProps } from "next/app";
import "../styles/globals.css";
import { AppLayout } from "@/components";

function App({ Component, pageProps }: AppProps) {
  return (
    <AppLayout>
      <Component {...pageProps} />
    </AppLayout>
  );
}

export default App;
