import { jwtVerify } from "jose";
import { NextProxy, NextRequest, NextResponse, ProxyConfig } from "next/server";

const publicRoutes = ["/login", "/register"];

const secret = new TextEncoder().encode(process.env.JWT_SECRET);

export const verifyAccessToken = async (token: string) => {
  try {
    const { payload } = await jwtVerify(token, secret, {
      algorithms: ["HS256"],
      issuer: "code-colosseum",
      requiredClaims: ["sub", "iat", "exp"],
    });

    return !!payload.sub;
  } catch (error) {
    console.error("JWT verification failed:", error);
    return false;
  }
};

export const proxy: NextProxy = async (request: NextRequest) => {
  const { pathname } = request.nextUrl;

  const isPublicRoute = publicRoutes.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`),
  );

  const token = request.cookies.get("access_token")?.value;

  const isAuthenticated = token ? await verifyAccessToken(token) : false;

  if (!isPublicRoute && !isAuthenticated) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  if (isPublicRoute && isAuthenticated) {
    return NextResponse.redirect(new URL("/dashboard", request.url));
  }

  return NextResponse.next();
};

export const config: ProxyConfig = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};
