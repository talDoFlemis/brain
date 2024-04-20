"use client";

import Image from "next/image";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import RegisterForm, { RegisterFormSchema } from "./register-form";

function RegisterSection() {
  async function submitForm(v: RegisterFormSchema) {
    await new Promise((resolve) => setTimeout(resolve, 0));
    console.log(v);
    return true;
  }

  return (
    <section className="col-span-4 lg:col-span-2 justify-self-center flex flex-col w-full py-12 px-4 gap-8 max-w-lg">
      <Image
        className="self-center"
        src="/brain-logo-white.svg"
        alt="brain.test logo"
        width={175}
        height={50}
        priority
      />
      <div className="self-center flex flex-col items-center gap-4">
        <h2 className="text-3xl text-foreground font-semibold">
          Registrar uma conta
        </h2>
        <h3 className="text-sm text-foreground/50 font-semibold">
          Digite suas credenciais
        </h3>
      </div>
      <RegisterForm submitForm={submitForm} />
      <p className="text-sm text-foreground">
        Ja possui uma conta?
        <Button variant="link" size="sm" asChild>
          <Link href="/sign-in">Logar</Link>
        </Button>
      </p>
    </section>
  );
}

export default RegisterSection;
