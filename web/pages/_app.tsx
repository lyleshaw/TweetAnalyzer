import { NextUIProvider } from '@nextui-org/react';
import { AppProps } from 'next/app'

function App({ Component, pageProps }: AppProps) {
  return (
    <NextUIProvider>
      <Component { ...pageProps } />
    </NextUIProvider>
  );
}

export default App;
