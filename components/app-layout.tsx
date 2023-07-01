import { ReactNode } from "react";
import { Footer } from "./footer";

interface IAppLayoutProps {
  children: ReactNode;
}
export function AppLayout({ children }: IAppLayoutProps) {
  return (
    <>
      <main>{children}</main>
      <Footer />
    </>
  );
}
