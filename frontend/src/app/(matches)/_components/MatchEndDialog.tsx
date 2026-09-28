"use client";

import { Trophy, XCircle } from "lucide-react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface MatchEndDialogProps {
  abandoned: boolean;
  open: boolean;
  winnerId: string | null;
}

export default function MatchEndDialog({
  open,
  abandoned,
  winnerId,
}: MatchEndDialogProps) {
  const router = useRouter();

  return (
    <Dialog open={open}>
      <DialogContent showCloseButton={false}>
        <DialogHeader className="items-center text-center">
          {abandoned ? (
            <XCircle className="size-10 text-zinc-400" />
          ) : (
            <Trophy className="size-10 text-yellow-500" />
          )}

          <DialogTitle>
            {abandoned ? "Match abandoned" : "Match finished"}
          </DialogTitle>

          <DialogDescription>
            {abandoned
              ? "Match ended without a winner."
              : `Winner ID: ${winnerId}`}
          </DialogDescription>
        </DialogHeader>

        <Button onClick={() => router.push("/dashboard")}>
          Back to dashboard
        </Button>
      </DialogContent>
    </Dialog>
  );
}
