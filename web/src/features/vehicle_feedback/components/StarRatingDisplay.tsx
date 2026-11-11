import { Star } from "lucide-react";
import { cn } from "@/shared/lib/utils";

interface Props {
  rating: number;
  className?: string;
}

export function StarRatingDisplay({ rating, className }: Props) {
  return (
    <span
      className={cn("inline-flex items-center gap-0.5", className)}
      aria-label={`${rating} de 5 estrellas`}
    >
      {[1, 2, 3, 4, 5].map((star) => (
        <Star
          key={star}
          aria-hidden
          className={cn(
            "size-4",
            star <= rating
              ? "fill-[var(--verde-500)] text-[var(--verde-700)]"
              : "text-muted-foreground",
          )}
        />
      ))}
    </span>
  );
}
