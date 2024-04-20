export { default } from "next-auth/middleware";

// Define the routes you want to protect here
export const config = { matcher: ["/"] }
