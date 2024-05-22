import Image from "next/image";
import Link from "next/link";
import { IoListOutline } from "react-icons/io5";
import { IoMdCheckmarkCircleOutline } from "react-icons/io";
import { FaPlus, FaRobot } from "react-icons/fa6";
import { FaPlay } from "react-icons/fa";
import { ReactNode } from "react";
import { IconType } from "react-icons/lib";
import { cn } from "@/lib/utils";

type SidebarLinkProps = {
  href: string;
  Icon: IconType;
  text: ReactNode;
  active: boolean;
};

function LinkContainer({ children }: { children: ReactNode }) {
  return (
    <div className="flex flex-col gap-4 w-full bg-foreground p-4 rounded-lg">
      {children}
    </div>
  );
}

function SidebarLink({ href, Icon, text, active }: SidebarLinkProps) {
  return (
    <Link
      className={cn(
        "flex items-center gap-2 px-2 py-2 text-lg text-secondary-foreground bg-foreground rounded-md",
        active ? "bg-secondary" : "",
      )}
      href={href}
    >
      <Icon className="w-6 h-6" />
      <span>{text}</span>
    </Link>
  );
}

function Sidebar() {
  return (
    <aside className="flex flex-col h-full gap-2 w-96 p-6">
      <div className="flex items-center gap-2 w-full bg-foreground p-6 rounded-lg">
        <Image
          className="self-center"
          src="/brain-logo.svg"
          alt="brain.test logo"
          width={50}
          height={50}
          priority
        />
        <span className="text-3xl uppercase text-white">Brain</span>
      </div>
      <LinkContainer>
        <SidebarLink
          href="#"
          active={true}
          Icon={IoListOutline}
          text="Quizzes criados"
        />
        <SidebarLink
          href="#"
          active={false}
          Icon={IoMdCheckmarkCircleOutline}
          text="Quizzes jogados"
        />
      </LinkContainer>

      <LinkContainer>
        <SidebarLink href="#" active={false} Icon={FaPlus} text="Criar Quiz" />
        <SidebarLink
          href="#"
          active={false}
          Icon={FaRobot}
          text={
            <>
              Quiz gerado por <span className="text-accent">AI</span>
            </>
          }
        />
      </LinkContainer>

      <LinkContainer>
        <SidebarLink
          href="#"
          active={false}
          Icon={FaPlay}
          text="Entrar na sessão"
        />
      </LinkContainer>
    </aside>
  );
}

export default Sidebar;
