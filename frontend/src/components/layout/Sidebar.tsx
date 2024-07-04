"use client";

import Image from "next/image";
import Link from "next/link";
import { IoListOutline } from "react-icons/io5";
import { IoMdCheckmarkCircleOutline } from "react-icons/io";
import { FaPlus } from "react-icons/fa6";
import { FaPlay } from "react-icons/fa";
import { ReactNode } from "react";
import { LuBrainCircuit } from "react-icons/lu";
import { IconType } from "react-icons/lib";
import { cn } from "@/lib/utils";
import { usePathname } from "next/navigation";

type SidebarLinkProps = {
  href: string;
  Icon: IconType;
  text: ReactNode;
  currentPathname: string;
};

function LinkContainer({ children }: { children: ReactNode }) {
  return (
    <div className="flex flex-col gap-4 w-full bg-foreground p-4 rounded-lg">
      {children}
    </div>
  );
}

function SidebarLink({ href, Icon, text, currentPathname }: SidebarLinkProps) {
  const active = href === currentPathname;
  return (
    <Link
      className={cn(
        "flex items-center gap-2 px-2 py-2 text-lg text-secondary-foreground bg-foreground rounded-md hover:bg-secondary transition-colors duration-300",
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
  const pathname = usePathname();
  return (
    <aside className="flex flex-col h-full gap-2 w-96 p-6">
      <div className="flex items-center gap-2 w-full bg-foreground p-6 rounded-lg justify-center">
        <Image
          className="self-center"
          src="/brain-logo.svg"
          alt="brain.test logo"
          width={50}
          height={50}
          priority
        />
        <span className="text-3xl uppercase text-white font-jua">Brain</span>
      </div>
      <LinkContainer>
        <SidebarLink
          href="/dashboard"
          currentPathname={pathname}
          Icon={IoListOutline}
          text="Quizzes criados"
        />
        <SidebarLink
          href="/dashboard/played"
          currentPathname={pathname}
          Icon={IoMdCheckmarkCircleOutline}
          text="Quizzes jogados"
        />
      </LinkContainer>

      <LinkContainer>
        <SidebarLink
          href="/dashboard/create-quiz"
          currentPathname={pathname}
          Icon={FaPlus}
          text="Criar Quiz"
        />
        <SidebarLink
          href="#"
          currentPathname={pathname}
          Icon={LuBrainCircuit}
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
          currentPathname={pathname}
          Icon={FaPlay}
          text="Entrar na sessão"
        />
      </LinkContainer>
    </aside>
  );
}

export default Sidebar;
