"use client";

import useSessionWithRefresh from "@/hooks/useSessionWithRefresh";
import { SIGN_OUT_CALLBACK_URL } from "@/utils/constants";
import { signOut } from "next-auth/react";

export default function Dashboard() {
  const session = useSessionWithRefresh();

  return (
    <main className="min-h-screen w-full px-4 py-8">
      <h1 className="text-4xl text-primary">
        Bem vindo {session.data?.user.username}
      </h1>
      <button onClick={() => signOut({ callbackUrl: SIGN_OUT_CALLBACK_URL })}>
        Sign out
      </button>
    </main>
  );
}
