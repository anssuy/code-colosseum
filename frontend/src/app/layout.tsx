import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import type { ReactNode } from "react";

import { AuthProvider } from "@/providers/AuthProvider";
import "./globals.css";

import { Toaster } from "@/components/ui/toast";
import { TooltipProvider } from "@/components/ui/tooltip";
import { MatchProvider } from "@/providers/MatchProvider";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Code Colosseum",
  description: "Compete in real-time coding battles.",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
      lang="en"
    >
      <body className="flex min-h-full flex-col font-geist-sans">
        <TooltipProvider>
          <AuthProvider>
            <MatchProvider>
              {children}
              <Toaster />
            </MatchProvider>
          </AuthProvider>
        </TooltipProvider>
      </body>
    </html>
  );
}
