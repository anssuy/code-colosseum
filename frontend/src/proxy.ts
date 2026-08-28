import { NextProxy, NextRequest, NextResponse, ProxyConfig } from "next/server";

const protectedRoutes = ["/dashboard"];
const guestRoutes = ["/login", "/register"];

export const proxy: NextProxy = (request: NextRequest) => {
  const { pathname } = request.nextUrl;
  const session = request.cookies.get("access_token");

  const isAuthenticated = !!session;

  const isProtectedRoute = protectedRoutes.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`),
  );

  const isGuestRoute = guestRoutes.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`),
  );

  if (isProtectedRoute && !isAuthenticated) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  if (isGuestRoute && isAuthenticated) {
    return NextResponse.redirect(new URL("/dashboard", request.url));
  }

  return NextResponse.next();
};

export const config: ProxyConfig = {
  matcher: ["/dashboard/:path*", "/login", "/register"],
};
