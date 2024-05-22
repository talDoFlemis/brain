import { ReactNode } from "react";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";

type DashboardLayoutProps = {
  children: ReactNode;
};

function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <div className="dark flex h-screen overflow-y-hidden bg-no-repeat bg-left-bottom bg-background bg-[url('/brain-surface-sidebar.svg')]">
      <Sidebar />
      <div className="flex h-screen w-full flex-col py-6 bg-no-repeat bg-right-bottom bg-[url('/brain-surface-bg.svg')]">
        <Header />
        <main className="flex h-screen w-full flex-col overflow-y-auto px-4 pt-8 xl:px-8">
          {children}
        </main>
      </div>
    </div>
  );
}

export default DashboardLayout;
