import { IoLogOutOutline } from "react-icons/io5";

function Header() {
  return (
    <header className="flex justify-end h-16 px-6 w-full bg-background">
      <div className="flex items-center gap-8 w-fit bg-foreground text-secondary-foreground p-4 rounded-lg">
        <span>Erico Bayern de munique</span>
        <IoLogOutOutline className="w-6 h-6" />
      </div>
    </header>
  );
}

export default Header;
