import { Star } from "lucide-react";
import { cn } from "@/shared/lib/utils";

const STARS = [1, 2, 3, 4, 5];

interface Props {
  value: number;
  onChange: (value: number) => void;
  name?: string;
}

// Grupo de radios dentro de un fieldset: operable con teclado y lector de
// pantalla, con toques de 48 px como mínimo.
export function StarRatingInput({ value, onChange, name = "rating" }: Props) {
  return (
    <fieldset className="space-y-2">
      <legend className="text-sm font-medium">¿Cómo fue tu viaje?</legend>
      <div className="flex items-center gap-1">
        {STARS.map((star) => (
          <label key={star} className="cursor-pointer rounded-md">
            <input
              type="radio"
              name={name}
              value={star}
              checked={value === star}
              onChange={() => onChange(star)}
              className="peer sr-only"
            />
            <Star
              aria-hidden
              className={cn(
                "size-12 rounded-md p-2 transition-colors peer-focus-visible:ring-2 peer-focus-visible:ring-ring peer-focus-visible:ring-offset-1",
                star <= value
                  ? "fill-[var(--verde-500)] text-[var(--verde-700)]"
                  : "text-muted-foreground",
              )}
            />
            <span className="sr-only">
              {star} {star === 1 ? "estrella" : "estrellas"}
            </span>
          </label>
        ))}
      </div>
    </fieldset>
  );
}
