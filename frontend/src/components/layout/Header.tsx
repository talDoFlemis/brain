"use client";
import { IoLogOutOutline } from "react-icons/io5";
import { Skeleton } from "@/components/ui/skeleton";
import { signOut } from "next-auth/react";
import useSessionWithRefresh from "@/hooks/useSessionWithRefresh";

function Header() {
  const { data } = useSessionWithRefresh();

  return (
    <header className="flex justify-end h-16 px-6 w-full bg-background">
      <div className="flex items-center gap-8 w-fit bg-foreground text-secondary-foreground p-4 rounded-lg">
        <span>{data && data.user.username}</span>
        <IoLogOutOutline
          className="w-6 h-6 cursor-pointer transition-colors duration-200 ease-in hover:text-accent"
          onClick={() => signOut()}
        />
      </div>
    </header>
  );
}

export default Header;

export function LoadingSkeleton() {
  return (
    <header className="flex justify-end h-16 px-6 w-full bg-background">
      <div className="flex items-center gap-8 w-fit bg-foreground text-secondary-foreground p-4 rounded-lg">
        <Skeleton className="h-6 w-20 bg-primary-foreground" />
        <Skeleton className="h-6 w-6 bg-primary-foreground rounded-full" />
      </div>
    </header>
  );
}
