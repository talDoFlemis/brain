import type { Metadata } from "next";

import LoginSection from "@/components/pages/signin/login-section";

export const metadata: Metadata = {
  title: "Brain Test • Sign in",
};

export default function SignIn() {
  return (
    <main className="dark grid grid-cols-4 min-h-screen w-full bg-background font-inter">
      <section className="hidden lg:block col-span-2 bg-cover bg-[url('/brain-surface.svg')]"></section>
      <LoginSection />
    </main>
  );
}
