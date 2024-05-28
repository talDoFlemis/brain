import type { Metadata } from "next";

import RegisterSection from "@/components/pages/signup/register-section";

export const metadata: Metadata = {
  title: "Brain Test • Sign up",
};

export default function SignUp() {
  return (
    <main className="dark grid grid-cols-4 min-h-screen w-full bg-background font-inter">
      <RegisterSection />
      <section className="hidden lg:block col-span-2 bg-cover bg-[url('/brain-surface.svg')]"></section>
    </main>
  );
}
