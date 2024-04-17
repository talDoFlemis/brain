import Image from "next/image";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import LoginForm from "./login-form";

function LoginSection() {
  return (
    <section className="col-span-4 lg:col-span-2 justify-self-center flex flex-col w-full py-8 px-4 gap-8 max-w-lg">
      <Image
        className="self-center"
        src="/brain-logo-white.svg"
        alt="brain.test logo"
        width={225}
        height={50}
        priority
      />
      <div className="self-center flex flex-col items-center gap-4">
        <h2 className="text-3xl text-foreground font-semibold">
          Logue em uma conta
        </h2>
        <h3 className="text-sm text-foreground/50 font-semibold">
          Digite suas credenciais
        </h3>
      </div>
      <LoginForm />
      <p className="text-sm text-foreground">
        Nao tem uma conta ainda?
        <Button variant="link" size="sm" asChild>
          <Link href="#">Registrar</Link>
        </Button>
      </p>
    </section>
  );
}

export default LoginSection;
