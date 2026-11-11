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
import {
  useFeedbackSummary,
  useUpdateFeedbackStatus,
  useVehicleFeedback,
  useVehicleQrList,
} from "../hooks/useVehicleFeedback";
import type { FeedbackStatus } from "../vehicleFeedback.types";
import { StarRatingDisplay } from "./StarRatingDisplay";

const PERIODS = [7, 30, 90];

const dateFormatter = new Intl.DateTimeFormat("es-PE", {
  dateStyle: "medium",
  timeStyle: "short",
});

const STATUS_LABELS: Record<FeedbackStatus, string> = {
  new: "Nueva",
  read: "Leída",
  hidden: "Oculta",
};

const STATUS_BADGES: Record<FeedbackStatus, "default" | "muted" | "secondary"> =
  {
    new: "default",
    read: "secondary",
    hidden: "muted",
  };

export function FeedbackTab() {
  const [days, setDays] = useState(30);
  const [fleetId, setFleetId] = useState("");
  const [status, setStatus] = useState<FeedbackStatus | "">("");

  const qrList = useVehicleQrList();
  const summary = useFeedbackSummary(days);
  const feedback = useVehicleFeedback({ fleetId, status });
  const updateStatus = useUpdateFeedbackStatus();

  const handleStatus = (id: string, next: FeedbackStatus, message: string) => {
    updateStatus.mutate(
      { id, status: next },
      {
        onSuccess: () => toast.success(message),
        onError: (err) =>
          toast.error(getErrorMessage(err, "No se pudo actualizar")),
      },
    );
  };

  return (
    <div className="space-y-6">
      <Card className="gap-0 py-0">
        <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-3 border-b border-border py-4">
          <CardTitle className="text-secondary-foreground">
            Promedio por vehículo
          </CardTitle>
          <div className="flex gap-1" role="group" aria-label="Periodo">
            {PERIODS.map((period) => (
              <Button
                key={period}
                type="button"
                size="sm"
                variant={days === period ? "default" : "outline"}
                aria-pressed={days === period}
                onClick={() => setDays(period)}
              >
                {period} días
              </Button>
            ))}
          </div>
        </CardHeader>
        <CardContent className="p-0">
          {summary.isLoading ? (
            <div className="space-y-2 p-4">
              <Skeleton className="h-8 w-full" />
              <Skeleton className="h-8 w-full" />
            </div>
          ) : summary.error ? (
            <div className="p-4">
              <Alert variant="destructive">
                <AlertDescription>
                  {getErrorMessage(summary.error, "Error cargando los promedios")}
                </AlertDescription>
              </Alert>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Bus</TableHead>
                  <TableHead>Ruta</TableHead>
                  <TableHead>Promedio</TableHead>
                  <TableHead>Calificaciones</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(summary.data ?? []).length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={4}
                      className="py-8 text-center text-muted-foreground"
                    >
                      No hay calificaciones en este periodo.
                    </TableCell>
                  </TableRow>
                )}
                {(summary.data ?? []).map((row) => (
                  <TableRow key={row.fleetId ?? row.busCode}>
                    <TableCell className="font-medium">{row.busCode}</TableCell>
                    <TableCell>{row.routeName}</TableCell>
                    <TableCell>
                      <span className="inline-flex items-center gap-2">
                        <StarRatingDisplay rating={Math.round(row.average)} />
                        <span className="font-semibold">
                          {row.average.toFixed(1)}
                        </span>
                      </span>
                    </TableCell>
                    <TableCell>{row.total}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <Card className="gap-0 py-0">
        <CardHeader className="border-b border-border py-4">
          <CardTitle className="text-secondary-foreground">
            Comentarios de pasajeros
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4 p-4">
          <div className="flex flex-wrap gap-3">
            <label className="check mt-0">
              Vehículo
              <select
                className="h-9 rounded-md border border-input bg-card px-2 py-0 text-sm"
                value={fleetId}
                onChange={(e) => setFleetId(e.target.value)}
              >
                <option value="">Todos</option>
                {(qrList.data ?? []).map((item) => (
                  <option key={item.fleetId} value={item.fleetId}>
                    {item.busCode}
                  </option>
                ))}
              </select>
            </label>
            <label className="check mt-0">
              Estado
              <select
                className="h-9 rounded-md border border-input bg-card px-2 py-0 text-sm"
                value={status}
                onChange={(e) =>
                  setStatus(e.target.value as FeedbackStatus | "")
                }
              >
                <option value="">Todos</option>
                <option value="new">Nuevas</option>
                <option value="read">Leídas</option>
                <option value="hidden">Ocultas</option>
              </select>
            </label>
          </div>

          {feedback.error && (
            <Alert variant="destructive">
              <AlertDescription>
                {getErrorMessage(feedback.error, "Error cargando los comentarios")}
              </AlertDescription>
            </Alert>
          )}

          {feedback.isLoading ? (
            <div className="space-y-2">
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
            </div>
          ) : (
            <div className="overflow-x-auto rounded-lg border border-border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Fecha</TableHead>
                    <TableHead>Bus</TableHead>
                    <TableHead>Conductor</TableHead>
                    <TableHead>Calificación</TableHead>
                    <TableHead>Comentario</TableHead>
                    <TableHead>Servicio</TableHead>
                    <TableHead>Estado</TableHead>
                    <TableHead className="text-right">Acciones</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(feedback.data ?? []).length === 0 && (
                    <TableRow>
                      <TableCell
                        colSpan={8}
                        className="py-8 text-center text-muted-foreground"
                      >
                        No hay comentarios para estos filtros.
                      </TableCell>
                    </TableRow>
                  )}
                  {(feedback.data ?? []).map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="whitespace-nowrap">
                        {dateFormatter.format(new Date(item.createdAt))}
                      </TableCell>
                      <TableCell className="font-medium whitespace-nowrap">
                        {item.busCode}
                      </TableCell>
                      <TableCell className="whitespace-nowrap">
                        {item.driverName || "—"}
                      </TableCell>
                      <TableCell>
                        <span className="inline-flex items-center gap-2">
                          <StarRatingDisplay rating={item.rating} />
                          <span>{item.rating}</span>
                        </span>
                      </TableCell>
                      <TableCell className="max-w-[320px] whitespace-pre-wrap">
                        {item.comment || "—"}
                      </TableCell>
                      <TableCell>
                        <Badge variant={item.busLive ? "default" : "muted"}>
                          {item.busLive ? "En servicio" : "Sin GPS activo"}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <Badge variant={STATUS_BADGES[item.status]}>
                          {STATUS_LABELS[item.status]}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-wrap justify-end gap-1">
                          {item.status === "new" && (
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              onClick={() =>
                                handleStatus(
                                  item.id,
                                  "read",
                                  "Calificación marcada como leída",
                                )
                              }
                            >
                              Marcar como leída
                            </Button>
                          )}
                          {item.status !== "hidden" ? (
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              className="text-destructive hover:text-destructive"
                              onClick={() =>
                                handleStatus(
                                  item.id,
                                  "hidden",
                                  "Calificación oculta del promedio",
                                )
                              }
                            >
                              Ocultar
                            </Button>
                          ) : (
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={() =>
                                handleStatus(
                                  item.id,
                                  "new",
                                  "Calificación restaurada",
                                )
                              }
                            >
                              Mostrar
                            </Button>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
