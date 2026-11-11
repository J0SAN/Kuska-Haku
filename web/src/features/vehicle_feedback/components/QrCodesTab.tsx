import { useState } from "react";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/api/httpClient";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/ui/alert-dialog";
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
  useGenerateVehicleQr,
  useRevokeVehicleQr,
  useRotateVehicleQr,
  useVehicleQrList,
} from "../hooks/useVehicleFeedback";
import type { VehicleQRSummary } from "../vehicleFeedback.types";
import { QrPreviewDialog } from "./QrPreviewDialog";

type PendingAction = "rotate" | "revoke";

interface ConfirmState {
  action: PendingAction;
  item: VehicleQRSummary;
}

export function QrCodesTab() {
  const { data: list = [], isLoading, error } = useVehicleQrList();
  const generate = useGenerateVehicleQr();
  const rotate = useRotateVehicleQr();
  const revoke = useRevokeVehicleQr();

  const [selected, setSelected] = useState<VehicleQRSummary | null>(null);
  const [previewOpen, setPreviewOpen] = useState(false);
  const [confirm, setConfirm] = useState<ConfirmState | null>(null);

  const openPreview = (item: VehicleQRSummary) => {
    setSelected(item);
    setPreviewOpen(true);
  };

  const handleGenerate = (item: VehicleQRSummary) => {
    generate.mutate(item.fleetId, {
      onSuccess: (qr) => {
        toast.success(`QR generado para ${item.busCode}`);
        openPreview({
          ...item,
          qrCodeId: qr.id,
          token: qr.token,
          createdAt: qr.createdAt,
        });
      },
      onError: (err) =>
        toast.error(getErrorMessage(err, "No se pudo generar el QR")),
    });
  };

  const handleConfirmed = () => {
    if (!confirm) return;
    const { action, item } = confirm;

    if (action === "rotate") {
      rotate.mutate(item.fleetId, {
        onSuccess: (qr) => {
          toast.success(`QR regenerado para ${item.busCode}`);
          setConfirm(null);
          openPreview({ ...item, qrCodeId: qr.id, token: qr.token });
        },
        onError: (err) =>
          toast.error(getErrorMessage(err, "No se pudo regenerar el QR")),
      });
      return;
    }

    revoke.mutate(item.fleetId, {
      onSuccess: () => {
        toast.success(`QR revocado para ${item.busCode}`);
        setConfirm(null);
      },
      onError: (err) =>
        toast.error(getErrorMessage(err, "No se pudo revocar el QR")),
    });
  };

  return (
    <div className="space-y-4">
      {error && (
        <Alert variant="destructive">
          <AlertDescription>
            {getErrorMessage(error, "Error cargando los códigos QR")}
          </AlertDescription>
        </Alert>
      )}

      <Card className="gap-0 py-0">
        <CardHeader className="border-b border-border py-4">
          <CardTitle className="text-secondary-foreground">
            Códigos QR por vehículo
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="space-y-2 p-4">
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
            </div>
          ) : (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Bus</TableHead>
                    <TableHead>Ruta</TableHead>
                    <TableHead>Conductor</TableHead>
                    <TableHead>QR</TableHead>
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
                        Todavía no hay vehículos registrados.
                      </TableCell>
                    </TableRow>
                  )}
                  {list.map((item) => {
                    const generating =
                      generate.isPending && generate.variables === item.fleetId;
                    return (
                      <TableRow key={item.fleetId}>
                        <TableCell className="font-medium whitespace-nowrap">
                          {item.busCode}
                        </TableCell>
                        <TableCell className="whitespace-nowrap">
                          {item.routeName}
                        </TableCell>
                        <TableCell className="whitespace-nowrap">
                          {item.driverName || "—"}
                        </TableCell>
                        <TableCell>
                          {item.token ? (
                            <Badge>● Activo</Badge>
                          ) : (
                            <Badge variant="muted">Sin QR</Badge>
                          )}
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-wrap justify-end gap-2">
                            {item.token ? (
                              <>
                                <Button
                                  type="button"
                                  variant="outline"
                                  size="sm"
                                  onClick={() => openPreview(item)}
                                >
                                  Ver
                                </Button>
                                <Button
                                  type="button"
                                  variant="outline"
                                  size="sm"
                                  onClick={() =>
                                    setConfirm({ action: "rotate", item })
                                  }
                                >
                                  Regenerar
                                </Button>
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="sm"
                                  className="text-destructive hover:text-destructive"
                                  onClick={() =>
                                    setConfirm({ action: "revoke", item })
                                  }
                                >
                                  Revocar
                                </Button>
                              </>
                            ) : (
                              <Button
                                type="button"
                                size="sm"
                                disabled={generating}
                                onClick={() => handleGenerate(item)}
                              >
                                {generating ? "Generando…" : "Generar QR"}
                              </Button>
                            )}
                          </div>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>

      <QrPreviewDialog
        item={selected}
        open={previewOpen}
        onOpenChange={setPreviewOpen}
      />

      <AlertDialog
        open={confirm !== null}
        onOpenChange={(open) => {
          if (!open) setConfirm(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {confirm?.action === "rotate"
                ? "Regenerar el código QR"
                : "Revocar el código QR"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              El QR impreso dejará de funcionar. Tendrás que reemplazar el
              sticker.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              className={
                confirm?.action === "revoke"
                  ? "bg-destructive text-primary-foreground hover:bg-destructive/90"
                  : undefined
              }
              onClick={handleConfirmed}
            >
              {confirm?.action === "rotate" ? "Regenerar" : "Revocar"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
