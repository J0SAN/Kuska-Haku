import { useState } from "react";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/api/httpClient";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/table";
import { useMarkSuggestionRead, useSuggestions } from "./hooks/useSuggestions";

const dateFormatter = new Intl.DateTimeFormat("es-PE", {
  dateStyle: "medium",
  timeStyle: "short",
});

export default function SuggestionsPage() {
  const [onlyNew, setOnlyNew] = useState(false);
  const { data: list = [], isLoading, error } = useSuggestions(onlyNew);
  const markRead = useMarkSuggestionRead();

  const handleMarkRead = (id: string) => {
    markRead.mutate(id, {
      onSuccess: () => toast.success("Sugerencia marcada como leída"),
      onError: (err) =>
        toast.error(getErrorMessage(err, "Error actualizando la sugerencia")),
    });
  };

  return (
    <section className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2>Sugerencias generales del servicio</h2>
          <p className="text-sm text-muted-foreground">
            Opiniones sobre el servicio en general. Para valorar un vehículo, el
            pasajero escanea el QR del bus.
          </p>
        </div>
        <label className="check mt-0">
          <input
            type="checkbox"
            checked={onlyNew}
            onChange={(e) => setOnlyNew(e.target.checked)}
          />
          Mostrar solo nuevas
        </label>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>
            {getErrorMessage(error, "Error cargando las sugerencias")}
          </AlertDescription>
        </Alert>
      )}

      <Card className="gap-0 py-0">
        <CardHeader className="border-b border-border py-4">
          <CardTitle className="text-secondary-foreground">
            Buzón de sugerencias generales
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="space-y-2 p-4">
              <Skeleton className="h-8 w-full" />
              <Skeleton className="h-8 w-full" />
              <Skeleton className="h-8 w-full" />
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Fecha</TableHead>
                  <TableHead>Usuario</TableHead>
                  <TableHead>Mensaje</TableHead>
                  <TableHead>Estado</TableHead>
                  <TableHead className="text-right">Acciones</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {list.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={5}
                      className="py-8 text-center text-muted-foreground"
                    >
                      No hay sugerencias para mostrar.
                    </TableCell>
                  </TableRow>
                )}
                {list.map((suggestion) => (
                  <TableRow key={suggestion.id}>
                    <TableCell className="whitespace-nowrap">
                      {dateFormatter.format(new Date(suggestion.createdAt))}
                    </TableCell>
                    <TableCell>{suggestion.userName}</TableCell>
                    <TableCell className="max-w-[420px] whitespace-pre-wrap">
                      {suggestion.message}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          suggestion.status === "new" ? "default" : "muted"
                        }
                      >
                        {suggestion.status === "new" ? "Nueva" : "Leída"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      {suggestion.status === "new" && (
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          disabled={
                            markRead.isPending &&
                            markRead.variables === suggestion.id
                          }
                          onClick={() => handleMarkRead(suggestion.id)}
                        >
                          Marcar como leída
                        </Button>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </section>
  );
}
