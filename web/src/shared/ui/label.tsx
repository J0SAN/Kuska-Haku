import * as React from "react";
import * as LabelPrimitive from "@radix-ui/react-label";
import { cn } from "@/shared/lib/utils";

function Label({
  className,
  ...props
}: React.ComponentProps<typeof LabelPrimitive.Root>) {
  return (
    <LabelPrimitive.Root
      data-slot="label"
      className={cn(
        "mt-0 flex items-center gap-2 text-sm leading-none font-medium select-none data-disabled:opacity-50",
        className,
      )}
      {...props}
    />
  );
}

export { Label };
