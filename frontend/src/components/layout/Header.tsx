import { authOptions } from "@/providers/auth";
import authService from "@/services/auth-service";
import { getServerSession } from "next-auth";
import { IoLogOutOutline } from "react-icons/io5";
import { Skeleton } from "@/components/ui/skeleton";

async function Header() {
  const session = await getServerSession(authOptions);
  const user = await authService.getUserInfo(session?.access_token ?? "");
  return (
    <header className="flex justify-end h-16 px-6 w-full bg-background">
      <div className="flex items-center gap-8 w-fit bg-foreground text-secondary-foreground p-4 rounded-lg">
        <span>{user.username}</span>
        <IoLogOutOutline className="w-6 h-6" />
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
