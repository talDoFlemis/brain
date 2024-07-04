export { auth as middleware } from "@/auth";

// Define the routes you want to protect here
export const config = { matcher: ["/dashboard/:path"] };
