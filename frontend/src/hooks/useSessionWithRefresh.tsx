import { useEffect } from "react";
import { signOut, useSession } from "next-auth/react";
import { SIGN_OUT_CALLBACK_URL } from "@/utils/constants";

// Checks whether there is an error within the session and signs user out
function useSessionWithRefresh() {
  const session = useSession();

  useEffect(() => {
    if (session.data?.error) {
      signOut({ callbackUrl: SIGN_OUT_CALLBACK_URL });
    }
  }, [session]);

  return session;
}

export default useSessionWithRefresh;
